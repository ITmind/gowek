package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {

	tests := []struct {
		name         string
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
	}{

		{
			name:         "get users",
			description:  "get HTTP status 200",
			route:        "/api/users",
			expectedCode: 200,
		},
		{
			name:         "get notes",
			description:  "get HTTP status 200",
			route:        "/api/notes",
			expectedCode: 200,
		},

		{
			name:         "page not found",
			description:  "get HTTP status 404, when route is not exists",
			route:        "/api/not-found",
			expectedCode: 404,
		},
	}

	app := NewApp()

	req := httptest.NewRequest("POST", "/api/token", strings.NewReader("user=admin&pass=admin"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode, "get token")

	tokenResponce := struct {
		Token string
	}{}

	json.NewDecoder(resp.Body).Decode(&tokenResponce)
	if tokenResponce.Token == "" {
		t.Fatal("token is empty")
	}
	resp.Body.Close()

	// Iterate through test single test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", test.route, nil)
			req.Header.Add("Authorization", "Bearer "+tokenResponce.Token)

			resp, _ := app.Test(req, 1)

			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})

	}
}
