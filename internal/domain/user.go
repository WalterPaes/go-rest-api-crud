package domain

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func (u *User) EncryptPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
