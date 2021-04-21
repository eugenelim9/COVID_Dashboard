package users

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestGetByID is a test function for the SQLStore's GetByID
func TestGetByID(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		idToGet      int64
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			2,
			true,
		},
		{
			"User With Large ID Found",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1234567890,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStore{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		// query used in your Store implementation
		query := regexp.QuoteMeta("SELECT * FROM users WHERE id=?")

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnError(ErrUserNotFound)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnRows(row)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// Unit test for SQL store's Insert function
func TestInsert(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		newUser      *User
	}{
		{
			"Insert user1",
			&User{
				0,
				"test@test.com",
				[]byte("passhash123"),
				"StevieG",
				"Steve",
				"G",
				"coolphoturl",
			},
			&User{
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "StevieG",
				FirstName: "Steve",
				LastName:  "G",
				PhotoURL:  "coolphoturl",
			},
		},
		{
			"Insert user2 with empty names",
			&User{ID: 1,
				Email:    "test@test.com",
				PassHash: []byte("passhash123"),
				UserName: "StevieG",
				PhotoURL: "coolphoturl",
			},
			&User{
				Email:    "test@test.com",
				PassHash: []byte("passhash123"),
				UserName: "StevieG",
				PhotoURL: "coolphoturl",
			},
		},
		{
			"Insert user3",
			&User{ID: 2,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "chilldude",
				FirstName: "John",
				LastName:  "Johnson",
				PhotoURL:  "coolphoturl",
			},
			&User{
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "chilldude",
				FirstName: "John",
				LastName:  "Johnson",
				PhotoURL:  "coolphoturl",
			},
		},
	}

	/*
		we can try insert each user wo id into the mock db. Then we can check that it returns the expected user,
		has correct row number and number of rows affected
		We can't specify id anyways so don't need to handle duplicates
		All the fields that are actually needed will be setted as well, or empty string so we gucci.
	*/

	// 2 checks in total
	// Just need to test for correctness
	// make sure the expectations checkout with result
	// correct user -> same as expected,
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("There was a problem opening a database connection: [%v]", err)
	}
	defer db.Close()
	mainSQLStore := &MySQLStore{db}

	for _, c := range cases {

		query := regexp.QuoteMeta("insert into users(email, pass_hash, username, first_name, last_name, photo_url) values (?,?,?,?,?,?)")
		mock.ExpectExec(query).WithArgs(c.newUser.Email, c.newUser.PassHash, c.newUser.UserName,
			c.newUser.FirstName, c.newUser.LastName, c.newUser.PhotoURL).
			WillReturnResult(sqlmock.NewResult(c.expectedUser.ID, 1))

		// test Insert()
		user, err := mainSQLStore.Insert(c.newUser)
		if err != nil {
			t.Errorf("Unexpected error [%s]: %v", c.name, err)
		}

		// see if expected user gets returned
		if !reflect.DeepEqual(user, c.expectedUser) {
			t.Errorf("Error, invalid match in test [%s]:\n\texpected %+v \n\treceived:%+v", c.name, *c.expectedUser, *user)
		}

		// see if all epectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// Unit test for SQL store's Update function
func TestUpdate(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		updateID     int64
		updates      *Updates
		expectedUser *User
		expectError  bool
	}{
		{
			"Update existing user",
			0,
			&Updates{
				FirstName: "Jack",
				LastName:  "City",
			},
			&User{
				0,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"Jack",
				"City",
				"photourl",
			},
			false,
		},
		{
			"Update invalid",
			0,
			&Updates{
				FirstName: "",
				LastName:  "",
			},
			&User{
				0,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			true,
		},
		{
			"Update valid ii",
			0,
			&Updates{
				FirstName: "Christian",
				LastName:  "",
			},
			&User{
				0,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"Christian",
				"lastname",
				"photourl",
			},
			false,
		},
	}

	// cases: exists, tryna update but it doesn't exist, invalid update name
	//

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()
		mainSQLStore := &MySQLStore{db}

		// Add expected row to db
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		updateQuery := regexp.QuoteMeta("UPDATE users SET first_name=?, last_name=? WHERE id=?")
		selectQuery := regexp.QuoteMeta("SELECT * FROM users WHERE id=?")

		if c.expectError == true {
			mock.ExpectQuery(selectQuery).WithArgs(c.updateID).WillReturnRows(row)
			user, err := mainSQLStore.Update(c.updateID, c.updates)
			if user != nil && err == nil {
				t.Errorf("Expected error but got [%v] instead", err)
			}
		} else {
			// no error
			mock.ExpectQuery(selectQuery).WithArgs(c.updateID).WillReturnRows(row)
			mock.ExpectExec(updateQuery).WithArgs(c.expectedUser.FirstName, c.expectedUser.LastName, c.updateID).
				WillReturnResult(sqlmock.NewResult(c.updateID, 1))

			// test update
			user, err := mainSQLStore.Update(c.updateID, c.updates)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			// see if expected user gets returned
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]:\n\texpected %+v \n\treceived:%+v", c.name, *c.expectedUser, *user)
			}
		}

		// see if all epectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

func TestDelete(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		id           int64
		expectError  bool
	}{
		{
			"Delete user1",
			&User{
				0,
				"test@test.com",
				[]byte("passhash123"),
				"StevieG",
				"Steve",
				"G",
				"coolphoturl",
			},
			0,
			false,
		},
		{
			"Delete user with long id",
			&User{
				12345,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			12345,
			false,
		},
		{
			"Delete user with not found id",
			&User{},
			0,
			true,
		},
	}
	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStore{db}
		query := regexp.QuoteMeta("DELETE FROM users WHERE id=?")
		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectExec(query).WithArgs(c.id).WillReturnResult(sqlmock.NewResult(0, 0))

			// Test Delete()
			err := mainSQLStore.Delete(c.id)
			if err == nil {
				t.Errorf("Test case: [%s] Expected error [%v] but got [%v] instead", c.name, ErrDeletingUser, err)
			}
		} else {

			// Set up an expected query
			mock.ExpectExec(query).WithArgs(c.id).WillReturnResult(sqlmock.NewResult(0, 1))

			// Test Delete()
			err := mainSQLStore.Delete(c.id)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}

		}

		// see if all epectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Test case: [%s] There were unfulfilled expectations: %s", c.name, err)
		}

	}
}

// TestGetByEmail is a test function for the SQLStore's GetByEmail
func TestGetByEmail(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		emailToGet   string
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"test@test.com",
			false,
		},
		{
			"User Not Found",
			&User{},
			"noemail@test.com",
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStore{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		// query used in your Store implementation
		query := regexp.QuoteMeta("SELECT * FROM users WHERE email=?")

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnError(ErrUserNotFound)

			// Test GetByID()
			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnRows(row)

			// Test GetByID()
			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByUsername is a test function for the SQLStore's GetByUsername
func TestGetByUsername(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name          string
		expectedUser  *User
		usernameToGet string
		expectError   bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"username",
			false,
		},
		{
			"User Not Found",
			&User{},
			"nouser",
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStore{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		// query used in your Store implementation
		query := regexp.QuoteMeta("SELECT * FROM users WHERE username=?")

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.usernameToGet).WillReturnError(ErrUserNotFound)

			// Test GetByUserName()
			user, err := mainSQLStore.GetByUserName(c.usernameToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.usernameToGet).WillReturnRows(row)

			// Test GetByID()
			user, err := mainSQLStore.GetByUserName(c.usernameToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}
