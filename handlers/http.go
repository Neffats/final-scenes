package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Neffats/final-scenes/models"
	"github.com/Neffats/final-scenes/stores"
)

type HTTP struct {
	Films  *stores.FilmStore
	Logger *log.Logger
}

func (h *HTTP) HandleGuess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.Printf("received bad guess request: unsupported method: %s\n", r.Method)
		http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		h.Logger.Printf(
			"received bad guess request: unsupported Content-Type: %s\n", contentType)
		http.Error(w, "Unsupported Content-Type", http.StatusBadRequest)
		return
	}

	var guess models.GuessAttempt
	var resp models.GuessResponse

	err := json.NewDecoder(r.Body).Decode(&guess)
	if err != nil {
		h.Logger.Printf("failed to unmarshal guess: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	hashedGuess := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.ToLower(guess.Guess))))
	if hashedGuess == guess.Question {
		resp.Answer = true
	} else {
		resp.Answer = false
	}

	byteResp, err := json.Marshal(resp)
	if err != nil {
		h.Logger.Printf("failed to marshal guess response: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(byteResp)
	return
}

func (h *HTTP) HandleTemplate(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.gohtml").Funcs(template.FuncMap{
		"inc": func(x int) int {
			return x + 1
		},
	}).ParseFiles("templates/index.gohtml")
	if err != nil {
		h.Logger.Printf("failed to parse template file: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	films := h.Films.All()
	err = t.Execute(w, films)
	if err != nil {
		h.Logger.Printf("failed to execute template file: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
