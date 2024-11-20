package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"telegram_bot/helpers"
	"time"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	botKey := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			handleCommand(&update, bot)

		}
	}
}

func handleCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	messageText := update.Message.Text
	chatID := update.Message.Chat.ID
	args := strings.Fields(messageText)
	command := strings.ToLower(args[0])

	if strings.HasPrefix(command, "/translate") {
		helpers.TranslateCommand(bot, chatID, args)
		return
	}
	switch command {
	case "/start":
		startCommand(bot, chatID, update.Message.From.FirstName)

	case "/help":
		helpCommand(bot, chatID)

	case "/joke":
		helpers.JokeCommand(bot, chatID)

	case "/weather":
		if len(args) < 2 {
			bot.Send(tgbotapi.NewMessage(chatID, "Please specify a city, e.g., /weather Lagos"))
			return
		}
		city := args[1]
		weatherCommand(bot, chatID, city)

	case "/time":
		timeCommand(bot, chatID)

	case "/ip":
		helpers.IpCommand(bot, chatID, args, update)

	default:
		echoCommand(bot, chatID, messageText)
	}
}

func startCommand(bot *tgbotapi.BotAPI, chatID int64, firstName string) {
	message := fmt.Sprintf("Hello, %s! Welcome to Golang Telegram Bot. Type /help to see what I can do.", firstName)
	bot.Send(tgbotapi.NewMessage(chatID, message))
}

func helpCommand(bot *tgbotapi.BotAPI, chatID int64) {
	helpText := `Here are the commands you can use:
/start - Start the bot
/help - List available commands
/joke - Get a random joke
/weather <city> - Get the current weather for a specified city
/time - Get the current date and time
/translate <target_lang> <text> - Translate text to a specified language
   Example: /translate es Hello, how are you?
   Supported languages: en, es, fr, de, it, ru, zh, ar
/ip <ip_address> - Get information about a specific IP address
   Example: /ip 192.168.1.1
Type anything else, and I will echo it back to you!`
	bot.Send(tgbotapi.NewMessage(chatID, helpText))
}

func weatherCommand(bot *tgbotapi.BotAPI, chatID int64, city string) {
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	weather, err := helpers.GetWeather(city, weatherAPIKey)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Sorry, I couldn't get the weather information."))
		log.Println(err)
		return
	}
	bot.Send(tgbotapi.NewMessage(chatID, weather))
}

func timeCommand(bot *tgbotapi.BotAPI, chatID int64) {
	currentTime := time.Now().Format("Monday, 02-Jan-2006 15:04:05 MST")
	bot.Send(tgbotapi.NewMessage(chatID, "Current time: "+currentTime))
}

func echoCommand(bot *tgbotapi.BotAPI, chatID int64, message string) {
	bot.Send(tgbotapi.NewMessage(chatID, "You said: "+message))
}
