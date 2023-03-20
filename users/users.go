package users

import "errors"

type User struct {
	Username string
	Password string
	Email    string
}
type authUser struct {
	username     string
	passwordHash string
	email        string
}

var authUserDB = map[string]authUser{}

var DefaultUserService userService

type userService struct {
}

func (userService) VerifyUser(user User) bool {
	authUser, ok := authUserDB[user.Username]
	if !ok {
		return false
	}
	err := 0
	if authUser.passwordHash == user.Password {
		err = 0
	}
	return err == 0
}

func (userService) CreateUser(newUser User) error {
	_, ok := authUserDB[newUser.Username]
	if ok {
		return errors.New("user already exists")
	}
	newAuthUser := authUser{
		username:     newUser.Username,
		passwordHash: newUser.Password,
		email:        newUser.Email,
	}
	authUserDB[newAuthUser.username] = newAuthUser
	return nil
}
