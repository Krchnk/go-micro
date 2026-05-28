package users

import "errors"

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	Create(name, email string) User
	Update(id int64, name, email string) (User, error)
	Delete(id int64) error
	List() []User
}
