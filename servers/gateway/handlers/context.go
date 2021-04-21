package handlers

import (
	"github.com/my/repo/servers/gateway/models/users"
	"github.com/my/repo/servers/gateway/sessions"
)

//define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

// HandlerContext will be a receiver for any http handler function that needs access to globals such as...
type HandlerContext struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
}
