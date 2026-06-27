package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"kernelscope/collector"
	"kernelscope/collector/timeline"
	"kernelscope/storage"
)

var Timeline = timeline.NewTimelineEngine(5000)
var PreviousSnapshot timeline.ProcessSnapshot
var HasSnapshot bool

func StartHTTPServer() {
	// 1. Mount the explicit target process endpoint
	http.HandleFunc("/api/processes", handleProcesses)

	// 2. Mount the new timeline route
	http.HandleFunc("/api/timeline", handleTimeline)

	// 3. Global fallback route for routing diagnostics
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("➡️ !!! HTTP API HANDLER RECEIVED A FETCH REQUEST !!! ⬅️")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "KernelScope Routing Error: Path '%s' is invalid.", r.URL.Path)
	})

	// Binding to ":8080" maps globally across the WSL bridge to 0.0.0.0
	fmt.Println("🚀 KernelScope Engine Online -> Listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server crash exception:", err)
	}
}

func handleProcesses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("➡️ !!! HTTP API HANDLER RECEIVED A FETCH REQUEST !!! ⬅️")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	processes, err := collector.CollectProcesses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	currentSnapshot := timeline.NewProcessSnapshot(processes)

	if HasSnapshot {
		events := timeline.DiffProcesses(PreviousSnapshot, currentSnapshot)

		for _, event := range events {
			// Add to the in-memory Ring Buffer ring cache
			Timeline.Add(event)

			// Persist the event to SQLite on disk
			if err := storage.SaveEvent(event); err != nil {
				fmt.Println("❌ Failed to save event:", err)
			}
		}
	}

	PreviousSnapshot = currentSnapshot
	HasSnapshot = true

	json.NewEncoder(w).Encode(processes)
}

func handleTimeline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	events := Timeline.List()

if len(events) == 0 {
	storedEvents, err := storage.LoadEvents(200)
	if err == nil {
		events = storedEvents
	}
}

json.NewEncoder(w).Encode(events)
}
