package auth

import (
	"time"
)

var (
	ExtraLatency time.Duration
)

type Service interface {
	Authenticate(user User) error
	Init()
}

func NewAuthService() Service {
	return &authService{}
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type authService struct {
}

func (u *authService) Init() {
}

func (u *authService) Authenticate(user User) error {
	time.Sleep(ExtraLatency)
	return nil
}
