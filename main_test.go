package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestAddingDevice(t *testing.T) {
	go main()
	time.Sleep(time.Second * 2)
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://localhost:8080/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	finalURL := resp.Request.URL.String()
	if finalURL != "http://example.org" {
		t.Errorf("Redirect to fallback did not work.")
	}

	// add one device
	resp, err = client.Get("http://localhost:8080/register?id=2323&name=device&address=192.168.100.151")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	resp, err = client.Get("http://localhost:8080/list.json")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	list, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var devices []Device
	if err := json.Unmarshal(list, &devices); err != nil {
		panic(err)
	}

	if len(devices) != 1 {
		t.Fatal("No device added to list.")
	}

	// update first device
	resp, err = client.Get("http://localhost:8080/register?id=232323&name=device23&address=192.168.100.151")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	resp, err = client.Get("http://localhost:8080/list.json")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	list, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err = json.Unmarshal(list, &devices); err != nil {
		panic(err)
	}

	if len(devices) != 1 {
		t.Fatal("No device added to list.")
	}

	if devices[0].Name != "device23" {
		t.Fatal("Name not updated.")
	}

	// add second
	resp, err = client.Get("http://localhost:8080/register?id=23232&name=device2&address=192.168.100.152")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	resp, err = client.Get("http://localhost:8080/list.json")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	list, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err = json.Unmarshal(list, &devices); err != nil {
		panic(err)
	}

	if len(devices) != 2 {
		t.Fatal("No second device added to list.")
	}

	// check if it redirects me to a device, redirect should fail
	resp, err = client.Get("http://localhost:8080/")
	if err == nil {
		t.Fatal(err)
	}

}
