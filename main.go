package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const lifetime time.Duration = 24 * time.Hour
const httpAddr = ":8080"

var devices struct {
	sync.Mutex
	d []Device
}

type Device struct {
	ExternalAddress string    `json:"-"`
	InternalAddress string    `json:"internaladdress"`
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Added           time.Time `json:"-"`
}

func main() {
	devices.d = make([]Device, 0)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Header.Get("Content-Type"))
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Please send json", 400)
			return
		}

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		var t struct {
			Id      string `json:"id"`
			Name    string `json:"name"`
			Address string `json:"address"`
		}

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		log.Println(t.Name)
	})

	http.HandleFunc("/register_get", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")
		ia := r.FormValue("address")
		if ia == "" {
			http.Error(w, `missing "address" URL parameter`, http.StatusBadRequest)
			return
		}

		if net.ParseIP(ia) == nil {
			http.Error(w, `"address" parameter is not a valid addresss`, http.StatusBadRequest)
			return
		}

		ea, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		devices.Lock()
		defer devices.Unlock()

		if i, ok := findDevice(ia, ea); ok {
			devices.d[i].Id = id
			devices.d[i].Name = name
			devices.d[i].Added = time.Now()
		} else {
			devices.d = append(devices.d, Device{
				ExternalAddress: ea,
				InternalAddress: ia,
				Id:              id,
				Name:            name,
				Added:           time.Now(),
			})
		}
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/devices.json", func(w http.ResponseWriter, r *http.Request) {
		ea, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		devices.Lock()
		defer devices.Unlock()

		ds := devicesFor(ea)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ds); err != nil {
			panic(err)
		}
	})

	http.Handle("/devices/", http.StripPrefix("/devices/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ea, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		devices.Lock()
		defer devices.Unlock()

		ds := devicesFor(ea)
		if len(ds) > 0 {
			http.Redirect(w, r, ds[0].InternalAddress, 302)
		} else {
			http.Redirect(w, r, "/devices", 302)
		}
	})

	go cleanup()

	fmt.Println("listen on", httpAddr)
	// TODO: use http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
	log.Fatal(http.ListenAndServe(httpAddr, nil))

}

func findDevice(ia string, ea string) (int, bool) {
	for i, d := range devices.d {
		if d.InternalAddress == ia && d.ExternalAddress == ea {
			return i, true
		}
	}
	return -1, false
}

func devicesFor(ea string) []Device {
	found := []Device{}
	for _, d := range devices.d {
		if d.ExternalAddress == ea {
			found = append(found, d)
		}
	}
	return found
}

func cleanup() {
	for {
		time.Sleep(time.Second * 5)
		devices.Lock()
		for i := len(devices.d) - 1; i >= 0; i-- {
			d := devices.d[i]
			if time.Since(d.Added) > lifetime {
				devices.d = append(devices.d[:i], devices.d[i+1:]...)
			}
		}
		devices.Unlock()
	}
}
