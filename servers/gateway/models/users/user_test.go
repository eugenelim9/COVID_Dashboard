package users

import (
	"crypto/md5"
	"fmt"
	"testing"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.
func TestNewUserValidate(t *testing.T) {
	// test cases
	cases := []struct {
		name  string
		nu    *NewUser
		valid bool
	}{
		{
			"valid user",
			&NewUser{
				Email:        "stevieG@gmail.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "thegoat",
				FirstName:    "Steven",
				LastName:     "Gerrard",
			},
			true,
		},
		// valid user 2
		{
			"valid user 2",
			&NewUser{
				Email:        "  A@m.com  ",
				Password:     "antherone",
				PasswordConf: "antherone",
				UserName:     "bigshaq",
				FirstName:    "Big",
				LastName:     "Shaq",
			},
			true,
		},
		// invalid email
		{
			"invalid email",
			&NewUser{
				Email:        "ohnobademail@",
				Password:     "antherone",
				PasswordConf: "antherone",
				UserName:     "bigshaq",
				FirstName:    "Big",
				LastName:     "Shaq",
			},
			false,
		},
		// no email
		{
			"no email",
			&NewUser{
				Email:        "",
				Password:     "antherone",
				PasswordConf: "antherone",
				UserName:     "bigshaq",
				FirstName:    "Big",
				LastName:     "Shaq",
			},
			false,
		},
		// short pass word
		{
			"too short pass word",
			&NewUser{
				Email:        "okay@email.com",
				Password:     "short",
				PasswordConf: "short",
				UserName:     "bigshaq",
				FirstName:    "Big",
				LastName:     "Shaq",
			},
			false,
		},
		// password mismatch
		{
			"password mismatch",
			&NewUser{
				Email:        "okay@email.com",
				Password:     "mismatch",
				PasswordConf: "matchmis",
				UserName:     "bigshaq",
				FirstName:    "Big",
				LastName:     "Shaq",
			},
			false,
		},
		// empty username
		{
			"empty username",
			&NewUser{
				Email:        "okay@email.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "",
				FirstName:    "Who",
				LastName:     "Dis",
			},
			false,
		},
		// space in username
		{
			"space in username",
			&NewUser{
				Email:        "okay@email.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "sike its the wrong number",
				FirstName:    "Who",
				LastName:     "Dis",
			},
			false,
		},
		// space in email
		{
			"space in email",
			&NewUser{
				Email:        "ok ay@email.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "Username",
				FirstName:    "Who",
				LastName:     "Dis",
			},
			false,
		},
		// space before/after email
		{
			"space before/after email",
			&NewUser{
				Email:        "  okay@email.com  ",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "SpaceBeforeAfterEmail",
				FirstName:    "Who",
				LastName:     "Dis",
			},
			true,
		},
		// White space in username
		{
			"White space in username",
			&NewUser{
				Email:        "  okay@email.com  ",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName: "SpaceBefo	reAfterEmail",
				FirstName: "Who",
				LastName:  "Dis",
			},
			false,
		},
	}
	for _, c := range cases {
		result := c.nu.Validate()
		if result == nil && c.valid == true {
			continue
		} else if result != nil && c.valid == false {
			continue
		}
		// bad
		t.Errorf("case %s -> New user valid = %t but evaluated as valid = %t", c.name, c.valid, !c.valid)
	}
}

func TestToUser(t *testing.T) {
	cases := []struct {
		name             string
		nu               *NewUser
		expectedPhotoURL string
		expectedPassword string
		expectError      bool
	}{
		{
			"Captial in email",
			&NewUser{
				Email:        "stevieG@gmail.com",
				Password:     "abc123",
				PasswordConf: "abc123",
				UserName:     "thegoat",
				FirstName:    "Steven",
				LastName:     "Gerrard",
			},
			gravatarBasePhotoURL + fmt.Sprintf("%x", md5.Sum([]byte("stevieg@gmail.com"))),
			"abc123",
			false,
		},
		// Spaces before and after email
		{
			"Spaces before and after email",
			&NewUser{
				Email:        "  stevieG@gmail.com  ",
				Password:     "abc123",
				PasswordConf: "abc123",
				UserName:     "thegoat",
				FirstName:    "Steven",
				LastName:     "Gerrard",
			},
			gravatarBasePhotoURL + fmt.Sprintf("%x", md5.Sum([]byte("stevieg@gmail.com"))),
			"abc123",
			false,
		},
		// Invalid new userName
		{
			"Invalid new userName",
			&NewUser{
				Email:        "  stevieG@gmail.com  ",
				Password:     "abc123",
				PasswordConf: "abc123",
				UserName:     "thego at",
				FirstName:    "Steven",
				LastName:     "Gerrard",
			},
			gravatarBasePhotoURL + fmt.Sprintf("%x", md5.Sum([]byte("stevieg@gmail.com"))),
			"abc123",
			true,
		},
	}
	for _, c := range cases {
		user, err := c.nu.ToUser()
		if err != nil {
			if c.expectError == false {
				t.Errorf("case %s: unexpected error validating case: %v", c.name, err)
			}
			// if expceted error came in, continue
			continue
		} else if c.expectError == true {
			t.Errorf("case %s: expected error but didn't get one: %v", c.name, err)
		}

		if user.PhotoURL != c.expectedPhotoURL {
			t.Errorf("case %+v -> User.PhotoURL: %s != expectedPhotoURL: %s", *c.nu, user.PhotoURL, c.expectedPhotoURL)
			continue
		}

		if err := user.Authenticate(c.expectedPassword); err != nil {
			t.Errorf("case %+v -> decrypted user passHash: %s != expectedPassword: %s", *c.nu, user.PhotoURL, c.expectedPassword)
		}

	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		name     string
		u        *User
		Expected string
	}{
		{"both valid",
			&User{FirstName: "Jim", LastName: "Kite"},
			"Jim Kite"},

		{"no firstname",
			&User{FirstName: "    ", LastName: "Kite"},
			"Kite"},

		{"no firstname ii",
			&User{LastName: "Kite"},
			"Kite"},

		{"no lastname",
			&User{FirstName: "Jim", LastName: ""},
			"Jim"},

		{"no lastname ii",
			&User{FirstName: "Jim"},
			"Jim"},

		{"no first and lastname",
			&User{FirstName: "", LastName: "    "},
			""},

		{"no first and lastname ii",
			&User{},
			""},
	}

	for _, c := range cases {
		result := c.u.FullName()
		if result != c.Expected {
			t.Errorf("case: %s -> expected: %s but received: %s", c.name, c.Expected, result)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	cases := []struct {
		name      string
		u         *User
		password  string
		expectErr bool
	}{
		{"Incorrect password",
			&User{},
			"thisisnottherightone",
			true,
		},

		{"Correct password",
			&User{},
			"thisistheone",
			false,
		},
		{"Empty password",
			&User{},
			"",
			true,
		},
	}

	for _, c := range cases {
		if err := c.u.SetPassword("thisistheone"); err != nil {
			t.Fatalf("Unexpected error setting passord")
		}

		authErr := c.u.Authenticate(c.password)
		if c.expectErr == false && authErr != nil {
			t.Errorf("case: %s -> expected error but authentication sucessful", c.name)
		}

		if c.expectErr == true && authErr == nil {
			t.Errorf("case: %s -> expected error but authentication sucessful", c.name)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name        string
		update      *Updates
		u           *User
		expectError bool
	}{
		{
			"Valid updates",
			&Updates{
				FirstName: "Ligma",
				LastName:  "Balls",
			},
			&User{
				FirstName: "name1",
				LastName:  "name2",
			},
			false,
		},
		{
			"Only 1 first name",
			&Updates{
				FirstName: "",
				LastName:  "Balls",
			},
			&User{
				FirstName: "name1",
				LastName:  "name2",
			},
			false,
		},
		{
			"Empty name",
			&Updates{
				FirstName: "",
				LastName:  "",
			},
			&User{
				FirstName: "name1",
				LastName:  "name2",
			},
			true,
		},
	}
	for _, c := range cases {
		err := c.u.ApplyUpdates(c.update)
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error validating case: %v", c.name, err)
			continue
		}
		if len(c.update.FirstName) != 0 && c.u.FirstName != c.update.FirstName && !c.expectError {
			t.Errorf("case %s: user firstname: %s != updated firstname: %s", c.name, c.u.FirstName, c.update.FirstName)
		}
		if len(c.update.LastName) != 0 && c.u.LastName != c.update.LastName && !c.expectError {
			t.Errorf("case %s: user lastname: %s != updated lastname: %s", c.name, c.u.LastName, c.update.LastName)
		}
	}
}
