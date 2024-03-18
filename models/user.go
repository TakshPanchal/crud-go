package models

import "database/sql"

type User struct {
	ID    string
	Name  string
	Email string
}

type UserModel struct {
	DB *sql.DB
}

// func (m *UserModel) Get(id int) (User, error) {

// }

// func (u *UserModel) Delete(id int) error {

// }

// func (m *UserModel) Create(u User) (User, error) {

// }

// func (m *UserModel) I(id int) (User, error) {
// }
