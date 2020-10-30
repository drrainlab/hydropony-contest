package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var userID uint
var testEmail = "test_user@test.com"

func Test_Create_New_User_Should_Be_Status_OK(t *testing.T) {

	router, srv := runServer()
	ts := httptest.NewServer(router)
	defer srv.Close()

	newUser := &User{
		Email:     testEmail,
		FirstName: "Bob",
		LastName:  "Marley",
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(newUser)

	// Create new user
	resp, err := http.Post(fmt.Sprintf("%s/users", ts.URL), "application/json", b)

	if err != nil {
		log.Fatal("Error", err)
	}

	var body struct {
		ID      uint              `json:"id"`
		Status  string            `json:"status"`
		Headers map[string]string `json:"headers"`
		Origin  string            `json:"origin"`
	}
	json.NewDecoder(resp.Body).Decode(&body)

	userID = body.ID

	if body.Status != "OK" {
		t.Fail()
	}

}

func Test_Create_Existing_User_Should_Be_Status_error(t *testing.T) {

	router, srv := runServer()
	ts := httptest.NewServer(router)
	defer srv.Close()

	newUser := &User{
		Email:     testEmail,
		FirstName: "Bob",
		LastName:  "Marley",
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(newUser)

	resp, err := http.Post(fmt.Sprintf("%s/users", ts.URL), "application/json", b)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 409 {
		t.Fail()
	}

}

func Test_Get_User_Info_Should_Return_Correct_Record(t *testing.T) {

	router, srv := runServer()
	ts := httptest.NewServer(router)
	defer srv.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users", ts.URL), nil)
	if err != nil {
		panic(err)
	}
	// trying to create existing user
	req.Header.Set("user-id", strconv.Itoa(int(userID)))

	resp, err := ts.Client().Do(req)
	if err != nil {
		panic(err)
	}

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)

	if err != nil {
		panic(err)
	}

	if user.Email != testEmail {
		t.Fail()
	}

}

func Test_Get_Nonexisting_User_Info_Should_Return_Not_Found(t *testing.T) {

	router, srv := runServer()
	ts := httptest.NewServer(router)
	defer srv.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users", ts.URL), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("user-id", strconv.Itoa(int(99999999)))

	resp, err := ts.Client().Do(req)
	if resp.StatusCode != 404 {
		t.Fail()
	}

}
