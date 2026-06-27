package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"kernelscope/collector"
)

func StartHTTPServer() {
	// 1. Mount the explicit target process endpoint
	http.HandleFunc("/api/processes", handleProcesses)

	// 2. Global fallback route for routing diagnostics
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

	json.NewEncoder(w).Encode(processes)
}