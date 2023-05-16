package database

import (
	"bookstore/model"
	"time"
)

func resetUserSecrets(user *model.User) {
	resetUserPassword(user)
	user.Role = ""

	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}
}

func resetUserPassword(user *model.User) {
	user.Password = ""
}

