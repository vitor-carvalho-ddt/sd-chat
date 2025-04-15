package web

import (
	"fmt"
	"net/http"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Accessing URL: %s\n", r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "frontend/index.html")
}
