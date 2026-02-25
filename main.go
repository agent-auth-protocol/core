package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// In-memory store for V1 (To be replaced with Redis/Postgres in production)
var (
	agentRegistry = make(map[string]ed25519.PublicKey)
	mu            sync.RWMutex
	// Server's private key for signing the JWTs (Simulated for V1)
	_, serverPrivateKey, _ = ed25519.GenerateKey(nil)
)

type RegisterRequest struct {
	AgentID   string `json:"agent_id"`
	PublicKey string `json:"public_key_hex"` // Ed25519 public key in hex
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	pubKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil || len(pubKeyBytes) != ed25519.PublicKeySize {
		http.Error(w, "Invalid public key", http.StatusBadRequest)
		return
	}

	mu.Lock()
	agentRegistry[req.AgentID] = pubKeyBytes
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Agent %s registered successfully for M2M auth.\n", req.AgentID)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	agentID := r.Header.Get("X-Agent-ID")

	mu.RLock()
	_, exists := agentRegistry[agentID]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Agent not registered", http.StatusUnauthorized)
		return
	}

	// In a full implementation, we would verify a cryptographic signature payload here.
	// For this V1 prototype, we issue a 5-minute JWT upon recognizing the AgentID.

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub": agentID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		"aud": "agent-infrastructure",
	})

	tokenString, err := token.SignedString(serverPrivateKey)
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": tokenString, "expires_in": "300"})
}

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/token", tokenHandler)

	fmt.Println("🛡️ AgentAuth-Core running on :8080")
	fmt.Println("Architecting Trust for the Autonomous Economy.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
