package models

import "database/sql"

type User struct {
	ID    string `json:id`
	Name  string `json:name`
	Email string `json:email`
}

type UserModel struct {
	DB *sql.DB
}

// func (m *UserModel) Get(id int) (User, error) {

// }

// func (u *UserModel) Delete(id int) error {

// }

func (m *UserModel) Insert(u User) (int, error) {
	stmt := `INSERT INTO users (name, email) 
			VALUES ($1, $2) RETURNING id;`
	id := -1
	err := m.DB.QueryRow(stmt, u.Name, u.Email).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// func (m *UserModel) I(id int) (User, error) {
// }
