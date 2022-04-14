package user

import (
	"bwastartup/helper"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	password_hash, errr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	helper.IfError(errr)

	user := User{
		Name:         input.Name,
		Occupation:   input.Occupation,
		Email:        input.Email,
		PasswordHash: string(password_hash),
		Role:         "user",
	}

	user, err := s.repository.Save(user)
	helper.IfError(err)
	return user, nil
}
