package users

import "errors"

//mport "assignments-towm1204/servers/gateway/models/users"

// TestUser is the test user which is the single default user in the fake mysql database

// FakeSQLStore is used for testing purpose where there  will be exactly one user in the database, defined as
// TestUser in this store.
type FakeSQLStore struct {
	TestUser *User
}

//GetByID returns the User with the given ID
func (fakestore *FakeSQLStore) GetByID(id int64) (*User, error) {
	if id == fakestore.TestUser.ID {
		return fakestore.TestUser, nil
	}
	return nil, errors.New("User not found")
}

//GetByEmail returns the User with the given email
func (fakestore *FakeSQLStore) GetByEmail(email string) (*User, error) {
	if email == fakestore.TestUser.Email {
		return fakestore.TestUser, nil
	}
	return nil, errors.New("user not found")

}

//GetByUserName returns the User with the given Username
func (fakestore *FakeSQLStore) GetByUserName(username string) (*User, error) {
	if username == fakestore.TestUser.UserName {
		return fakestore.TestUser, nil
	}
	return nil, errors.New("user not found")
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (fakestore *FakeSQLStore) Insert(user *User) (*User, error) {
	if fakestore.TestUser == nil {
		user.ID = 0
		return user, nil
	}
	user.ID = fakestore.TestUser.ID + 1
	return user, nil
}

//Log logs user logins into our database table
func (fakestore *FakeSQLStore) Log(userID int64, ip string) error {
	return nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (fakestore *FakeSQLStore) Update(id int64, updates *Updates) (*User, error) {
	if id == fakestore.TestUser.ID {
		firstName := fakestore.TestUser.FirstName
		lastName := fakestore.TestUser.LastName
		if len(updates.FirstName) != 0 {
			firstName = updates.FirstName
		}
		if len(updates.LastName) != 0 {
			lastName = updates.FirstName
		}
		return &User{fakestore.TestUser.ID, fakestore.TestUser.Email, fakestore.TestUser.PassHash, fakestore.TestUser.UserName,
			firstName, lastName, fakestore.TestUser.PhotoURL}, nil
	}
	return nil, errors.New("user not found")
}

//Delete deletes the user with the given ID
func (fakestore *FakeSQLStore) Delete(id int64) error {
	if id != fakestore.TestUser.ID {
		return errors.New("user not found")
	}
	return nil
}
