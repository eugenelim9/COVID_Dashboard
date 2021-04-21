package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	//validate the new user according to these rules:
	//- Email field must be a valid email address (hint: see mail.ParseAddress)
	_, emailAddrErr := mail.ParseAddress(nu.Email)
	if emailAddrErr != nil {
		return errors.New("Invalid Email Address")
	}

	//- Password must be at least 6 characters
	// if len(nu.Password) < 6 {
	// 	return errors.New("Password must be at least 6 characters")
	// }
	//- Password and PasswordConf must match
	if nu.Password != nu.PasswordConf {
		return errors.New("Password and its confirmation does not match")
	}

	// Username mst not contain space
	hasSpace, err := regexp.MatchString(`\s`, nu.UserName)
	if err != nil {
		return err
	}

	//- UserName must be non-zero length and may not contain spaces
	if len(nu.UserName) == 0 || hasSpace == true {
		return errors.New("Username must be non-zero length and may not contain spaces")
	}
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	//call Validate() to validate the NewUser and
	//return any validation errors that may occur.
	if err := nu.Validate(); err != nil {
		return nil, err
	}
	//if valid, create a new *User and set the fields
	//based on the field values in `nu`.
	//Leave the ID field as the zero-value; your Store
	//implementation will set that field to the DBMS-assigned
	//primary key value.
	incomingUser := &User{Email: nu.Email, UserName: nu.UserName, FirstName: nu.FirstName, LastName: nu.LastName}

	//Set the PhotoURL field to the Gravatar PhotoURL
	//for the user's email address.
	//see https://en.gravatar.com/site/implement/hash/
	//and https://en.gravatar.com/site/implement/images/
	formattedEmail := strings.TrimSpace(strings.ToLower(nu.Email))
	photoHash := md5.Sum([]byte(formattedEmail))
	incomingUser.PhotoURL = fmt.Sprintf("%s%x", gravatarBasePhotoURL, photoHash)

	//also call .SetPassword() to set the PassHash
	//field of the User to a hash of the NewUser.Password
	if err := incomingUser.SetPassword(nu.Password); err != nil {
		return nil, err
	}

	return incomingUser, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	//implement according to comment above
	result := []string{}
	if len(strings.Replace(u.FirstName, " ", "", -1)) != 0 {
		result = append(result, u.FirstName)
	}
	if len(strings.Replace(u.LastName, " ", "", -1)) != 0 {
		result = append(result, u.LastName)
	}
	return strings.Join(result, " ")
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//use the bcrypt package to generate a new hash of the password
	//https://godoc.org/golang.org/x/crypto/bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = hash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//use the bcrypt package to compare the supplied
	//password with the stored PassHash
	//https://godoc.org/golang.org/x/crypto/bcrypt
	return bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	//set the fields of `u` to the values of the related
	//field in the `updates` struct
	if len(updates.FirstName) == 0 && len(updates.LastName) == 0 {
		return errors.New("Updates not applied, both fields are empty")
	}
	if len(updates.FirstName) != 0 {
		u.FirstName = updates.FirstName
	}

	if len(updates.LastName) != 0 {
		u.LastName = updates.LastName
	}
	return nil
}
