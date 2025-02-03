package censor

import (
	"strings"
	"sync"
)

type Service interface {
	UpdateBannedWords(words []string)
	CensorMessage(content string) string
}

type censor struct {
	bannedWords []string
	mutex       sync.RWMutex
}

func NewCensor(initialBannedWords []string) Service {
	return &censor{
		bannedWords: initialBannedWords,
	}
}

func (c *censor) UpdateBannedWords(words []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.bannedWords = words
}

func (c *censor) CensorMessage(content string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	censoredContent := content
	for _, word := range c.bannedWords {
		replacement := strings.Repeat("*", len(word))
		censoredContent = strings.ReplaceAll(censoredContent, word, replacement)
	}
	return censoredContent
}
