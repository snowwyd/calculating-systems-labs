package main

import (
	"log"
	"runtime"
	"time"
)

const serverWaitDuration = 2 * time.Second

func printResults(bestResult Result, elapsed time.Duration) {
	log.Println("\n=== Результаты вычислений ===")
	log.Printf("\nЛучшая найденная последовательность:")
	log.Printf("  Порядок задач: %v\n", bestResult.Sequence)
	log.Printf("  Длительность проекта: %d единиц времени\n", bestResult.Duration)
	log.Printf("  При ограничении ресурса: %d\n", maxResource)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("Ожидание запуска API сервера...")
	time.Sleep(serverWaitDuration)

	// Инициализация компонентов
	apiClient := NewAPIClient(apiURL)
	scheduler := NewScheduler(maxResource)
	optimizer := NewOptimizer(scheduler, numSequences, numWorkers, progressInterval)

	// Создание тестового заказа
	log.Println("Создание тестового заказа...")
	orderID, err := apiClient.CreateTestOrder()
	if err != nil {
		log.Fatalf("Ошибка создания заказа: %v", err)
	}
	log.Printf("Создан заказ ID: %d\n", orderID)

	// Получение данных заказа
	log.Println("Получение данных заказа...")
	order, err := apiClient.GetOrder(orderID)
	if err != nil {
		log.Fatalf("Ошибка получения заказа: %v", err)
	}

	if len(order.Tasks) == 0 {
		log.Fatal("В заказе нет задач для обработки")
	}

	log.Printf("Количество задач: %d\n", len(order.Tasks))

	// Запуск оптимизации
	log.Printf("Начало расчёта %d последовательностей с использованием %d воркеров...\n",
		numSequences, numWorkers)

	startTime := time.Now()
	bestResult := optimizer.FindOptimalSequence(order.Tasks)
	elapsed := time.Since(startTime)

	// Вывод результатов
	printResults(bestResult, elapsed)
}
