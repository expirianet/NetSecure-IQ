package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// MikroTikData represents the structure of the data sent from MikroTik
type MikroTikData struct {
	MAC    string `json:"mac"`
	Status string `json:"status"` // "online" or "offline"
}

// Simulated database of known MikroTik MAC addresses
var knownMACs = map[string]bool{
	"AA:BB:CC:DD:EE:FF": true,
	"11:22:33:44:55:66": true,
}

func main() {
	http.HandleFunc("/api/ping", handlePing)
	fmt.Println("Go Backend is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var data MikroTikData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Verify if the MAC address is known
	if !knownMACs[data.MAC] {
		http.Error(w, "Unknown MAC address", http.StatusUnauthorized)
		return
	}

	fmt.Printf("Received data from %s with status: %s\n", data.MAC, data.Status)

	// Save to InfluxDB
	if err := saveToInflux(data.MAC, data.Status); err != nil {
		log.Printf("InfluxDB write failed: %v\n", err)
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received and stored"))
}

func saveToInflux(mac, status string) error {
	// Adjust these parameters to match your InfluxDB config
	const (
		token  = "my-token"
		bucket = "netsecure"
		org    = "netsecure-org"
		url    = "http://influxdb:8086"
	)

	client := influxdb2.NewClient(url, token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(org, bucket)

	point := influxdb2.NewPoint(
		"device_status",
		map[string]string{"mac": mac},
		map[string]interface{}{"status": status},
		time.Now(),
	)

	return writeAPI.WritePoint(context.Background(), point)
}
