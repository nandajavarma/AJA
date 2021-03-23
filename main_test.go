package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)

	if err != nil{
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(handler)

	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v instead of %v", status, http.StatusOK)

	}

	expected := "Hello, World!"

	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("handler returned wrong body: got %v instead of %v", actual, expected)
	}
}


func TestRouter(t *testing.T) {

	r := newRouter()

	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/app")

	if err != nil{
		t.Fatal(err)
	}

	if  resp.StatusCode != http.StatusOK {
		t.Errorf("wrong status code: got %v instead of %v", resp.StatusCode, http.StatusOK)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		t.Fatal(err)
	}
	respString := string(b)
	expected := "Hello, World!"

	if respString != expected {
		t.Errorf("handler returned wrong body: got %v instead of %v", respString, expected)
	}
}
