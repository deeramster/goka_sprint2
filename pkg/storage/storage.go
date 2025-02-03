package storage

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"sync"

	"github.com/deeramster/goka_sprint2/pkg/models"
)

type Storage interface {
	LoadBlockedUsers(user string) (models.BlockedUsers, error)
	SaveBlockedUsers(user string, blocked models.BlockedUsers) error
}

type FileStorage struct {
	mutex sync.RWMutex
	dir   string
}

func NewFileStorage(dir string) (*FileStorage, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Ошибка создания директории %s: %v", dir, err)
		return nil, err
	}
	return &FileStorage{dir: dir}, nil
}

func (fs *FileStorage) LoadBlockedUsers(user string) (models.BlockedUsers, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	filePath := fs.filePath(user)

	// Если файл не существует, возвращаем пустой список
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return models.BlockedUsers{Users: []string{}}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Ошибка чтения файла %s: %v", filePath, err)
		return models.BlockedUsers{}, err
	}

	var blocked models.BlockedUsers
	if err := json.Unmarshal(data, &blocked); err != nil {
		log.Printf("Ошибка десериализации данных из %s: %v", filePath, err)
		return models.BlockedUsers{}, err
	}

	return blocked, nil
}

func (fs *FileStorage) SaveBlockedUsers(user string, blocked models.BlockedUsers) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	data, err := json.MarshalIndent(blocked, "", "  ")
	if err != nil {
		log.Printf("Ошибка сериализации списка заблокированных для %s: %v", user, err)
		return err
	}

	filePath := fs.filePath(user)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Printf("Ошибка записи файла %s: %v", filePath, err)
		return err
	}

	log.Printf("Список заблокированных для %s обновлен", user)
	return nil
}

func (fs *FileStorage) filePath(user string) string {
	return path.Join(fs.dir, user+"_blocked.json")
}
