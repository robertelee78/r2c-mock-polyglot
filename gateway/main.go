// snapstore is a deliberately vulnerable polyglot service used to exercise
// r2c against a small, client-shaped repository. The Go gateway is the
// primary analysis unit and invokes the C, Rust, and Python sibling units.
package main

import (
	"example.com/snapstore/gateway/internal/version"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/import":
		// CWE-918: an unauthenticated caller controls the complete outbound URL.
		sourceURL := r.URL.Query().Get("url")
		response, err := http.Get(sourceURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer response.Body.Close()
		_, _ = io.Copy(w, response.Body)

	case "/trim":
		// Direct argv execution is intentional: it creates a typed boundary to
		// the C unit without introducing a second shell-injection finding here.
		output, err := exec.Command("exiftrim", r.URL.Query().Get("file")).CombinedOutput()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(output)

	case "/jobs":
		output, err := exec.Command("worker", r.URL.Query().Get("name")).CombinedOutput()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(output)

	case "/reports":
		output, err := exec.Command(
			"python3", "report/report.py", r.URL.Query().Get("owner"),
		).CombinedOutput()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(output)

	default:
		http.Error(w, "unknown endpoint", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Printf("snapstore %s listening on http://127.0.0.1:8080\n", version.Value)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
