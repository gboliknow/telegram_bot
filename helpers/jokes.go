package helpers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// JokeResponse represents the structure of joke API response
type JokeResponse struct {
    Type       string `json:"type"`
    Setup      string `json:"setup"`
    Punchline string `json:"punchline"`
}

// Advanced joke fetching function with multiple sources and error handling
func JokeCommand(bot *tgbotapi.BotAPI, chatID int64) {
    // Array of joke APIs to try
    jokeAPIs := []string{
        "https://official-joke-api.appspot.com/jokes/random",
        "https://v2.jokeapi.dev/joke/Programming?type=twopart",
    }

    for _, apiURL := range jokeAPIs {
        joke, err := fetchJokeFromAPI(apiURL)
        if err == nil {
            // Format joke based on API response type
            var jokeMessage string
            switch {
            case joke.Setup != "" && joke.Punchline != "":
                jokeMessage = fmt.Sprintf("%s\n\n*Punchline*: %s", joke.Setup, joke.Punchline)
            case joke.Punchline != "":
                jokeMessage = joke.Punchline
            default:
                continue
            }

            // Send joke with some formatting
            msg := tgbotapi.NewMessage(chatID, jokeMessage)
            msg.ParseMode = "Markdown"
            bot.Send(msg)
            return
        }
    }

    // Fallback to local jokes if all API calls fail
    localJokes := []string{
        "Why did the Golang developer get fired? Because he couldn't C!",
        "Concurrency jokes are not funny if no one's listening.",
        "Why do Go developers carry umbrellas? Because it drizzles!",
    }
    
    msg := tgbotapi.NewMessage(chatID, localJokes[time.Now().Unix()%int64(len(localJokes))])
    bot.Send(msg)
}

// Fetch joke from a specific API with timeout and error handling
func fetchJokeFromAPI(apiURL string) (JokeResponse, error) {
    // Create a client with timeout
    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    // Make the request
    resp, err := client.Get(apiURL)
    if err != nil {
        return JokeResponse{}, err
    }
    defer resp.Body.Close()

    // Check HTTP status code
    if resp.StatusCode != http.StatusOK {
        return JokeResponse{}, fmt.Errorf("API returned status %d", resp.StatusCode)
    }

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return JokeResponse{}, err
    }

    // Parse JSON response
    var joke JokeResponse
    err = json.Unmarshal(body, &joke)
    if err != nil {
        return JokeResponse{}, err
    }

    return joke, nil
}