package repo

import (
	"log/slog"
)

func GetAllUsers() ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT id FROM users ")
	if err != nil {
		slog.Error(err.Error())
		return []User{}, err
	}

	return users, nil
}

func GetUser(login string) (User, error) {
	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE login=? ", login)
	if err != nil {
		slog.Error(err.Error())
		return User{}, err
	}

	return user, nil
}

func AddUser(user *User) error {
	_, err := db.NamedExec("INSERT INTO users (login, email, hash, isadmin) VALUES (:login, :email, :hash, :isadmin)", &user)
	if err != nil {
		slog.Error(err.Error())
	}

	return err
}
