package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Railway!")
	})

	addr := "0.0.0.0:" + port
	fmt.Println("🚀 Server starting on port", port)
	fmt.Println("🌐 Frontend available at: http://" + addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
