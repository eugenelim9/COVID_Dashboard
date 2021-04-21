package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/my/repo/servers/gateway/models/users"
	"github.com/my/repo/servers/gateway/sessions"
	"golang.org/x/crypto/bcrypt"
)

var contentTypeJSON = "application/json"

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

// UsersHandler handles requests for user resourses
func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), contentTypeJSON) {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		incomingUser := &users.NewUser{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(incomingUser); err != nil {
			http.Error(w, "error decoding json", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// ensure new user is valid and convert newUser into User
		user, err := incomingUser.ToUser()
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}
		// check for duplicate email and id
		checkEmailUser, err := ctx.UserStore.GetByEmail(user.Email)
		if err == nil && len(checkEmailUser.Email) != 0 {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}

		checkUserNameUser, err := ctx.UserStore.GetByUserName(user.UserName)
		if err == nil && len(checkUserNameUser.UserName) != 0 {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}

		// creating new user in database
		savedUser, err := ctx.UserStore.Insert(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}
		newSessionState := SessionState{time.Now(), savedUser}

		// begin new session
		_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, newSessionState, w)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		// respond to client
		w.Header().Add("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		if err := enc.Encode(savedUser); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "request error", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificUserHandler handles request for specific users
func (ctx *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	// The current user must be authenticated to call this handler regardless of HTTP method.
	// If the user is not authenticated, respond immediately with an http.StatusUnauthorized (401) error status code
	// Authenticate user
	sessionState := &SessionState{}
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(w, "sesssion unauthorized please log in", http.StatusUnauthorized)
		return
	}

	method := r.Method
	if method == "GET" {
		// get user id
		urlSlice := strings.Split(r.URL.Path, "/")
		userIDString := urlSlice[len(urlSlice)-1]
		var userID int64
		if userIDString == "me" {
			userID = sessionState.User.ID
		} else {
			userID, err = strconv.ParseInt(userIDString, 10, 64)
			if err != nil {
				http.Error(w, "user does not exist", http.StatusNotFound)
				return
			}
		}

		// return user if found and StatusOK, if not found return StatusNotFound
		user, err := ctx.UserStore.GetByID(userID)
		if err != nil || len(user.UserName) == 0 {
			http.Error(w, "user does not exist", http.StatusNotFound)
			return
		}
		w.Header().Add("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		if err := enc.Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}
	} else if method == "PATCH" {
		// If the user ID in the request URL is not "me" or does not match the currently-authenticated user,
		// immediately respond with an http.StatusForbidden (403) error status code and appropriate error message.
		urlSlice := strings.Split(r.URL.Path, "/")
		userIDString := urlSlice[len(urlSlice)-1]
		var userID int64
		if userIDString == "me" {
			userID = sessionState.User.ID
		} else {
			userID, err = strconv.ParseInt(userIDString, 10, 64)
			if err != nil {
				http.Error(w, "user does not exist", http.StatusNotFound)
				return
			}
			// check if it matches with current user
			if userID != sessionState.User.ID {
				http.Error(w, "request not authorized", http.StatusForbidden)
				return
			}
		}
		// If the request's Content-Type header does not start with application/json, respond with status code
		// http.StatusUnsupportedMediaType (415), and a message indicating that the request body must be in JSON.
		if !strings.HasPrefix(r.Header.Get("Content-Type"), contentTypeJSON) {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		// The request body should contain JSON that can be decoded into the users.Updates struct.
		// Use that to update the user's profile.
		userUpdates := &users.Updates{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(userUpdates); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		// close response body
		defer r.Body.Close()

		user, err := ctx.UserStore.Update(userID, userUpdates)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		if err := enc.Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

	} else {
		http.Error(w, "request error", http.StatusMethodNotAllowed)
		return
	}
}

// SessionsHandler handles reqeust for "sessions" resources and allows clients to begin new session with existing credentials
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), contentTypeJSON) {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		// The request body should contain JSON that can be decoded into a users.Credentials struct. Use those
		// credentials to find the user profile and authenticate. If you don't find the user profile, do something that would
		// take about the same amount of time as authenticating, and then respond with a http.StatusUnauthorized error status code
		// and generic message "invalid credentials". Respond with the same error if you find the profile but fail to authenticate.
		cred := &users.Credentials{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(cred); err != nil {
			http.Error(w, "error decoding json", http.StatusBadRequest)
			return
		}

		// close body
		defer r.Body.Close()

		// get user given that email
		user, err := ctx.UserStore.GetByEmail(cred.Email)
		if err != nil {
			// user not found do fake comparison return error
			bcrypt.CompareHashAndPassword([]byte("dummypassword"), []byte("dummypassword"))
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// do the auth
		if err := user.Authenticate(cred.Password); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// If authentication is successful, begin a new session.
		newSessionState := &SessionState{time.Now(), user}
		_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, newSessionState, w)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		// Insert Log
		ip := GetIP(r)
		ctx.UserStore.Log(user.ID, ip)

		// Respond to client
		w.Header().Add("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		if err := enc.Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "request error", http.StatusMethodNotAllowed)
		return
	}
}

// GetIP is a helper function for getting IP address from the request
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		ips := strings.Split(forwarded, ", ")
		return ips[0]
	}
	return r.RemoteAddr
}

// SpecificSessionHandler handles requests related to a specific authenticated session
func (ctx *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		// If the last path segment does not equal "mine", immediately respond with an http.StatusForbidden (403) error
		// status code and appropriate error message.
		urlSlice := strings.Split(r.URL.Path, "/")
		lastSegment := urlSlice[len(urlSlice)-1]
		if lastSegment != "mine" {
			http.Error(w, "request not authorized", http.StatusForbidden)
			return
		}
		// end current session
		if _, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}
		// respond with plain text
		w.Write([]byte("signed out"))
	} else {
		http.Error(w, "reqeust error", http.StatusMethodNotAllowed)
		return
	}
}
