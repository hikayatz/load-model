package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	. "github.com/hikayatz/load-model/helpers"

)

func TestPostPerson(t *testing.T) {

	var jsonReq = `{
    "name": "hikayat", 
    "address": "jakarta pusat", 
    "age" : "28",
    "is_married": "true",
     "weight": "58.9",
     "birth_day": "13-08-1992"
}`
	var jsonByte = []byte(jsonReq)
	req, err := http.NewRequest("POST", "/person", bytes.NewBuffer(jsonByte))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(HeaderContentType, MIMEApplicationJSON)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePostPerson)
	handler.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("there are something erros with status code: got %v want %v",
			status, http.StatusOK)
	}
	var expected = `{"name":"hikayat","address":"jakarta pusat","phone":"","age":28,"is_married":true,"weight":58.9,"birth_day":"1992-08-13T00:00:00Z"}`
	// Check the response body
	if res.Body.String() != expected{
		t.Errorf("response body not unexpected, response body is %v and expected response %v", res.Body.String(), expected )
	}

}
