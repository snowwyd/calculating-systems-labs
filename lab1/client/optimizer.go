package main

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	numSequences     = 1000000
	numWorkers       = 8
	progressInterval = 100000
)

// Optimizer выполняет поиск оптимальной последовательности выполнения задач
type Optimizer struct {
	scheduler        *Scheduler
	numSequences     int
	numWorkers       int
	progressInterval int
}

// NewOptimizer создаёт новый оптимизатор
func NewOptimizer(scheduler *Scheduler, numSequences, numWorkers, progressInterval int) *Optimizer {
	return &Optimizer{
		scheduler:        scheduler,
		numSequences:     numSequences,
		numWorkers:       numWorkers,
		progressInterval: progressInterval,
	}
}

// FindOptimalSequence ищет оптимальную последовательность выполнения задач
func (o *Optimizer) FindOptimalSequence(tasks []Task) Result {
	taskIDs := extractTaskIDs(tasks)
	return o.runParallelComputation(tasks, taskIDs)
}

func (o *Optimizer) runParallelComputation(tasks []Task, taskIDs []int) Result {
	jobs := make(chan int, o.numSequences)
	results := make(chan Result, o.numSequences)

	var wg sync.WaitGroup
	for w := 0; w < o.numWorkers; w++ {
		wg.Add(1)
		go o.worker(tasks, taskIDs, jobs, results, &wg)
	}

	go func() {
		for i := 0; i < o.numSequences; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return o.collectResults(results)
}

func (o *Optimizer) worker(tasks []Task, taskIDs []int, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(rand.Int())))

	for range jobs {
		sequence := make([]int, len(taskIDs))
		copy(sequence, taskIDs)
		shuffleSequence(sequence, rng)

		duration := o.scheduler.CalculateProjectDuration(tasks, sequence)
		results <- Result{
			Sequence: sequence,
			Duration: duration,
		}
	}
}

func (o *Optimizer) collectResults(results <-chan Result) Result {
	bestResult := Result{Duration: math.MaxInt}
	count := 0

	for result := range results {
		count++
		if result.Duration < bestResult.Duration {
			bestResult = result
		}

		if count%o.progressInterval == 0 {
			progress := float64(count) / float64(o.numSequences) * 100
			log.Printf("Обработано %d последовательностей (%.1f%%)...", count, progress)
		}
	}

	return bestResult
}

// Вспомогательные функции

func extractTaskIDs(tasks []Task) []int {
	taskIDs := make([]int, len(tasks))
	for i, task := range tasks {
		taskIDs[i] = task.ID
	}
	return taskIDs
}

func shuffleSequence(sequence []int, rng *rand.Rand) {
	for i := len(sequence) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		sequence[i], sequence[j] = sequence[j], sequence[i]
	}
}
