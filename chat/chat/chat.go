// Package chat contains utilities for working with the chat database.
package chat

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Chat represents a chat room.
type Chat struct {
	ID        string `gorm:"primaryKey"`
	Title     string
	Messages  []Message
	Users     []*User `gorm:"many2many:user_languages;"`
	Creator   User
	CreatorId string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Chat.BeforeCreate sets the ID of the Chat that is being created as a new
// UUID.
func (c *Chat) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return
}

// Message represents a message in a Chat room.
type Message struct {
	ID        string `gorm:"primaryKey"`
	Author    User
	Text      string
	Chat      Chat
	ChatId    string
	AuthorId  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Message.BeforeCreate sets the ID of the Message that is being created as a
// new UUID.
func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.NewString()
	return
}

type ChatRepo interface {
	Save(Chat, string) (*Chat, error)
	GetById(string) (*Chat, error)
	SaveMessage(string, string, string) (*Message, error)
	FindByTitle(string) (*[]Chat, error)
}

type chatRepository struct {
	ChatRepo
	DB *gorm.DB
}

// ChatRepository creates a repository for working with chat rooms in the DB.
func ChatRepository(db *gorm.DB) ChatRepo {
	return &chatRepository{DB: db}
}

// Save saves a chat room to the DB.
func (cr *chatRepository) Save(chat Chat, creatorId string) (*Chat, error) {
	chat.CreatorId = creatorId
	if result := cr.DB.Create(&chat); result.Error != nil {
		return nil, result.Error
	}
	return &chat, nil
}

// GetById gets a chat from the DB by its ID field.
func (cr *chatRepository) GetById(id string) (*Chat, error) {
	chat := &Chat{}
	if result := cr.DB.Model(&Chat{}).
		Preload("Users").
		Preload("Messages").
		First(chat, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}
	return chat, nil
}

// SaveMessage saves a message with the given text by a specified user to the
// specified chat room.
func (cr *chatRepository) SaveMessage(chatId string, text string, userId string) (*Message, error) {
	message := &Message{Text: text, ChatId: chatId, AuthorId: userId}
	if result := cr.DB.Create(&message); result.Error != nil {
		return nil, result.Error
	}
	return message, nil
}

func (cr *chatRepository) FindByTitle(title string) (*[]Chat, error) {
	chats := &[]Chat{}
	if result := cr.DB.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title)).Find(chats); result.Error != nil {
		return nil, result.Error
	}
	return chats, nil
}
