package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	token  = os.Getenv("TELEGRAM_TOKEN")
	chatID = os.Getenv("TELEGRAM_CHAT_ID")
)

func SendItems(items []Item) {
	for _, item := range items {
		txt := fmt.Sprintf(`🔔 *%s* from %s
Mention: %s
URL: %s
Time: %s
`, item.Title, item.Provider, item.Mention, item.URL, item.Timestamp.Format("2006-01-02 15:04:05"))
		sendMessage(txt)
	}
}

func SendError(err error) {
	sendMessage(err.Error())
}

func sendMessage(msg string) {
	body := map[string]any{
		"chat_id": chatID,
		"text":    msg,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Println("error marshalling telegram body:", err)
		return
	}

	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+token+"/sendMessage", bytes.NewReader(bodyBytes))
	if err != nil {
		log.Println("error creating telegram request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := Client.Do(req)
	if err != nil {
		log.Println("error sending telegram message:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		errBody, _ := io.ReadAll(res.Body)
		log.Println("error sending telegram message:", res.Status, "body:", string(errBody))
		return
	}
}
