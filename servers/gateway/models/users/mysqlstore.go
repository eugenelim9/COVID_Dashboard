package users

import (
	"database/sql"
	"errors"
	"time"

	// import for sql driver
	_ "github.com/go-sql-driver/mysql"
)

// ErrDeletingUser returned when theres an error deleting a user
var ErrDeletingUser = errors.New("error deleting user")

//MySQLStore represents a MySql store
type MySQLStore struct {
	Db *sql.DB
}

// GetByID returns User with given ID
func (ms *MySQLStore) GetByID(id int64) (*User, error) {
	result := &User{}
	rows, err := ms.Db.Query("SELECT * FROM users WHERE id=?", id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Email, &result.PassHash, &result.UserName, &result.FirstName, &result.LastName, &result.PhotoURL)
		if err != nil {
			return nil, err
		}
		// take only first row
		break
	}
	return result, nil
}

// Insert inserts the user into the database, and returns
// the newly-inserted User, complete with the DBMS-assigned ID
func (ms *MySQLStore) Insert(user *User) (*User, error) {
	insq := "insert into users(email, pass_hash, username, first_name, last_name, photo_url) values (?,?,?,?,?,?)"
	res, execErr := ms.Db.Exec(insq, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL)
	if execErr != nil {
		return nil, errors.New("could not insert new user")
	}
	id, insertIDErr := res.LastInsertId()
	if insertIDErr != nil {
		return nil, errors.New("could not insert new user")
	}
	user.ID = id
	return user, nil
}

//Log logs user logins into our database table
func (ms *MySQLStore) Log(userID int64, ip string) error {
	insq := "insert into userLog(userID, inTime, clientIP) values (?,?,?)"
	_, execErr := ms.Db.Exec(insq, userID, time.Now(), ip)
	if execErr != nil {
		return errors.New("could not insert new log")
	}

	return nil
}

// GetByEmail returns User with given email
func (ms *MySQLStore) GetByEmail(email string) (*User, error) {
	result := &User{}
	rows, err := ms.Db.Query("SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Email, &result.PassHash, &result.UserName, &result.FirstName, &result.LastName, &result.PhotoURL)
		if err != nil {
			return nil, err
		}
		// take only first row
		break
	}
	return result, nil
}

//GetByUserName returns *User with given username
func (ms *MySQLStore) GetByUserName(email string) (*User, error) {
	result := &User{}
	rows, err := ms.Db.Query("SELECT * FROM users WHERE username=?", email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Email, &result.PassHash, &result.UserName, &result.FirstName, &result.LastName, &result.PhotoURL)
		if err != nil {
			return nil, err
		}
		// take only first row
		break
	}
	return result, nil
}

//ErrUpdatingUser is returned when user can be found not updated
var ErrUpdatingUser = errors.New("error updating user")

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ms *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	user, getErr := ms.GetByID(id)
	if getErr != nil {
		return nil, getErr
	}
	// update user
	if err := user.ApplyUpdates(updates); err != nil {
		return nil, ErrUpdatingUser
	}

	insq := "UPDATE users SET first_name=?, last_name=? WHERE id=?"
	if _, err := ms.Db.Exec(insq, user.FirstName, user.LastName, id); err != nil {
		return nil, err
	}

	return user, nil
}

//Delete deletes the user with the given ID
func (ms *MySQLStore) Delete(id int64) error {
	insq := "DELETE FROM users WHERE id=?"
	result, err := ms.Db.Exec(insq, id)
	if err != nil {
		return errors.New("unexpected error executing query")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.New("unexpected error evaluating result")
	}
	if affected < 1 {
		return ErrDeletingUser
	}

	return nil
}
