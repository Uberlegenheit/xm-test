package services

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"xm-task/dmodels"
	"xm-task/smodels"
)

func (s *ServiceFacade) SignInOrRegister(user smodels.User) (bool, error) {
	usr, err := s.dao.GetUserByEmail(user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password+os.Getenv("PASSWORD_SALT")), bcrypt.DefaultCost)
			if err != nil {
				return false, fmt.Errorf("pass bcrypt.GenerateFromPassword: %v", err)
			}

			usr, err = s.dao.CreateUser(dmodels.User{
				ID:       uuid.NewV4(),
				Email:    user.Email,
				Password: string(hashedPassword),
			})
			if err != nil {
				return false, fmt.Errorf("dao.CreateUser: %v", err)
			}

			return true, nil
		}
		return false, fmt.Errorf("dao.GetUserByEmail: %v", err)
	}

	return bcrypt.CompareHashAndPassword([]byte(usr.Password+os.Getenv("PASSWORD_SALT")), []byte(user.Password)) == nil, nil
}

func (s *ServiceFacade) GetUserByEmail(email string) (dmodels.User, error) {
	user, err := s.dao.GetUserByEmail(email)
	if err != nil {
		return dmodels.User{}, fmt.Errorf("dao.GetUserByEmail: %v", err)
	}

	return user, err
}
