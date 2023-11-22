package pkg

import (
	"github.com/google/uuid"
)

func MakeUUID() string {
	return uuid.New().String()
}
