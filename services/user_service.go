package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"mall-seckill/datamodels"
	"mall-seckill/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (*datamodels.User, bool)
	AddUser(user *datamodels.User) (int64, error)
}

type UserService struct {
	UserRepository repositories.IUser
}

func NewUserService(repository repositories.IUser) IUserService {
	return &UserService{repository}
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (*datamodels.User, bool) {
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return &datamodels.User{}, false
	}
	isOk, _ := ValidatePassword(pwd, user.HashPassword)
	if !isOk {
		return &datamodels.User{}, false
	}
	return user, true
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return userId, err
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误")
	}
	return true, nil
}
