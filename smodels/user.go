package smodels

import "regexp"

const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

type User struct {
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Validate() bool {
	regex := regexp.MustCompile(emailRegex)
	return regex.MatchString(u.Email)
}
