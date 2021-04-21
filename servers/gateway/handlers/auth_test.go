package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/my/repo/servers/gateway/models/users"
	"github.com/my/repo/servers/gateway/sessions"
)

// private function to verify correct and incorrect input for TestSessionsHandler
func verifyUsersHandlerOutput(c struct {
	name               string
	ctx                *HandlerContext
	method             string
	contentType        string
	expectedStatusCode int
	expectedUser       *users.User
	toCreateUser       *users.NewUser
}, rr *httptest.ResponseRecorder, t *testing.T) {
	if c.expectedUser != nil {
		// check for correct contenttype
		if !strings.HasPrefix(rr.Header().Get("Content-Type"), contentTypeJSON) {
			t.Errorf("case [%s] unexpected error wrong Content-Type header -> expected: %s received: %s",
				c.name, contentTypeJSON, rr.Header().Get("Content-Type"))
		}
		// statuscode
		if rr.Code != c.expectedStatusCode {
			t.Errorf("case [%s] unexpected status code -> expected: %d received: %d", c.name, c.expectedStatusCode,
				rr.Result().StatusCode)
		}
		// body response
		resultingUser := &users.User{}
		if err := json.Unmarshal(rr.Body.Bytes(), resultingUser); err != nil {
			t.Errorf("case [%s] unexpected error %s", c.name, err)
		}
		if resultingUser.UserName != c.expectedUser.UserName {
			t.Errorf("case [%s] unexpected output mismatch -> expected: user with username %s received :%s", c.name,
				c.expectedUser.UserName, resultingUser.UserName)
		}
		if resultingUser.ID != c.expectedUser.ID {
			t.Errorf("case [%s] unexpected output mismatch -> expected: user with userID %v received :%v", c.name,
				c.expectedUser.ID, resultingUser.ID)
		}
		if resultingUser.PhotoURL != c.expectedUser.PhotoURL {
			t.Errorf("case [%s] unexpected output mismatch -> expected: user with PhotoURL %s received :%s", c.name,
				c.expectedUser.PhotoURL, resultingUser.PhotoURL)
		}
		// expecting error and no user
	} else {
		if rr.Result().StatusCode != c.expectedStatusCode {
			t.Errorf("case [%s] unexpected status code -> expected: %d received :%d", c.name, c.expectedStatusCode,
				rr.Result().StatusCode)
		}
	}
}

func TestUsersHandler(t *testing.T) {

	signingKey := "the key"
	newUser := &users.NewUser{Email: "test@user.com", Password: "password", PasswordConf: "password",
		UserName: "LigmaB", FirstName: "Ligma", LastName: "Balls"}
	createdUser, err := newUser.ToUser()
	createdUser.ID = 1
	if err != nil {
		log.Fatalf("%s", err)
	}
	// Create session using memstore
	memstore := sessions.NewMemStore(0, 0)
	sid, err := sessions.NewSessionID(signingKey)
	if err != nil {
		log.Fatalf("%s", err)
	}
	// TODO:
	if err := memstore.Save(sid, newUser); err != nil {
		log.Fatalf("%s", err)
	}
	// Make user store
	userStore := &users.FakeSQLStore{TestUser: &users.User{ID: 0}}

	// Create context - two with different signin keys
	//noUserContext := &HandlerContext{"different key", memstore, userStore}
	// ensuring that each input has the correct output, header, status code, body in both correct and incorrect cases
	// content type not json

	cases := []struct {
		name               string
		ctx                *HandlerContext
		method             string
		contentType        string
		expectedStatusCode int
		expectedUser       *users.User
		toCreateUser       *users.NewUser
	}{
		{
			"valid POST request",
			&HandlerContext{signingKey, memstore, userStore},
			"POST",
			contentTypeJSON,
			http.StatusCreated,
			createdUser,
			newUser,
		},
		{
			"Invalid Method request",
			&HandlerContext{signingKey, memstore, userStore},
			"PATCH",
			contentTypeJSON,
			http.StatusMethodNotAllowed,
			nil,
			newUser,
		},
		{
			"Invalid header request",
			&HandlerContext{signingKey, memstore, userStore},
			"POST",
			"text/plain",
			http.StatusUnsupportedMediaType,
			nil,
			newUser,
		},
		{
			"POST wiht no user body in request",
			&HandlerContext{signingKey, memstore, userStore},
			"POST",
			contentTypeJSON,
			http.StatusBadRequest,
			nil,
			nil,
		},
	}

	for _, c := range cases {
		// Create JSON of user that will be passed into body
		userJSON, err := json.Marshal(c.toCreateUser)
		bodyUser := bytes.NewReader(userJSON)
		req, err := http.NewRequest(c.method, "", bodyUser)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", c.contentType)
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.ctx.UsersHandler)
		handler.ServeHTTP(rr, req)

		// if status := rr.Code; status != c.expectedStatusCode {
		// 	t.Errorf("handler returned wrong status code: got %v want %v",
		// 		status, c.expectedStatusCode)
		// }
		// returnedUser := &users.User{}
		// dec := json.NewDecoder(rr.Body)
		// dec.Decode(returnedUser)
		// Test to see if the userID being created and returned is the same.
		// if c.expectedUser.ID != returnedUser.ID {
		// 	t.Errorf("Wrong ID returned: got %v want %v",
		// 		returnedUser.ID, c.expectedUser.ID)
		// }
		verifyUsersHandlerOutput(c, rr, t)

	}
}

// private function to verify correct and incorrect input for TestSpecificUser GET and PATCH
func verifySpecificUserHandlerOutput(c struct {
	name               string
	ctx                *HandlerContext
	method             string
	contentType        string
	requestedUserID    string
	expectedStatusCode int
	expectedUser       *users.User
}, rr *httptest.ResponseRecorder, t *testing.T) {
	if c.expectedUser != nil {
		// check for correct contenttype
		if !strings.HasPrefix(rr.Header().Get("Content-Type"), contentTypeJSON) {
			t.Errorf("case [%s] unexpected error wrong Content-Type header -> expected: %s received: %s",
				c.name, contentTypeJSON, rr.Header().Get("Content-Type"))
		}
		// statuscode
		if rr.Result().StatusCode != c.expectedStatusCode {
			t.Errorf("case [%s] unexpected status code -> expected: %d received: %d", c.name, c.expectedStatusCode,
				rr.Result().StatusCode)
		}
		// body response
		resultingUser := &users.User{}
		if err := json.Unmarshal(rr.Body.Bytes(), resultingUser); err != nil {
			t.Errorf("case [%s] unexpected error %s", c.name, err)
		}
		if resultingUser.UserName != c.expectedUser.UserName {
			t.Errorf("case [%s] unexpected output mismatch -> expected: user with username %s received :%s", c.name,
				c.expectedUser.UserName, resultingUser.UserName)
		}
		// expecting error and no user
	} else {
		if rr.Result().StatusCode != c.expectedStatusCode {
			t.Errorf("case [%s] unexpected status code -> expected: %d received :%d", c.name, c.expectedStatusCode,
				rr.Result().StatusCode)
		}
	}
}

// TODO: could be refractored and simplified
func TestSpecificUserHandler(t *testing.T) {
	signingKey := "the key"
	testUser := &users.User{ID: 0, Email: "test@user.com", PassHash: []byte("password"),
		UserName: "StevieG", FirstName: "Steven", LastName: "Gerrard", PhotoURL: "coolimageurl"}

	// make memstore for session and input one user to Memstore
	sStore := sessions.NewMemStore(0, 0)
	sid, err := sessions.NewSessionID(signingKey)
	if err != nil {
		t.Fatalf("unexpected test error %s", err)
	}
	if err := sStore.Save(sid, SessionState{time.Now(), testUser}); err != nil {
		t.Fatalf("unexpected test error %s", err)
	}
	// make user store
	uStore := &users.FakeSQLStore{TestUser: testUser}

	// make context that will work for all cases
	context := &HandlerContext{signingKey, sStore, uStore}
	noUserContext := &HandlerContext{"different key", sStore, uStore}

	// user update and updated user for PATCH
	userUpdate := &users.Updates{FirstName: "jack", LastName: "mack"}

	updatedUser := &users.User{ID: 0, Email: "test@user.com", PassHash: []byte("password"),
		UserName: "StevieG", FirstName: "Steven", LastName: "Gerrard", PhotoURL: "coolimageurl"}
	if err := updatedUser.ApplyUpdates(userUpdate); err != nil {
		t.Fatalf("unexpected  test error %s", err)
	}
	cases := []struct {
		name               string
		ctx                *HandlerContext
		method             string
		contentType        string
		requestedUserID    string
		expectedStatusCode int
		expectedUser       *users.User
	}{
		{
			"valid GET request",
			context,
			"GET",
			"",
			"0",
			http.StatusOK,
			testUser,
		},
		{
			"valid GET request with /me",
			context,
			"GET",
			"",
			"me",
			http.StatusOK,
			testUser,
		},
		{
			"GET request with no current user",
			noUserContext,
			"GET",
			"",
			"0",
			http.StatusUnauthorized,
			nil,
		},
		{
			"GET with error parsing id",
			context,
			"GET",
			"",
			"zzz1",
			http.StatusNotFound,
			nil,
		},
		{
			"unimplemented method",
			context,
			"POST",
			"",
			"0",
			http.StatusMethodNotAllowed,
			nil,
		},
		{
			"GET request user not found",
			context,
			"GET",
			"",
			"2",
			http.StatusNotFound,
			nil,
		},
		{
			"valid PATCH",
			context,
			"PATCH",
			contentTypeJSON,
			"0",
			http.StatusOK,
			updatedUser,
		},
		{
			"valid PATCH with /me",
			context,
			"PATCH",
			"application/json",
			"me",
			http.StatusOK,
			updatedUser,
		},
		{
			"PATCH unmatched user",
			context,
			"PATCH",
			contentTypeJSON,
			"2",
			http.StatusForbidden,
			nil,
		},
		{
			"PATCH wrong content type",
			context,
			"PATCH",
			"text/html",
			"0",
			http.StatusUnsupportedMediaType,
			nil,
		},
		{
			"PATCH with nouser context",
			noUserContext,
			"PATCH",
			contentTypeJSON,
			"0",
			http.StatusUnauthorized,
			nil,
		},
		{
			"PATCH with error parsing id",
			context,
			"PATCH",
			contentTypeJSON,
			"zzz1",
			http.StatusNotFound,
			nil,
		},
	}
	for _, c := range cases {
		// set up handler and request recorder
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.ctx.SpecificUserHandler)
		if c.method == "GET" {
			req, err := http.NewRequest(c.method, fmt.Sprintf("/v1/users/%s", c.requestedUserID), nil)
			if err != nil {
				t.Errorf("case [%s] unexpected error %s", c.name, err)
			}
			req.Header.Set("Authorization", "Bearer "+fmt.Sprintf("%v", sid))
			handler.ServeHTTP(rr, req)
			// private method to verify correct and incorrect input
			verifySpecificUserHandlerOutput(c, rr, t)
		} else if c.method == "PATCH" {
			userJSON, err := json.Marshal(userUpdate)
			if err != nil {
				t.Errorf("case [%s] unexpected error %s", c.name, err)
			}
			body := bytes.NewReader(userJSON)
			req, err := http.NewRequest(c.method, fmt.Sprintf("/v1/users/%s", c.requestedUserID), body)
			if err != nil {
				t.Errorf("case [%s] unexpected error %s", c.name, err)
			}
			req.Header.Set("Content-Type", c.contentType)
			req.Header.Set("Authorization", "Bearer "+fmt.Sprintf("%v", sid))
			handler.ServeHTTP(rr, req)
			// private method to verify correct and incorrect input
			verifySpecificUserHandlerOutput(c, rr, t)
		} else {
			// handle unsupported method
			req, err := http.NewRequest(c.method, "", nil)
			if err != nil {
				t.Errorf("case [%s] unexpected error %s", c.name, err)
			}
			req.Header.Set("Authorization", "Bearer "+fmt.Sprintf("%v", sid))
			handler.ServeHTTP(rr, req)
			if rr.Result().StatusCode != c.expectedStatusCode {
				t.Errorf("case [%s] expected status code -> expected: %d received: %d", c.name,
					c.expectedStatusCode, rr.Result().StatusCode)
			}
		}
	}
}

// private function to verify correct and incorrect input for TestSessionsHandler
func verifySessionsHandlerOutput(c struct {
	name               string
	method             string
	contentType        string
	ctx                *HandlerContext
	inputCredentials   *users.Credentials
	expectedStatusCode int
	expectedOutput     *users.User
}, rr *httptest.ResponseRecorder, t *testing.T) {
	if c.expectedOutput != nil {
		// check for correct contenttype
		if !strings.HasPrefix(rr.Header().Get("Content-Type"), contentTypeJSON) {
			t.Errorf("case [%s] unexpected error wrong Content-Type header -> expected: %s received: %s",
				c.name, contentTypeJSON, rr.Header().Get("Content-Type"))
		}
		// statuscode
		if rr.Result().StatusCode != c.expectedStatusCode {
			t.Errorf("case [%s] unexpected status code -> expected: %d received: %d", c.name, c.expectedStatusCode,
				rr.Result().StatusCode)
		}
		// body response
		resultingUser := &users.User{}
		if err := json.Unmarshal(rr.Body.Bytes(), resultingUser); err != nil {
			t.Errorf("case [%s] unexpected error %s", c.name, err)
		}
		if resultingUser.UserName != c.expectedOutput.UserName {
			t.Errorf("case [%s] unexpected output mismatch -> expected: user with username %s received :%s", c.name,
				c.expectedOutput.UserName, resultingUser.UserName)
		}
		// expecting error and no user
	} else {
		if rr.Result().StatusCode != c.expectedStatusCode {
			t.Errorf("case [%s] unexpected status code -> expected: %d received :%d", c.name, c.expectedStatusCode,
				rr.Result().StatusCode)
		}
	}
}

func TestSessionsHandler(t *testing.T) {
	// test user that exists in user store
	testUser := &users.User{ID: 0, Email: "test@user.com",
		UserName: "StevieG", FirstName: "Steven", LastName: "Gerrard", PhotoURL: "coolimageurl"}

	// set password properly
	if err := testUser.SetPassword("password"); err != nil {
		t.Fatalf("unexpected test error %s", err)
	}

	cases := []struct {
		name               string
		method             string
		contentType        string
		ctx                *HandlerContext
		inputCredentials   *users.Credentials
		expectedStatusCode int
		expectedOutput     *users.User
	}{
		{"Valid POST request",
			"POST",
			contentTypeJSON,
			&HandlerContext{"the key",
				sessions.NewMemStore(0, 0),
				&users.FakeSQLStore{TestUser: testUser},
			},
			&users.Credentials{Email: "test@user.com", Password: "password"},
			http.StatusCreated,
			testUser,
		},
		{"Non POST request",
			"PATCH",
			contentTypeJSON,
			&HandlerContext{"the key",
				sessions.NewMemStore(0, 0),
				&users.FakeSQLStore{TestUser: testUser},
			},
			&users.Credentials{Email: "test@user.com", Password: "password"},
			http.StatusMethodNotAllowed,
			nil,
		},
		{"POST request wrong Content-Type",
			"POST",
			"text/html",
			&HandlerContext{"the key",
				sessions.NewMemStore(0, 0),
				&users.FakeSQLStore{TestUser: testUser},
			},
			&users.Credentials{Email: "test@user.com", Password: "password"},
			http.StatusUnsupportedMediaType,
			nil,
		},
		{"POST request user not found",
			"POST",
			contentTypeJSON,
			&HandlerContext{"the key",
				sessions.NewMemStore(0, 0),
				&users.FakeSQLStore{TestUser: testUser},
			},
			&users.Credentials{Email: "invalid@user.com", Password: "password"},
			http.StatusUnauthorized,
			nil,
		},
		{"POST request invalid password",
			"POST",
			contentTypeJSON,
			&HandlerContext{"the key",
				sessions.NewMemStore(0, 0),
				&users.FakeSQLStore{TestUser: testUser},
			},
			&users.Credentials{Email: "test@user.com", Password: "ehhhhhhh"},
			http.StatusUnauthorized,
			nil,
		},
	}

	for _, c := range cases {
		// setting up handler, response recorder and request to be passed in
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.ctx.SessionsHandler)
		// marshalling credentials into byte slice
		credJSON, err := json.Marshal(c.inputCredentials)
		if err != nil {
			t.Errorf("case [%s] unexpected error %s", c.name, err)
		}
		body := bytes.NewReader(credJSON)
		req, err := http.NewRequest(c.method, "", body)
		if err != nil {
			t.Errorf("case [%s] unexpected error %s", c.name, err)
		}
		req.Header.Set("Content-Type", c.contentType)
		// serve http
		handler.ServeHTTP(rr, req)
		// check output if we need to
		verifySessionsHandlerOutput(c, rr, t)
	}
}

func verifySpecificSessionHandlerOutput(c struct {
	name               string
	method             string
	ctx                *HandlerContext
	requestedUserID    string
	expectedStatusCode int
	expectedOutput     string
}, rr *httptest.ResponseRecorder, t *testing.T) {
	if c.expectedOutput != "" {

		if rr.Body.String() != c.expectedOutput {
			t.Errorf("case [%s] unexpected output mismatch -> expected: %s received :%s", c.name,
				"signed out", rr.Body.String())
		}
		// statuscode
		if rr.Result().StatusCode != c.expectedStatusCode {
			t.Errorf("case [%s] unexpected status code -> expected: %d received: %d", c.name, c.expectedStatusCode,
				rr.Result().StatusCode)
		}

		// expecting error and no user
	} else {
		if rr.Result().StatusCode != c.expectedStatusCode {
			t.Errorf("case [%s] unexpected status code -> expected: %d received :%d", c.name, c.expectedStatusCode,
				rr.Result().StatusCode)
		}
	}
}

func TestSpecificSessionHandler(t *testing.T) {
	// test user that exists in user store
	signingKey := "the key"
	testUser := &users.User{ID: 0, Email: "test@user.com", PassHash: []byte("password"),
		UserName: "StevieG", FirstName: "Steven", LastName: "Gerrard", PhotoURL: "coolimageurl"}

	// make memstore for session and input one user to Memstore
	sStore := sessions.NewMemStore(0, 0)
	sid, err := sessions.NewSessionID(signingKey)
	if err != nil {
		t.Fatalf("unexpected test error %s", err)
	}
	if err := sStore.Save(sid, SessionState{time.Now(), testUser}); err != nil {
		t.Fatalf("unexpected test error %s", err)
	}

	cases := []struct {
		name            string
		method          string
		ctx             *HandlerContext
		requestedUserID string
		// inputCredentials   *users.Credentials
		expectedStatusCode int
		expectedOutput     string
	}{
		{"Valid DELETE request",
			"DELETE",
			&HandlerContext{"the key",
				sessions.NewMemStore(0, 0),
				&users.FakeSQLStore{TestUser: testUser},
			},
			"mine",
			http.StatusOK,
			"signed out",
		},
		{"Invalid request user URL",
			"DELETE",
			&HandlerContext{"the key",
				sessions.NewMemStore(0, 0),
				&users.FakeSQLStore{TestUser: testUser},
			},
			"0",
			http.StatusForbidden,
			"",
		},
		{"Invalid MethodL",
			"GET",
			&HandlerContext{"the key",
				sessions.NewMemStore(0, 0),
				&users.FakeSQLStore{TestUser: testUser},
			},
			"0",
			http.StatusMethodNotAllowed,
			"",
		},
	}

	for _, c := range cases {
		// setting up handler, response recorder and request to be passed in
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.ctx.SpecificSessionHandler)
		req, err := http.NewRequest(c.method, fmt.Sprintf("/v1/sessions/%s", c.requestedUserID), nil)
		if err != nil {
			t.Errorf("case [%s] unexpected error %s", c.name, err)
		}
		req.Header.Set("Authorization", "Bearer "+fmt.Sprintf("%v", sid))
		// serve http

		handler.ServeHTTP(rr, req)
		// check output if we need to
		verifySpecificSessionHandlerOutput(c, rr, t)
	}
}
