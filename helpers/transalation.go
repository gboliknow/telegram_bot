package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var supportedLanguages = map[string]string{
	"en": "English",
	"es": "Spanish",
	"fr": "French",
	"de": "German",
	"it": "Italian",
	"ru": "Russian",
	"zh": "Chinese",
	"ar": "Arabic",
}

type TranslationRequest struct {
	Q      string `json:"q"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type TranslationResponse struct {
	TranslatedText string `json:"translatedText"`
}

func translateText(text, sourceLang, targetLang string) (string, error) {
	// Default to auto-detect if source language not specified
	if sourceLang == "" {
		sourceLang = "auto"
	}

	// Prepare request payload
	reqBody, err := json.Marshal(TranslationRequest{
		Q:      text,
		Source: sourceLang,
		Target: targetLang,
	})
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}

	resp, err := http.Post(
		"https://libretranslate.de/translate",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", fmt.Errorf("error making request to translation API: %v", err)
	}
	defer resp.Body.Close()

	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to translate, status code: %d", resp.StatusCode)
	}

	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var translationResp TranslationResponse
	err = json.Unmarshal(body, &translationResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	return translationResp.TranslatedText, nil
}

func TranslateCommand(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	if len(args) < 3 {
		bot.Send(tgbotapi.NewMessage(chatID, "Usage: /translate <target_lang> <text>\nSupported languages: en, es, fr, de, it, ru, zh, ar"))
		return
	}

	targetLang := args[1]
	text := strings.Join(args[2:], " ")

	if _, exists := supportedLanguages[targetLang]; !exists {
		bot.Send(tgbotapi.NewMessage(chatID, "Unsupported language. Use: en, es, fr, de, it, ru, zh, ar"))
		return
	}

	translatedText, err := translateText(text, "", targetLang)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Translation failed. Please try again."))
		log.Println(err)
		return
	}

	responseMsg := fmt.Sprintf("Translation to %s:\n%s", supportedLanguages[targetLang], translatedText)
	bot.Send(tgbotapi.NewMessage(chatID, responseMsg))
}
