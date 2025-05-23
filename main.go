package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"validator/processor"
	"validator/stemmer"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		return
	}
	switch mode {
	case "web":
		startWebServer()
	case "tg":
		startTgBot()
	case "all":
		go startWebServer()
		go startTgBot()
		select {}
	default:
		fmt.Println("Unknown mode")
	}
}

func startTgBot() {
	bot, err := tgbotapi.NewBotAPI("7841984143:AAGcRRzW1Nsdy4yy6yUvUwgWWVioBgjDf9E")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true // Set to true for debugging purposes

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Respond to the message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, processor.CheckText(update.Message.Text))
			msg.ReplyToMessageID = update.Message.MessageID

			// Send the message
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func startWebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Handle the POST request from the button
	// http.HandleFunc("/submit", handleSubmit)
	http.HandleFunc("/checker", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/checker.html")
	})

	http.HandleFunc("/check", handleCheck)

	http.HandleFunc("/stemmer", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/stemmer.html")
	})

	http.HandleFunc("/stem", handleStem)

	// Start the server on port 8080
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// handleSubmit is called when the button is clicked
func handleCheck(w http.ResponseWriter, r *http.Request) {
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

func handleStem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Read the input text from the POST request
		inputText := r.FormValue("text")
		language := r.FormValue("language")
		language = strings.ToLower(language)

		// Send a response back to the client
		if inputText != "" {
			fmt.Fprint(w, stemmer.StemText(inputText, language))
		} else {
			fmt.Fprint(w, "No glossary provided")
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
