package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	body := bytes.NewBufferString("{\"name\":\"Testdevice\",\"address\":\"192.168.100.151\"}")
	req, err := http.NewRequest("POST", "/api/register", body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.RemoteAddr = "80.2.3.41:321"

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterDevice)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v - %v",
			status, rr.Body)
	}

	// Check the response body is what we expect.
	expected := "Successfully added, visit https://nupnp.com for more.\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestList(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/devices", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Fake remoteAddr, does not have to be tested, because user cannot change it.
	// If it is 127.0.0.1, then the server thinks its a proxy.
	req.RemoteAddr = "80.2.3.41:321"

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListDevices)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v - %v",
			status, rr.Body)
	}

	log.Print(rr.Body.String())

}
