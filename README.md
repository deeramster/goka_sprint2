# Message Stream Processor

## Описание проекта

Система обработки потоков сообщений с функциями:
- Блокировки нежелательных пользователей
- Цензуры сообщений

### Компоненты системы

1. **Message Processor**:
   - Обработка входящих сообщений
   - Фильтрация сообщений от заблокированных пользователей
   - Цензура контента

2. **Storage**:
   - Файловое хранение списков заблокированных пользователей
   - Персистентность данных между сессиями

3. **Kafka Infrastructure**:
   - Топик `messages` для входящих сообщений
   - Топик `filtered_messages` для обработанных сообщений
   - Топик `blocked_users` для команд блокировки

### Архитектура

Система использует микросервисную архитектуру с распределенной обработкой сообщений через Apache Kafka и библиотеку Goka.

## Требования

- Docker
- Docker Compose
- Go 1.23

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone git@github.com:deeramster/goka_sprint2.git
cd goka_sprint2
```

2. Запустите инфраструктуру:
```bash
docker-compose up -d
```

3. Соберите и запустите процессор сообщений:
```bash
cd cmd/message_processor
go build main.go
```

## Тестирование

### Тестовые сценарии

1. **Блокировка пользователя**
   - Отправьте команду блокировки в топик `blocked_users`
   ```json
   {
     "user": "alice",
     "block": "bob"
   }
   ```

2. **Отправка сообщения**
   - Отправьте сообщение в топик `messages`
   ```json
   {
     "from": "bob", 
     "to": "alice", 
     "content": "Hello, bad words are not allowed!"
   }
   ```

3. **Цензура сообщений**
   - Проверьте, что сообщения с запрещенными словами цензурируются

### Инструкция по тестированию

1. Используйте Kafka UI (http://localhost:8080) для мониторинга топиков
2. Отправляйте тестовые сообщения через Kafka CLI или UI
3. Проверяйте результаты в топике `filtered_messages`

## Примеры команд

### Блокировка пользователя
```bash
kafka-console-producer --bootstrap-server localhost:9094 \
  --topic blocked_users 
```

### Отправка сообщения
```bash
kafka-console-producer --bootstrap-server localhost:9094 \
  --topic messages 
```

## Конфигурация

- Список запрещенных слов: `["bad", "words", "censored"]`
- Брокеры Kafka: `localhost:9094,localhost:9095,localhost:9096`


