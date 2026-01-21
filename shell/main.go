package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
)

// TaskStatus 任务状态
type TaskStatus int

const (
	StatusPending   TaskStatus = iota // 0
	StatusRunning                     // 1
	StatusSuccess                     // 2
	StatusFailed                      // 3
	StatusTimeout                     // 4
	StatusCancelled                   // 5
)

func (s TaskStatus) String() string {
	switch s {
	case StatusPending:
		return "待处理"
	case StatusRunning:
		return "运行中"
	case StatusSuccess:
		return "成功"
	case StatusFailed:
		return "失败"
	case StatusTimeout:
		return "超时"
	case StatusCancelled:
		return "取消"
	default:
		return "未知"
	}
}

// Task 任务定义
type Task struct {
	ID           string        // 任务ID
	Name         string        // 任务名称
	Cmd          string        // 执行命令
	Args         []string      // 命令参数
	Timeout      time.Duration // 超时时间
	RetryCount   int           // 重试次数
	RetryDelay   time.Duration // 重试延迟
	MaxOutput    int           // 最大输出行数
	Env          []string      // 环境变量
	WorkDir      string        // 工作目录
	Dependencies []string      // 依赖的任务ID
}

// TaskResult 任务执行结果
type TaskResult struct {
	TaskID     string        // 任务ID
	TaskName   string        // 任务名称
	Status     TaskStatus    // 任务状态
	StartTime  time.Time     // 开始时间
	EndTime    time.Time     // 结束时间
	Duration   time.Duration // 持续时间
	ExitCode   int           // 退出code
	Output     string        // 输出内容
	Error      error         // 错误信息
	RetryCount int           // 重试次数
}

// Scheduler 调度器
type Scheduler struct {
	maxWorkers      int                    // 最大并发数
	tasks           map[string]*Task       // 所有任务
	taskResults     map[string]*TaskResult // 所有任务结果
	taskQueue       chan *Task             // 任务队列
	taskResultQueue chan *TaskResult       // 结果队列
	wg              sync.WaitGroup         // 等待组
	mu              sync.Mutex             // 读写锁
	ctx             context.Context        // 上下文
	cancel          context.CancelFunc     // 取消函数
	isRunning       bool                   // 是否正在运行
	completedTasks  map[string]bool        // 已完成任务
}

// NewScheduler 创建调度器
func NewScheduler(maxWorkers int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		maxWorkers:      maxWorkers,
		tasks:           make(map[string]*Task),
		taskResults:     make(map[string]*TaskResult),
		taskQueue:       make(chan *Task, 100),
		taskResultQueue: make(chan *TaskResult, 100),
		ctx:             ctx,
		cancel:          cancel,
		completedTasks:  make(map[string]bool),
	}
}

// AddTask 添加任务
func (s *Scheduler) AddTask(task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.ID == "" {
		task.ID = fmt.Sprintf("task-%d", len(s.tasks)+1)
	}
	if task.Name == "" {
		task.Name = task.ID
	}
	if task.Timeout == 0 {
		task.Timeout = 5 * time.Minute
	}
	if task.MaxOutput == 0 {
		task.MaxOutput = 5000
	}
	s.tasks[task.ID] = task
	return nil
}

// checkDependencies 检查任务依赖
func (s *Scheduler) checkDependencies() error {
	// 根据 Dependencies 去检查需要的任务是否存在在队列中
	for _, task := range s.tasks {
		for _, depID := range task.Dependencies {
			if _, exists := s.tasks[depID]; !exists {
				return fmt.Errorf("任务 %s 依赖的任务 %s 不存在", task.ID, depID)
			}
		}
	}
	return nil
}

// taskDispatcher 任务分发器
func (s *Scheduler) taskDispatcher() {
	// 检查依赖
	if err := s.checkDependencies(); err != nil {
		log.Printf("检查依赖失败, %v", err)
		return
	}
	// 没有依赖的任务加入队列
	// 有依赖的任务会在没有依赖的任务完成后执行
	s.mu.Lock()
	for _, task := range s.tasks {
		if len(task.Dependencies) == 0 {
			s.taskQueue <- task
		}
	}
	s.mu.Unlock()
}

// copyAndLog 复制并输出记录
func (s *Scheduler) copyAndLog(dst io.Writer, src io.Reader, prefix, taskName string) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		io.WriteString(dst, line+"\n")
		// 实时日志
		log.Printf("[%s] %s: %s", taskName, prefix, line)
	}
}

// runCommand 执行shell命令
func (s *Scheduler) runCommand(task *Task, output io.Writer) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, task.Timeout)
	defer cancel()

	// 创建命令
	var cmd *exec.Cmd
	if len(task.Args) > 0 {
		cmd = exec.CommandContext(ctx, task.Cmd, task.Args...)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", task.Cmd)
	}

	// 设置工作目录
	if task.WorkDir != "" {
		cmd.Dir = task.WorkDir
	}

	// 设置环境变量
	if len(task.Env) > 0 {
		cmd.Env = append(os.Environ(), cmd.Env...)
	}

	// 设置输出
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return -1, err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return -1, err
	}

	// 并发读取 stdout 和 stderr
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.copyAndLog(output, stdoutPipe, "STDOUT", task.Name)
	}()
	go func() {
		defer wg.Done()
		s.copyAndLog(output, stderrPipe, "STDERR", task.Name)
	}()
	wg.Wait()

	// 等待命令完成
	err = cmd.Wait()
	exitCode := cmd.ProcessState.ExitCode()
	if ctx.Err() == context.DeadlineExceeded {
		return exitCode, fmt.Errorf("任务执行超时(限时: %v)", task.Timeout)
	}
	return exitCode, err
}

// trimOutput 限制输出大小
func (s *Scheduler) trimOutput(output string, maxLines int) string {
	lines := bytes.Split([]byte(output), []byte("\n"))
	if len(lines) <= maxLines {
		return output
	}

	// 保留开头和结尾
	keep := maxLines / 2
	firstPart := lines[:keep]
	lastPart := lines[len(lines)-keep:]

	var result []byte
	result = append(result, bytes.Join(firstPart, []byte("\n"))...)
	result = append(result, []byte("\n... (忽略中间内容) ...\n")...)
	result = append(result, bytes.Join(lastPart, []byte("\n"))...)

	return string(result)
}

// executeTask 执行单个任务
func (s *Scheduler) executeTask(workerID int, task *Task) *TaskResult {
	result := &TaskResult{
		TaskID:     task.ID,
		TaskName:   task.Name,
		Status:     StatusRunning,
		StartTime:  time.Now(),
		RetryCount: 0,
	}

	log.Printf("Worker-%d 开始执行%s: %s", workerID, task.Name, task.Cmd)

	// 执行命令
	var output bytes.Buffer
	var err error
	var exitCode int

	for attempt := 0; attempt <= task.RetryCount; attempt++ {
		if attempt > 0 {
			log.Printf("任务 %s 第 %d 次重试...", task.Name, attempt)
			time.Sleep(task.RetryDelay)
		}

		result.RetryCount = attempt
		output.Reset()
		exitCode, err = s.runCommand(task, &output)

		if err == nil {
			result.Status = StatusSuccess
			break
		}

		if attempt == task.RetryCount {
			result.Status = StatusFailed
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.ExitCode = exitCode
	result.Output = s.trimOutput(output.String(), task.MaxOutput)
	result.Error = err

	return result
}

// worker 工作协程
func (s *Scheduler) worker(id int) {
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case task := <-s.taskQueue:
			result := s.executeTask(id, task)
			s.taskResultQueue <- result
		}
	}
}

// printResult 打印任务结果
func (s *Scheduler) printResult(result *TaskResult) {
	var statusColor *color.Color

	switch result.Status {
	case StatusSuccess:
		statusColor = color.New(color.FgGreen, color.Bold)
	case StatusFailed:
		statusColor = color.New(color.FgRed, color.Bold)
	case StatusCancelled:
		statusColor = color.New(color.FgYellow, color.Bold)
	default:
		statusColor = color.New(color.FgWhite)
	}

	statusColor.Printf("\n任务完成: %s (%s)\n", result.TaskName, result.TaskID)
	fmt.Printf("  状态: %s", result.Status)
	fmt.Printf("  耗时: %v", result.Duration)
	fmt.Printf("  开始: %s", result.StartTime.Format(time.DateTime))
	fmt.Printf("  结束: %s", result.EndTime.Format(time.DateTime))
	fmt.Printf("  退出码: %d", result.ExitCode)
	fmt.Printf("  重试次数: %d", result.RetryCount)

	if result.Error != nil {
		fmt.Printf("  错误: %v\n", result.Error)
	}

	if result.Output != "" {
		fmt.Println("  输出预览:")
		lines := bytes.SplitN([]byte(result.Output), []byte("\n"), 6)
		for i, line := range lines {
			if i >= 5 {
				fmt.Println("    ...(更多输出请查看完整日志)...")
			}
			if len(line) > 0 {
				fmt.Printf("    %s\n", line)
			}
		}
	}
	fmt.Println()
}

// checkDependentTasks 检查依赖任务
func (s *Scheduler) checkDependentTasks() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		// 如果任务已经在队列或已完成则跳过
		if s.completedTasks[task.ID] {
			continue
		}
		// 未进行任务依赖项是否全部满足
		allDepsCompleted := true
		for _, depID := range task.Dependencies {
			if !s.completedTasks[depID] {
				allDepsCompleted = false
				break
			}
		}
		// 如果依赖项项目全部满足，加入队列
		if allDepsCompleted && len(task.Dependencies) > 0 {
			// 标记已调度
			if !s.completedTasks[task.ID] {
				select {
				case s.taskQueue <- task:
					s.completedTasks[task.ID] = true
				default:
					log.Printf("队列任务已满, 任务 %s 等待调度", task.Name)
				}
			}
		}
	}
}

// resultProcessor 处理任务结果
func (s *Scheduler) resultProcessor() {
	for result := range s.taskResultQueue {
		s.mu.Lock()
		s.taskResults[result.TaskID] = result
		s.completedTasks[result.TaskID] = true
		s.mu.Unlock()
		// 打印结果
		s.printResult(result)
		// 检查是否有依赖此任务的任务可以执行
		s.checkDependentTasks()
	}
}

// AddTasks 批量添加任务
func (s *Scheduler) AddTasks(tasks ...*Task) {
	for _, task := range tasks {
		s.AddTask(task)
	}
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("程序已经在运行")
	}
	s.isRunning = true
	s.mu.Unlock()

	// 启动work
	for i := 0; i < s.maxWorkers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	// 启动结果处理器
	go s.resultProcessor()

	// 启动任务调器
	go s.taskDispatcher()

	log.Printf("调度器启动，最大并发数: %d", s.maxWorkers)
	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	log.Println("停止调度器...")
	s.cancel()
	s.wg.Wait()
	close(s.taskQueue)
	close(s.taskResultQueue)
	s.isRunning = false
	log.Println("调度器已停止")
}

// GetResults 获取所有任务结果
func (s *Scheduler) GetResults() map[string]*TaskResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	results := make(map[string]*TaskResult)
	for k, v := range s.taskResults {
		results[k] = v
	}
	return results
}

// PrintSummary 打印汇总报告
func (s *Scheduler) PrintSummary() {
	results := s.GetResults()

	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("任务执行汇总报告")
	fmt.Println(strings.Repeat("-", 60))

	var totalTime time.Duration
	successCount := 0
	failedCount := 0

	for _, result := range results {
		totalTime += result.Duration
		if result.Status == StatusSuccess {
			successCount++
		} else {
			failedCount++
		}
	}

	fmt.Printf("任务总数: %d\n", len(results))
	fmt.Printf("成功: %d\n", successCount)
	fmt.Printf("失败: %d\n", failedCount)
	fmt.Printf("总耗时: %v\n", totalTime)
	fmt.Printf("平均耗时: %v\n", totalTime/time.Duration(len(results)))

	// 打印详细结果表格
	fmt.Println("\n详细结果:")
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%-20s %-15s %-12s %-10s %-30s\n", "任务名称", "状态", "耗时", "退出码", "开始时间")
	fmt.Println(strings.Repeat("-", 100))
	for _, result := range results {
		statusStr := result.Status.String()
		if result.Status == StatusSuccess {
			statusStr = color.GreenString(statusStr)
		} else {
			statusStr = color.RedString(statusStr)
		}

		fmt.Printf("%-20s %-15s %-12v %-10d %-30s\n", result.TaskName, statusStr, result.Duration.Round(time.Millisecond), result.ExitCode, result.StartTime.Format(time.DateTime))
		fmt.Println(strings.Repeat("-", 100))
	}
}

func main() {
	// 创建调度器
	scheduler := NewScheduler(3)

	// 定义任务
	tasks := []*Task{
		{
			ID:         "Test A",
			Name:       "测试脚本A",
			Cmd:        "sh ./shell/test.sh 5 测试脚本A",
			Timeout:    10 * time.Minute,
			RetryDelay: 3 * time.Second,
			RetryCount: 2,
		},
		{
			ID:         "Test B",
			Name:       "测试脚本B",
			Cmd:        "sh ./shell/test.sh 3 测试脚本B",
			Timeout:    10 * time.Minute,
			RetryDelay: 3 * time.Second,
			RetryCount: 2,
		},
		{
			ID:         "Test C",
			Name:       "测试脚本C",
			Cmd:        "sh ./shell/test.sh 2",
			Timeout:    10 * time.Minute,
			RetryDelay: 3 * time.Second,
			RetryCount: 2,
		},
		{
			ID:           "Test D",
			Name:         "测试脚本D",
			Cmd:          "sh ./shell/test.sh 1 测试脚本D",
			Timeout:      10 * time.Minute,
			Dependencies: []string{"Test A", "Test B"},
			RetryDelay:   3 * time.Second,
			RetryCount:   2,
		},
	}

	// 添加任务
	scheduler.AddTasks(tasks...)

	// 启动调度
	if err := scheduler.Start(); err != nil {
		log.Fatal("启动失败:", err)
	}

	// 等待所有任务完成
	fmt.Println("调度器运行中, Ctrl+C 停止")

	// 监听中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 等待完成或者收到中断信号
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n接收到中断信号，正在停止...")
			scheduler.Stop()
			scheduler.PrintSummary()
			return
		case <-ticker.C:
			// 检查是否所有任务都已完成
			results := scheduler.GetResults()
			if len(results) == len(tasks) {
				scheduler.Stop()
				scheduler.PrintSummary()
				return
			}
		}
	}
}
