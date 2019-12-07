package useCase

import (
	"github.com/omise/go-tamboon/src/entities"
	"github.com/omise/go-tamboon/src/repositries"
)

type UserChannel interface {
	GetUserChannel(user chan entities.UserInfo) error
}

type UserStruct struct {
	Repo repositries.FileSys
}

func InitUserUsecase(user UserStruct) UserStruct {
	return user
}

func (userInst *UserStruct) StartToPushDataToUserChannel(fileName string) (error) {
	return userInst.Repo.ReadFile(fileName)
}
