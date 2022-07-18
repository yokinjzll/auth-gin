package mysql

import (
	"database/sql"
	"errors"
	"gin-chat/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, password string) (int, error) {
	hashpass, err := HashPassword(password)
	if err != nil {
		return 0, err
	}
	stmp := `INSERT INTO user (username, password)
	VALUES(?, ?)`

	result, err := m.DB.Exec(stmp, username, hashpass)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *UserModel) Auth(username, password string) (int, error) {

	stmt := `SELECT id, username, password FROM user WHERE username = ?`
	row := m.DB.QueryRow(stmt, username)
	u := &models.User{}
	err := row.Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		} else {
			return 0, err
		}
	}

	checkhash := CheckPasswordHash(password, u.Password)
	if !checkhash {
		return 0, err
	}
	return u.ID, nil

}
