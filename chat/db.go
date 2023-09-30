package main

import (
	"fmt"

	"github.com/tim117/chatserver/chat"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect connects the gorm.io/gorm client to the chat.db sqlite database.
func Connect() *gorm.DB {
	fmt.Println("../chat.db")
	database, err := gorm.Open(sqlite.Open("../chat.db"))

	if err != nil {
		panic("Unable to connect to the database.")
	}

	database.AutoMigrate(&chat.User{}, &chat.Chat{}, &chat.Message{})
	return database
}

// SeedData seeds the gorm.io/gorm DB wiht chat and user data and panics if an
// error occurs.
func SeedData(db *gorm.DB) {
	userRepo := chat.UserRepository(db)

	const systemUsername = "system"
	system, err := userRepo.GetByUsername(systemUsername)
	if err != nil {
		fmt.Printf("Failed while trying to find user with username %s.\n", systemUsername)
		panic(err.Error())
	}
	if system == nil {
		fmt.Printf("No user with the username %s found. Creating a new user.\n", systemUsername)
		system, err = userRepo.Save(chat.User{Username: systemUsername})
		if err != nil {
			fmt.Printf("Failed while trying to create user with username %s.\n", systemUsername)
			panic(err.Error())
		}
	}

	chatRepo := chat.ChatRepository(db)
	publicChatTitle := "public"
	chats, err := chatRepo.FindByTitle(publicChatTitle)
	if err != nil {
		fmt.Printf("Failed while trying to find chat room with title %s.\n", publicChatTitle)
		panic(err.Error())
	}
	public := First(*chats, func(c chat.Chat) bool {
		return c.Title == publicChatTitle
	})
	if public == nil {
		fmt.Printf("No chat room with the title %s found. Creating a new chat room.\n", publicChatTitle)
		public, err = chatRepo.Save(chat.Chat{Title: publicChatTitle}, system.ID)
		if err != nil {
			fmt.Printf("Failed while trying to create chat room with title %s.\n", publicChatTitle)
			panic(err.Error())
		}
	}
	if len(public.Messages) == 0 {
		fmt.Printf("No messages in %s. Sending a first message.\n", publicChatTitle)
		_, err = chatRepo.SaveMessage(public.ID, "Welcome! This is the public chat room. Everyone is invited 😇", system.ID)
		if err != nil {
			fmt.Printf("Failed to send a message in %s.\n", publicChatTitle)
			panic(err.Error())
		}
	} else {
		fmt.Printf("There are already %d messages in %s. First message is %s.\n", len(public.Messages), publicChatTitle, public.Messages[0].Text)
	}
	fmt.Println()
}

func First[T any](list []T, evaluate func(T) bool) *T {
	for _, item := range list {
		if evaluate(item) {
			return &item
		}
	}
	return nil
}
