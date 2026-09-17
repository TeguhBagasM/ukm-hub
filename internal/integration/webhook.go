package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// SendFormWebhook mengirim payload registrasi ke webhook secara asynchronous.
// Kegagalan webhook tidak pernah menggagalkan pembuatan registrasi.
func SendFormWebhook(url string, payload any) {
	if url == "" {
		return
	}
	go func() {
		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[webhook] failed to marshal payload: %v", err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			log.Printf("[webhook] failed to build request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[webhook] request to %s failed: %v", url, err)
			return
		}
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("[webhook] delivered to %s (status %d)", url, resp.StatusCode)
			return
		}
		log.Printf("[webhook] %s returned status %d: %s", url, resp.StatusCode, string(respBody))
	}()
}
