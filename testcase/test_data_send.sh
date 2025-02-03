#!/bin/bash

# Отправка команд блокировки
docker exec -it 9665b28da839 kafka-console-producer.sh \
    --topic blocked_users \
    --bootstrap-server localhost:9094 << EOF
{"user": "alice", "block": "bob"}
{"user": "charlie", "block": "dave"}
EOF

# Отправка тестовых сообщений
docker exec -it 9665b28da839 kafka-console-producer.sh \
    --topic messages \
    --bootstrap-server localhost:9094 << EOF
{"from": "bob", "to": "alice", "content": "Привет, это тестовое сообщение которое должно быть заблокировано"}
{"from": "dave", "to": "charlie", "content": "Это сообщение содержит bad words и должно быть отцензурировано"}
{"from": "eve", "to": "frank", "content": "Нормальное сообщение без цензуры"}
EOF

# Чтение результатов из filtered_messages
docker exec -it 9665b28da839 kafka-console-consumer.sh \
    --topic filtered_messages \
    --bootstrap-server localhost:9094 \
    --from-beginning \
    --max-messages 10