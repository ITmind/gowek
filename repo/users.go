package repo

import (
	"log/slog"
)

func GetAllUsers() ([]User, error) {
	var users []User
	err := DB.Select(&users, "SELECT id FROM users ")
	if err != nil {
		slog.Error(err.Error())
		return []User{}, err
	}

	return users, nil
}

func GetUser(login string) (User, error) {
	var user User
	err := DB.Get(&user, "SELECT id FROM users WHERE login=? ", login)
	if err != nil {
		slog.Error(err.Error())
		return User{}, err
	}

	return user, nil
}

func AddUser(user *User) error {
	_, err := DB.NamedExec("INSERT INTO users (login, email, hash) VALUES (:login, :email, :hash)", &user)
	if err != nil {
		slog.Error(err.Error())
	}

	return err
}
