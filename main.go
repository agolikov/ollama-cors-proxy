package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Get Ollama URL from environment variables or construct from host and port
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "localhost"
	}

	ollamaPort := os.Getenv("OLLAMA_PORT")
	if ollamaPort == "" {
		ollamaPort = "11434"
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://" + ollamaHost + ":" + ollamaPort
	}

	// Get proxy port from environment variable or use default
	proxyPort := os.Getenv("PROXY_PORT")
	if proxyPort == "" {
		proxyPort = "11435"
	}

	targetURL, err := url.Parse(ollamaURL)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify the request before forwarding it
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Forward original headers to Ollama
		req.Header.Set("Host", targetURL.Host)
		req.Header.Set("Origin", targetURL.Scheme+"://"+targetURL.Host)
		req.Host = targetURL.Host
	}

	// Set CORS headers in the response
	// Add this type before the main function
	handleWithCORS := func(w http.ResponseWriter, r *http.Request) {
		// Create wrapped response writer to capture status
		sw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}

		// Log incoming request
		log.Printf("📨 %s %s://%s%s %s", r.Method, r.URL.Scheme, r.Host, r.URL.Path, r.RemoteAddr)

		// CORS headers
		sw.Header().Set("Access-Control-Allow-Origin", "*")
		sw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		sw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Content-Length, Accept-Encoding, X-CSRF-Token, Origin, Cache-Control, X-Requested-With")
		sw.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		sw.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			log.Printf("✈️  CORS Preflight request")
			sw.WriteHeader(http.StatusOK)
			return
		}

		proxy.ServeHTTP(sw, r)

		// Log response status with status text
		statusText := http.StatusText(sw.status)
		log.Printf("📫 [%d %s] %s %s", sw.status, statusText, r.Method, r.URL.Path)
	}

	http.HandleFunc("/", handleWithCORS)

	log.Printf("🚀 CORS Proxy listening on :%s -> %s", proxyPort, ollamaURL)
	log.Fatal(http.ListenAndServe(":"+proxyPort, nil))
}
