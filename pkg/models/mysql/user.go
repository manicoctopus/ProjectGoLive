package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"ProjectGoLive/pkg/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(c *models.User) (int, error) {
	stmt := `INSERT INTO 
		user (UserName, UserEmail, HashedPassword, UserContact, IsBOwner, IsVerified, Created) 
		VALUES (?, ?, ?, ?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, c.UserName, c.UserEmail, c.HashedPassword, c.UserContact, c.IsBOwner, c.IsVerified)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, c.UserEmail) {
				return -1, models.ErrDuplicateEmail
			}
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (m *UserModel) Update(c *models.User) error {
	stmt := `UPDATE user 
			 SET UserName=?, UserEmail=?, HashedPassword=?, UserContact=?, IsBOwner=?, IsVerified=?
			 WHERE UserID = ?`

	_, err := m.DB.Exec(stmt, c.UserName, c.UserEmail, c.HashedPassword, c.UserContact, c.IsBOwner, c.IsVerified, c.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Delete(id uint32) error {
	stmt := `DELETE FROM users WHERE 
	id = ?`
	results, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	} else if rows, _ := results.RowsAffected(); rows == 0 {
		return models.ErrNoRecord
	} else {
		return err
	}
}

func (m *UserModel) Retrieve(id uint32) (*models.User, error) {
	stmt := `SELECT 
				UserID, UserName, UserEmail, HashedPassword, UserContact, IsBOwner, IsVerified, Created
			FROM user
			WHERE UserID = ?`

	u := &models.User{}
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&u.UserID,
		&u.UserName,
		&u.UserEmail,
		&u.HashedPassword,
		&u.UserContact,
		&u.IsBOwner,
		&u.IsVerified,
		&u.Created,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return u, nil
}

func (m *UserModel) RetrieveAll() ([]*models.User, error) {
	stmt :=
		`SELECT 
			UserID, UserName, UserEmail, HashedPassword, UserContact, IsBOwner, IsVerified, Created
		FROM user`

	users := []*models.User{}

	rows, _ := m.DB.Query(stmt)

	for rows.Next() {
		u := &models.User{}
		err := rows.Scan(
			&u.UserID,
			&u.UserName,
			&u.UserEmail,
			&u.HashedPassword,
			&u.UserContact,
			&u.IsBOwner,
			&u.IsVerified,
			&u.Created,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (m *UserModel) AuthenticateUser(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	results := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := results.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, err
}
