package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const lifetime time.Duration = 24 * time.Hour
const httpAddr = ":8180"

var devices struct {
	sync.Mutex
	d []Device
}

type Device struct {
	ExternalAddress string    `json:"-"`
	InternalAddress string    `json:"internaladdress"`
	Port            int       `json:"port,omitempty"` // optional
	Name            string    `json:"name"`
	Added           time.Time `json:"added"`
}

func main() {
	devices.d = make([]Device, 0)

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/api/register", RegisterDevice)
	http.HandleFunc("/api/devices", ListDevices)
	http.Handle("/", http.FileServer(http.Dir("public")))

	go cleanup()

	fmt.Println("listen on", httpAddr)
	// Note: use TLS
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

func RegisterDevice(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Please send json", 400)
		return
	}

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var t struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Port    int    `json:"port"`
	}

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	t.Address = strings.Trim(t.Address, " ")

	// Prevent simple loopback mistake
	if t.Address == "127.0.0.1" || t.Address == "::1" {
		http.Error(w, `Loopback is not allowed`, http.StatusBadRequest)
		return
	}

	if net.ParseIP(t.Address) == nil {
		http.Error(w, `"address" is not a valid IP address`, http.StatusBadRequest)
		return
	}

	// TODO: validate parameter name required and no html/js
	ea, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Check if proxy was configured.
	if ea == "127.0.0.1" {
		xrealip := r.Header.Get("x-real-ip")
		if xrealip != "" {
			ea = xrealip
		} else {
			log.Println("127.0.0.1 tried to add an address, this can happen when proxy is not configured correctly.")
			http.Error(w, `Host 127.0.0.1 is not allowed to register devices`, http.StatusBadRequest)
			http.NotFound(w, r)
			return
		}
	}

	devices.Lock()
	defer devices.Unlock()

	if i, ok := findDevice(t.Address, ea); ok {
		devices.d[i].Name = t.Name
		devices.d[i].Port = t.Port
		devices.d[i].Added = time.Now()
		log.Println(time.Now(), "updated", t.Address)
	} else {
		devices.d = append(devices.d, Device{
			ExternalAddress: ea,
			InternalAddress: t.Address,
			Port:            t.Port,
			Name:            t.Name,
			Added:           time.Now(),
		})
		log.Println(time.Now(), "added", t.Address)
	}

	fmt.Fprintf(w, "Successfully added, visit https://nupnp.com for more.\n")
}

func ListDevices(w http.ResponseWriter, r *http.Request) {
	ea, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Check if proxy was configured.
	if ea == "127.0.0.1" {
		xrealip := r.Header.Get("x-real-ip")
		if xrealip != "" {
			ea = xrealip
		} else {
			log.Println("127.0.0.1 tried to access an address, this can happen when proxy is not configured correctly.")
			http.NotFound(w, r)
			return
		}
	}

	devices.Lock()
	defer devices.Unlock()

	ds := devicesFor(ea)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ds); err != nil {
		panic(err)
	}
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
