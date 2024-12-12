package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type Users struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO snippetbox.users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 && strings.Contains(mySQLErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
	}

	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	stmt := `SELECT id, hashed_password FROM snippetbox.users WHERE email=?`

	var userId int
	var hashedPassword string
	err := m.DB.QueryRow(stmt, email).Scan(&userId, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrInvalidCredentials
		} else {
			return -1, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return -1, ErrInvalidCredentials
		}
	}
	return userId, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
