package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IPInfo struct {
	IP      string `json:"query"`
	Country string `json:"country"`
	City    string `json:"city"`
	Region  string `json:"regionName"`
	ISP     string `json:"isp"`
}

func fetchIPInfo(ip string) (*IPInfo, error) {
    url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Log the raw response body for debugging
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    log.Printf("Raw API Response: %s", body)

    // Decode the response body into IPInfo struct
    var ipInfo IPInfo
    if err := json.Unmarshal(body, &ipInfo); err != nil {
        return nil, fmt.Errorf("failed to unmarshal IP info: %w", err)
    }

    // Check for missing fields
    if ipInfo.Country == "" || ipInfo.City == "" {
        return nil, fmt.Errorf("IP information not available for %s", ip)
    }

    return &ipInfo, nil
}


func IpCommand(bot *tgbotapi.BotAPI, chatID int64, args []string, update *tgbotapi.Update) {
	var ipAddress string
	if len(args) < 2 {
		// Use the sender's username if no IP address is provided
		ipAddress = update.Message.From.UserName
	} else {
		// Use the IP address provided by the user
		ipAddress = args[1]
	}

	ipInfo, err := fetchIPInfo(ipAddress)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "IP lookup failed."))
		log.Println(err)
		return
	}

	// Send back IP info to the user
	responseMsg := fmt.Sprintf(
		"IP: %s\nCountry: %s\nCity: %s\nRegion: %s\nISP: %s",
		ipInfo.IP, ipInfo.Country, ipInfo.City, ipInfo.Region, ipInfo.ISP)
	bot.Send(tgbotapi.NewMessage(chatID, responseMsg))
}
