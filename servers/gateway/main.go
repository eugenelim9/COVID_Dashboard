package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/my/repo/servers/gateway/handlers"
	"github.com/my/repo/servers/gateway/models/users"
	"github.com/my/repo/servers/gateway/sessions"
)

// Director is a custom director to modify reverse proxy requests
type Director func(r *http.Request)

// CustomDashDirector is a wrapper class for reverse proxies so we can pass request on to dash microservices
func CustomDashDirector(targets []string, ctx *handlers.HandlerContext) Director {
	var counter int32
	counter = 0
	return func(r *http.Request) {
		// round robin selecting
		target := targets[counter%int32(len(targets))]
		atomic.AddInt32(&counter, 1)
		// find currently authenticated user
		sessionState := &handlers.SessionState{}
		_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, sessionState)
		if err != nil {
			r.Header.Del("X-User")
		} else {
			userByteSlice, err := json.Marshal(sessionState.User)
			if err == nil {
				r.Header.Set("X-User", string(userByteSlice[:]))
			}
		}
		r.Host = target
		r.URL.Host = target
		r.URL.Scheme = "http"
	}
}

//main is the main entry point for the server
func main() {
	/* - Read the ADDR environment variable to get the address
	the server should listen on. If empty, default to ":80" */
	addr := os.Getenv("ADDR")
	sessKey := os.Getenv("SESSIONKEY")
	redisaddr := os.Getenv("REDDISADDR")
	dsn := os.Getenv("DSN")
	dashboardAddresses := strings.Split(os.Getenv("DASHBOARDADDR"), ",")

	if len(addr) == 0 {
		addr = ":8443"
	}
	// new session store
	if len(redisaddr) == 0 {
		redisaddr = "127.0.0.1:6379"
	}
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	if len(tlsCertPath) == 0 || len(tlsCertPath) == 0 {
		log.Fatalln("TLSKEY and/or TLSCERT environment variables not set")
	}

	sessStore := sessions.NewRedisStore(redis.NewClient(&redis.Options{
		Addr: redisaddr,
	}), time.Hour)

	// new user store
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("error opening database")
	}
	userStore := &users.MySQLStore{Db: db}
	defer db.Close()

	if err := userStore.Db.Ping(); err != nil {
		log.Printf("error pinging database: %v\n", err)
	} else {
		log.Printf("successfully connected!\n")
	}

	// creating new context
	ctx := handlers.HandlerContext{SigningKey: sessKey, SessionStore: sessStore, UserStore: userStore}
	/*
		- Create a new mux for the web server. */
	mux := http.NewServeMux()

	dashProxy := &httputil.ReverseProxy{Director: CustomDashDirector(dashboardAddresses, &ctx)}

	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", ctx.SpecificSessionHandler)
	mux.Handle("/v1/dashboards", dashProxy)
	mux.Handle("/v1/dashboards/", dashProxy)
	mux.Handle("/v1/data", dashProxy)

	wrappedMux := handlers.NewCORS(mux)

	/*
		- Start a web server listening on the address you read from
		  the environment variable, using the mux you created as
		  the root handler. Use log.Fatal() to report any errors
		  that occur when trying to start the web server.
	*/
	log.Printf("Listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}
