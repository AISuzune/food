package dao

import (
	"main/app/internal/dao/user"
)

var insUser = user.Group{}

func User() *user.Group {
	return &insUser
}
