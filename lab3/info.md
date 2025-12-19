# Как отвечать

Покажи ему файл `docker-compose.yml`
```yaml
  rabbitmq: # вот тут у тебя создается образ брокера сообщений. Скажи, что взял Rabbit, а не Kafka, потому что он проще, Kafka ту мач для этого проекта
    image: rabbitmq:3-management-alpine
    container_name: lab3_rabbitmq
    ports:
      - "5672:5672"      # AMQP порт - порт для брокера сообщений
      - "15672:15672"    # Management UI - интерфейс, где ты можешь администрировать брокер
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: password
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
```

Если спросит че такое брокер сообщений - скажи, что это очередь, в которую мы записываем некие сообщения, чтобы они грамотно доставлялись между микросервисами. Здесь используется для уведомлений, чтобы они приходили в прямом эфире

Про остальные сервисы скажи, что тут сделал минимальную структуру и все сделал в main.go
Про БД скажи, что решил сделать максимально просто, и поэтому сделал хранилище памяти как Go-структуру.

Открой `order-service/main.go`

```go
func initRabbitMQ() error { // эта функция создает подключение к брокеру
	var err error
	
	// Повторные попытки подключения. Для надежности
	for i := 0; i < 10; i++ {
		rabbitConn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("Попытка подключения к RabbitMQ (%d/10)...", i+1)
		time.Sleep(2 * time.Second)
	}
	
	if err != nil {
		return fmt.Errorf("не удалось подключиться к RabbitMQ: %w", err)
	}

	rabbitCh, err = rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("не удалось открыть канал: %w", err)
	}

	_, err = rabbitCh.QueueDeclare( // тут объявляешь параметры нашего брокера
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("не удалось объявить очередь: %w", err)
	}

	log.Println("Подключено к RabbitMQ")
	return nil
}
```


А вот эта функция публикует сообщение в брокер событий. В нашем случае это - уведомление
```go
func publishOrderNotification(order *Order) error {
	message, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = rabbitCh.Publish( // вот тут это происходит
        // параметры для отправки бла бла
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate


		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}

	log.Printf("Отправлено уведомление о заказе %s в RabbitMQ", order.OrderID)
	return nil
}
```


А вот тут (это в сервисе уведомлений в main.go) у нас происходит подписка на наш брокер сообщений, чтобы сервис уведомлений получал события из брокера
```go
	msgs, err := ch.Consume( 
        // тут параметры бла бла
		queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Не удалось подписаться на очередь: %v", err)
	}
```

чуть ниже в этом файле есть штука снизу. Она читает все сообщения из очереди и записывает их в Go-шный канал, откуда над ними уже проводится работа
```go
	go func() {
		for msg := range msgs {
			var order Order
			if err := json.Unmarshal(msg.Body, &order); err != nil {
				log.Printf("Ошибка парсинга сообщения: %v", err)
				continue
			}

			processOrder(&order)
		}
	}()

	<-forever
```

# Как запустить
1. Проверь, очистил ли ты предыдущие контейнеры (в папке lab2 пропиши `docker-compose down -v`)
2. Затем в lab3 запусти докер через терминал `docker-compose up -d`
3. Создай 3 терминала под каждый сервис и запусти их через `go run .`
4. В терминале с докером в папке `lab3` запусти тесты `test-simple.bat`
5. Во всех терминалах `Ctrl+C` + закрыть
6. Бла бла очистить докер