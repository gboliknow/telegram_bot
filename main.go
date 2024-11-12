package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	botKey := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	jokes := []string{
		"Why did the Golang developer get fired? Because he couldn’t C!",
		"Concurrency jokes are not funny if no one’s listening.",
		"Why do Go developers carry umbrellas? Because it drizzles!",
	}

	rand.Seed(time.Now().UnixNano())

	for update := range updates {
		if update.Message != nil {
			switch {
			case update.Message.Text == "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, "+update.Message.From.FirstName+"! Welcome to Golang Telegram Bot. Type /help to see available commands.")
				bot.Send(msg)

			case update.Message.Text == "/help":
				helpText := "Here are the commands you can use:\n" +
					"/start - Start the bot\n" +
					"/help - List available commands\n" +
					"/joke - Get a random joke\n" +
					"/weather <city> - Get the current weather for a specified city\n" +
					"/time - Get the current date and time\n" +
					"Type anything else, and I will echo it back to you!"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
				bot.Send(msg)

			case update.Message.Text == "/joke":
				joke := jokes[rand.Intn(len(jokes))] // Pick a random joke
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, joke)
				bot.Send(msg)

			case strings.HasPrefix(update.Message.Text, "/weather"):
				parts := strings.Split(update.Message.Text, " ")
				if len(parts) < 2 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please specify a city, e.g., /weather Lagos")
					bot.Send(msg)
					continue
				}

				city := parts[1]
				weather, err := getWeather(city, weatherAPIKey)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I couldn't get the weather information.")
					bot.Send(msg)
					log.Println(err)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, weather)
					bot.Send(msg)
				}

			case update.Message.Text == "/time":
				currentTime := getTime()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, currentTime)
				bot.Send(msg)

			default:

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You said: "+update.Message.Text)
				bot.Send(msg)
			}
		}
	}
}

func getTime() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("Monday, 02-Jan-2006 15:04:05 MST")
	return fmt.Sprintf("Current date and time: %s", formattedTime)
}

func getWeather(city string, key string) (string, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, key)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get weather data")
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	main := data["main"].(map[string]interface{})
	weather := data["weather"].([]interface{})[0].(map[string]interface{})
	temp := main["temp"].(float64)
	description := weather["description"].(string)

	return fmt.Sprintf("The weather in %s: %0.1f°C, %s.", city, temp, description), nil
}
