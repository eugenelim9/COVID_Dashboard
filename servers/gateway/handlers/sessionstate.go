package handlers

import (
	"time"

	"github.com/my/repo/servers/gateway/models/users"
)

//define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

// SessionState holds session information keeping track of time and user
type SessionState struct {
	BeginTime time.Time   `json:"beginTime"`
	User      *users.User `json:"user"`
}
