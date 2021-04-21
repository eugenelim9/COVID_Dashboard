package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// private handler to test cors handler. It will return StatusAccepted instead of StatusOK for difference
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("this handler is just for test, does nothing and returns StatusAccepted"))
}

func TestCorsServeHTTP(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		expectedStatusCode int
	}{
		{
			"preflight request",
			"OPTIONS",
			http.StatusOK,
		}, {
			"normal request",
			"GET",
			http.StatusAccepted,
		},
	}

	for _, c := range cases {
		req, err := http.NewRequest(c.method, "", nil)
		if err != nil {
			t.Fatalf("case [%s] unexpected error making new request: %s", c.name, err)
		}
		rr := httptest.NewRecorder()
		wrappedCorsHandler := NewCORS(http.HandlerFunc(testHandler))
		wrappedCorsHandler.ServeHTTP(rr, req)
		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("case [%s] Access-Control-Allow-Origin header -> expected: %s received: %s", c.name, "*",
				rr.Header().Get("Access-Control-Allow-Origin"))
		}
		if rr.Header().Get("Access-Control-Allow-Methods") != "GET, PUT, POST, PATCH, DELETE" {
			t.Errorf("case [%s] Access-Control-Allow-Methods header -> expected: %s received: %s", c.name,
				"GET, PUT, POST, PATCH, DELETE", rr.Header().Get("Access-Control-Allow-Methods"))
		}
		if rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
			t.Errorf("case [%s] Access-Control-Allow-Headers header -> expected: %s received: %s", c.name,
				"Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
		}
		if rr.Header().Get("Access-Control-Expose-Headers") != "Authorization" {
			t.Errorf("case [%s] Access-Control-Expose-Headers header -> expected: %s received: %s", c.name, "Authorization",
				rr.Header().Get("Access-Control-Expose-Headers"))
		}
		if rr.Header().Get("Access-Control-Max-Age") != "600" {
			t.Errorf("case [%s] Access-Control-Max-Age header -> expected: %s received: %s", c.name, "600",
				rr.Header().Get("Access-Control-Max-Age"))
		}

		if rr.Result().StatusCode != c.expectedStatusCode {
			t.Errorf("case [%s] returned wrong status code -> expected: %d received: %d", c.name, c.expectedStatusCode,
				rr.Result().StatusCode)
		}
	}
}
