package processor

import (
	"errors"
	"log"
	"strings"

	"github.com/deeramster/goka_sprint2/pkg/censor"
	"github.com/deeramster/goka_sprint2/pkg/models"
	"github.com/deeramster/goka_sprint2/pkg/storage"
)

type MessageProcessor struct {
	storage storage.Storage
	censor  censor.Service
}

func NewMessageProcessor(storage storage.Storage, censor censor.Service) *MessageProcessor {
	return &MessageProcessor{
		storage: storage,
		censor:  censor,
	}
}

func (mp *MessageProcessor) ProcessMessage(msg *models.Message) *models.Message {
	// Строгая проверка входных данных
	if msg == nil || msg.From == "" || msg.To == "" {
		log.Printf("Некорректное сообщение: %+v", msg)
		return nil
	}

	// Загружаем список заблокированных для получателя
	blockedUsers, err := mp.storage.LoadBlockedUsers(msg.To)
	if err != nil {
		log.Printf("Критическая ошибка загрузки списка заблокированных для %s: %v", msg.To, err)
		// В случае ошибки - блокируем сообщение
		return nil
	}

	// Точная проверка блокировки
	for _, blockedUser := range blockedUsers.Users {
		if strings.EqualFold(blockedUser, msg.From) {
			log.Printf("БЛОКИРОВКА: Сообщение от %s полностью заблокировано для %s", msg.From, msg.To)
			return nil
		}
	}

	// Цензура сообщения
	censoredContent := mp.censor.CensorMessage(msg.Content)

	return &models.Message{
		From:    msg.From,
		To:      msg.To,
		Content: censoredContent,
	}
}

func (mp *MessageProcessor) HandleBlockCommand(cmd *models.BlockCommand) error {
	// Проверяем корректность команды
	if cmd.User == "" || cmd.BlockUser == "" {
		return errors.New("некорректная команда блокировки: пустой пользователь")
	}

	// Проверяем, не пытается ли пользователь заблокировать сам себя
	if cmd.User == cmd.BlockUser {
		return errors.New("нельзя заблокировать самого себя")
	}

	// Загружаем текущий список заблокированных
	blockedUsers, err := mp.storage.LoadBlockedUsers(cmd.User)
	if err != nil {
		log.Printf("Ошибка загрузки списка заблокированных для %s: %v", cmd.User, err)
		return err
	}

	// Проверяем, не заблокирован ли пользователь уже
	for _, blockedUser := range blockedUsers.Users {
		if blockedUser == cmd.BlockUser {
			// Пользователь уже заблокирован
			log.Printf("%s уже заблокирован пользователем %s", cmd.BlockUser, cmd.User)
			return nil
		}
	}

	// Добавляем пользователя в список заблокированных
	blockedUsers.Users = append(blockedUsers.Users, cmd.BlockUser)

	// Сохраняем обновленный список
	if err := mp.storage.SaveBlockedUsers(cmd.User, blockedUsers); err != nil {
		log.Printf("Ошибка сохранения списка заблокированных для %s: %v", cmd.User, err)
		return err
	}

	log.Printf("Пользователь %s заблокирован пользователем %s", cmd.BlockUser, cmd.User)
	return nil
}
