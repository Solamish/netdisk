package model

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

func (user *User) SignUp() error{
	var usr User
	DB.Where("username = ?", user.Username).First(&usr)
	if usr.Username == user.Username {
		err := errors.New("username exited")
		return err
	}
	return DB.Create(user).Error
}

func SignIn(username string, password string) bool{
	var user User
	DB.Where("username = ?", username).First(&user)
	if user.Password == password {
		return true
	}
	return false
}