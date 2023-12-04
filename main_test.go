package main

import (
	"bytes"
	"encoding/json"
	"gowek/repo"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApi(t *testing.T) {

	t.Setenv("DB_URL", t.TempDir()+"/test.db")
	t.Setenv("DB_TYPE", "sqlite")

	app := NewApp()
	defer app.Close()

	req := httptest.NewRequest("POST", "/token", strings.NewReader("username=admin&password=admin"))
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

	t.Run("chekUser", func(t *testing.T) {
		addUser(t, app, tokenResponce.Token)
		chekUser(t, app, tokenResponce.Token)
	})

	t.Run("chekNotes", func(t *testing.T) {
		addNote(t, app, tokenResponce.Token)
		chekNote(t, app, tokenResponce.Token)
	})

}

func addUser(t *testing.T, app *Gowek, token string) {
	jsonData := map[string]string{"Login": "testuser", "Email": "testemail", "Hash": "testpass"}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonBytes))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	//при уменьшении таймаута валится ошибка!!!
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	require.Equalf(t, 201, resp.StatusCode, "user created")
}

func chekUser(t *testing.T, app *Gowek, token string) {

	req := httptest.NewRequest("GET", "/api/users/testuser", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var user repo.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "testuser", user.Login)
}

func addNote(t *testing.T, app *Gowek, token string) {
	jsonData := map[string]any{"Date": time.Now().Format("2006-01-02T15:04:05Z07:00"), "Rating": 5, "Note": "testemail"}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/notes", bytes.NewBuffer(jsonBytes))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	//при уменьшении таймаута валится ошибка!!!
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	require.Equalf(t, 201, resp.StatusCode, "Note created")
}

func chekNote(t *testing.T, app *Gowek, token string) {

	req := httptest.NewRequest("GET", "/api/notes/1", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var note repo.Note
	err = json.NewDecoder(resp.Body).Decode(&note)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, uint(0), note.User_id)
}
