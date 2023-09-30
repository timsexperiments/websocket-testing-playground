package cmd

import "gorm.io/gorm"

type LocalStorage interface {
	GetUsername() string
}

type localStorage struct {
	gorm.Model
	Username string
}

func GetUsername() {

}

func CraeteLocalStorage() {}
