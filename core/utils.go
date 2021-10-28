package core

import (
	"github.com/google/uuid"
)

func GetUUID() string {
	u2, _ := uuid.NewUUID()
	return u2.String()
}
