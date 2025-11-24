package main

const maxResource = 10

// Scheduler отвечает за расчёт длительности проекта
type Scheduler struct {
	maxResource int
}

// NewScheduler создаёт новый планировщик
func NewScheduler(maxResource int) *Scheduler {
	return &Scheduler{maxResource: maxResource}
}

// CalculateProjectDuration рассчитывает длительность проекта для заданной последовательности
func (s *Scheduler) CalculateProjectDuration(tasks []Task, sequence []int) int {
	if len(sequence) == 0 {
		return 0
	}

	executions := make(map[int]TaskExecution)
	resourceUsage := make(map[int]int)

	for _, taskID := range sequence {
		task, found := s.findTaskByID(tasks, taskID)
		if !found {
			continue
		}

		earliestStart := s.calculateEarliestStart(task, executions)
		startTime := s.findAvailableStartTime(earliestStart, task.Duration, task.Resource, resourceUsage)
		s.reserveResources(startTime, task.Duration, task.Resource, resourceUsage)

		executions[task.ID] = TaskExecution{
			StartTime: startTime,
			EndTime:   startTime + task.Duration,
			Resource:  task.Resource,
		}
	}

	return s.calculateMaxEndTime(executions)
}

func (s *Scheduler) findTaskByID(tasks []Task, taskID int) (Task, bool) {
	for _, t := range tasks {
		if t.ID == taskID {
			return t, true
		}
	}
	return Task{}, false
}

func (s *Scheduler) calculateEarliestStart(task Task, executions map[int]TaskExecution) int {
	earliestStart := 0
	for _, predID := range task.Pred {
		if exec, exists := executions[predID]; exists {
			if exec.EndTime > earliestStart {
				earliestStart = exec.EndTime
			}
		}
	}
	return earliestStart
}

func (s *Scheduler) canStartAtTime(startTime, duration, resource int, resourceUsage map[int]int) bool {
	for t := startTime; t < startTime+duration; t++ {
		if resourceUsage[t]+resource > s.maxResource {
			return false
		}
	}
	return true
}

func (s *Scheduler) findAvailableStartTime(earliestStart, duration, resource int, resourceUsage map[int]int) int {
	startTime := earliestStart
	for !s.canStartAtTime(startTime, duration, resource, resourceUsage) {
		startTime++
	}
	return startTime
}

func (s *Scheduler) reserveResources(startTime, duration, resource int, resourceUsage map[int]int) {
	for t := startTime; t < startTime+duration; t++ {
		resourceUsage[t] += resource
	}
}

func (s *Scheduler) calculateMaxEndTime(executions map[int]TaskExecution) int {
	maxEndTime := 0
	for _, exec := range executions {
		if exec.EndTime > maxEndTime {
			maxEndTime = exec.EndTime
		}
	}
	return maxEndTime
}
