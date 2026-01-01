package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"
)

// TerrorsClient est un client simple pour envoyer des erreurs à Terrors
type TerrorsClient struct {
	BaseURL string
	AppID   string
}

// ErrorRequest représente une requête d'erreur
type ErrorRequest struct {
	AppID       string `json:"appId"`
	Message     string `json:"message"`
	Stack       string `json:"stack"`
	Fingerprint string `json:"fingerprint"`
	URL         string `json:"url"`
	Timestamp   int64  `json:"ts"`
	Type        string `json:"type"`
}

// NewTerrorsClient crée un nouveau client Terrors
func NewTerrorsClient(baseURL, appID string) *TerrorsClient {
	return &TerrorsClient{
		BaseURL: baseURL,
		AppID:   appID,
	}
}

// computeFingerprint calcule le fingerprint d'une erreur
func (c *TerrorsClient) computeFingerprint(message, stack string) string {
	// Prendre les premières lignes du stack pour le fingerprint
	stackLines := ""
	if stack != "" {
		lines := bytes.Split([]byte(stack), []byte("\n"))
		if len(lines) > 1 {
			stackLines = string(lines[1]) // Deuxième ligne du stack
		}
	}

	raw := message + "\n" + stackLines
	hash := sha1.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])
}

// CaptureError capture une erreur et l'envoie à Terrors
func (c *TerrorsClient) CaptureError(err error, url string) error {
	message := err.Error()
	stack := string(debug.Stack())

	fingerprint := c.computeFingerprint(message, stack)

	errorReq := ErrorRequest{
		AppID:       c.AppID,
		Message:     message,
		Stack:       stack,
		Fingerprint: fingerprint,
		URL:         url,
		Timestamp:   time.Now().UnixMilli(),
		Type:        "error",
	}

	jsonData, err := json.Marshal(errorReq)
	if err != nil {
		return fmt.Errorf("erreur sérialisation JSON: %w", err)
	}

	resp, err := http.Post(c.BaseURL+"/jason", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erreur envoi HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erreur serveur: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// CapturePanic capture un panic et l'envoie à Terrors
func (c *TerrorsClient) CapturePanic(url string) {
	if r := recover(); r != nil {
		var err error
		switch v := r.(type) {
		case error:
			err = v
		case string:
			err = fmt.Errorf("%s", v)
		default:
			err = fmt.Errorf("%v", v)
		}

		// Envoyer l'erreur (en arrière-plan, ne pas bloquer)
		go func() {
			if err := c.CaptureError(err, url); err != nil {
				fmt.Printf("❌ Erreur lors de l'envoi à Terrors: %v\n", err)
			}
		}()
	}
}

// Exemple d'utilisation
func main() {
	// Initialiser le client
	client := NewTerrorsClient("http://localhost:3000", "app_xxxxxxxx")

	// Exemple 1: Capturer une erreur simple
	fmt.Println("Exemple 1: Erreur simple")
	err := fmt.Errorf("database connection failed: connection timeout")
	if err := client.CaptureError(err, "http://localhost:8080/api/users"); err != nil {
		fmt.Printf("Erreur: %v\n", err)
	}

	// Exemple 2: Utiliser avec defer pour capturer les panics
	fmt.Println("\nExemple 2: Protection avec defer")
	func() {
		defer client.CapturePanic("http://localhost:8080/api/process")

		// Code qui peut paniquer
		panic("something went wrong!")
	}()

	// Exemple 3: Middleware pour un serveur HTTP
	fmt.Println("\nExemple 3: Middleware HTTP")
	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		defer client.CapturePanic(r.URL.String())

		// Votre code ici
		// Si une erreur se produit, elle sera capturée
		panic("test error")
	})

	fmt.Println("✅ Exemples prêts !")
	fmt.Println("Pour tester, lancez le serveur Terrors et exécutez ce programme.")
}
