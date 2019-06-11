package project

import (
	// 
)

// Service contains services to manage data in the repository.
type Service interface {
	UserService

	Close() error
}

type UserService interface {
	FetchIndexUser(*PageRequest) ([]*User, int, error)
	FetchShowUser(string) (*User, error)
	FetchStoreUser(*User) (*User, error)
	FetchUpdateUser(string, *User) (*User, error)
	FetchDestroyUser(string) error
}
