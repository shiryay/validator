package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"validator/processor"
)

func main() {
	var mode string

	fs := flag.NewFlagSet("mode", flag.ExitOnError)
	fs.StringVar(&mode, "m", "", "server type")
	if len(os.Args) < 3 {
		fs.Usage()
		return
	}
	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing flag: ", err)
	} else if mode == "web" {
		fmt.Println("Running in web mode")
	} else if mode == "tg" {
		fmt.Println("Running as telegram bot")
	} else {
		fmt.Println("Unknown mode")
	}
}

func setupWebServer() {
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
			fmt.Fprint(w, processor.CheckText(inputText))
		} else {
			fmt.Fprint(w, "You have to provide the text to check")
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
