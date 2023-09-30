package chat

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is a representation of a User of the chat server.
type User struct {
	ID           string    `gorm:"primaryKey"`
	Username     string    `gorm:"unique"`
	Chats        []Chat    `gorm:"many2many:user_chats"`
	CreatedChats []Chat    `gorm:"foreignKey:CreatorId"`
	Messages     []Message `gorm:"foreignKey:AuthorId"`
}

// User.BeforeCreate adds a UUID in the User.ID field before creating it.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}

type UserRepo interface {
	Save(user User) (*User, error)
	GetByUsername(username string) (*User, error)
}

type userRepository struct {
	UserRepo
	DB *gorm.DB
}

// UserRepository creates a repository for working with users in the DB.
func UserRepository(db *gorm.DB) UserRepo {
	return &userRepository{DB: db}
}

// Save saves a User to the DB.
func (ur *userRepository) Save(user User) (*User, error) {
	if result := ur.DB.Create(&user); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByUsername gets a user based on their username.
func (ur *userRepository) GetByUsername(username string) (*User, error) {
	user := &User{}
	if result := ur.DB.Model(&User{}).Preload("Chats").First(user, "username = ?", username); result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
