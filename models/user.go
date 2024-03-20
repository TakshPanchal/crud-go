package models

import (
	"database/sql"
	"errors"
	"strconv"
)

type User struct {
	ID    int    `json:id`
	Name  string `json:name`
	Email string `json:email`
}

type UserModel struct {
	db *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{db: db}
}

func (m *UserModel) Get(id int) (*User, error) {
	stmt := `SELECT id, name, email FROM users 
			WHERE id=$1`
	var user User
	err := m.db.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Delete(id int) error {
	stmt := `DELETE FROM users
			WHERE id=$1;`
	_, err := m.db.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Insert(u *User) (int, error) {
	stmt := `INSERT INTO users (name, email) 
			VALUES ($1, $2) RETURNING id;`
	id := -1
	err := m.db.QueryRow(stmt, u.Name, u.Email).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (m *UserModel) UpdateFields(id int, updatedUser *User) error {
	stmt := "UPDATE users SET "
	args := []interface{}{id}

	if len(updatedUser.Name) != 0 {
		stmt += "name = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, updatedUser.Name)
	}

	if len(updatedUser.Email) != 0 {
		stmt += "email = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, updatedUser.Email)
	}

	if len(args) == 1 {
		return errors.New("no fields for update")
	}

	stmt = stmt[:len(stmt)-2] // Remove the trailing comma and space
	stmt += " WHERE id = $1"

	_, err := m.db.Exec(stmt, args...)

	return err
}
