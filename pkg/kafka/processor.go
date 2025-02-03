package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/deeramster/goka_sprint2/pkg/models"
	"github.com/deeramster/goka_sprint2/pkg/processor"
	"github.com/lovoo/goka"
)

type KafkaProcessor struct {
	msgProcessor *processor.MessageProcessor
	brokers      []string
}

func NewKafkaProcessor(brokers []string, msgProcessor *processor.MessageProcessor) *KafkaProcessor {
	return &KafkaProcessor{
		msgProcessor: msgProcessor,
		brokers:      brokers,
	}
}

type messageCodec struct{}

func (c *messageCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *messageCodec) Decode(data []byte) (interface{}, error) {
	var m models.Message
	err := json.Unmarshal(data, &m)
	return &m, err
}

type blockCommandCodec struct{}

func (c *blockCommandCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *blockCommandCodec) Decode(data []byte) (interface{}, error) {
	var cmd models.BlockCommand
	err := json.Unmarshal(data, &cmd)
	return &cmd, err
}

func (kp *KafkaProcessor) Run(ctx context.Context) error {
	group := goka.DefineGroup("message-processor",
		goka.Input("messages", new(messageCodec), kp.handleMessage),
		goka.Input("blocked_users", new(blockCommandCodec), kp.handleBlockCommand),
		goka.Output("filtered_messages", new(messageCodec)),
	)

	proc, err := goka.NewProcessor(kp.brokers, group)
	if err != nil {
		return err
	}

	return proc.Run(ctx)
}

func (kp *KafkaProcessor) handleMessage(ctx goka.Context, msg interface{}) {
	message, ok := msg.(*models.Message)
	if !ok {
		log.Printf("Invalid message format: %v", msg)
		return
	}

	processed := kp.msgProcessor.ProcessMessage(message)
	if processed != nil {
		ctx.Emit("filtered_messages", message.To, processed)
	}
}

func (kp *KafkaProcessor) handleBlockCommand(ctx goka.Context, msg interface{}) {
	cmd, ok := msg.(*models.BlockCommand)
	if !ok {
		log.Printf("Некорректный формат команды блокировки: %v", msg)
		return
	}

	// Логируем входящую команду блокировки
	log.Printf("Получена команда блокировки: %+v", cmd)

	if err := kp.msgProcessor.HandleBlockCommand(cmd); err != nil {
		log.Printf("Ошибка обработки команды блокировки: %v", err)
		// В зависимости от требований можно реализовать отправку сообщения об ошибке
	}
}
