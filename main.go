package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"validator/helper"
)

func main() {
	// Serve static files from the "static" directory
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Handle the POST request from the button
	http.HandleFunc("/submit", handleSubmit)

	// Start the server on port 8080
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

// handleSubmit is called when the button is clicked
func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Read the input text from the POST request
		inputText := r.FormValue("text")

		// Send a response back to the client
		if inputText != "" {
			fmt.Fprint(w, helper.CheckText(inputText))
		} else {
			fmt.Fprint(w, "You have to provide the text to check")
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}
