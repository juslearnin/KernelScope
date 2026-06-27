package main

import (
	"kernelscope/server"
	"kernelscope/storage"
	"log"
)

func main() {
err := storage.InitDatabase()
if err != nil {
	log.Fatal(err)
}
	server.StartHTTPServer()
}

