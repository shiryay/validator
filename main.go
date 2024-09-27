package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"validator/processor"

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
	// Serve static files from the "static" directory
	http.Handle("/", http.FileServer(http.Dir("/home/samsepiol/go/src/validator/static")))

	// Handle the POST request from the button
	http.HandleFunc("/submit", handleSubmit)

	// Start the server on port 8080
	fmt.Println("Server started on port 8080")
	// http.ListenAndServe(":8080", nil)
	http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/n0access.freemyip.com/fullchain.pem", "/etc/letsencrypt/live/n0access.freemyip.com/privkey.pem", nil)
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
