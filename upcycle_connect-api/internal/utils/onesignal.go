package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendPushNotification(playerID, title, message string) {
	if playerID == "" {
		return
	}

	appID := os.Getenv("ONESIGNAL_APP_ID")
	apiKey := os.Getenv("ONESIGNAL_API_KEY")
	if appID == "" || apiKey == "" {
		return
	}

	payload := map[string]any{
		"app_id":             appID,
		"include_player_ids": []string{playerID},
		"headings":           map[string]string{"en": title, "fr": title},
		"contents":           map[string]string{"en": message, "fr": message},
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("OneSignal request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("OneSignal send error:", err)
		return
	}
	defer resp.Body.Close()
}
