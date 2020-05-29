package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/Neffats/final-scenes/models"
)

var (
	logger = log.New(os.Stdout, "logger: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	store := NewStore("films.json")
	err := store.Init()
	if err != nil {
		logger.Fatalf("failed to initialise film store: %v", err)
	}

	h := &HTTPHandler{
		Films: store,
	}

	port := os.Getenv("FINAL_SCENES_PORT")
	if port == "" {
		logger.Fatal("FINAL_SCENES_PORT environment variable missing.")
	}
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./www"))
	mux.Handle("/static/", LogWrapper(http.StripPrefix("/static/", fs)))
	mux.Handle("/guess/", LogWrapper(http.HandlerFunc(h.HandleGuess)))
	mux.Handle("/", LogWrapper(http.HandlerFunc(h.HandleTemplate)))

	srv := &http.Server{
		Addr:     fmt.Sprintf("0.0.0.0:%s", port),
		Handler:  mux,
		ErrorLog: logger,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Println("listening....")
	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

type HTTPHandler struct {
	Films *FilmStore
}

func LogWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}

func (h *HTTPHandler) HandleGuess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Printf("received bad guess request: unsupported method: %s\n", r.Method)
		http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		logger.Printf(
			"received bad guess request: unsupported Content-Type: %s\n", contentType)
		http.Error(w, "Unsupported Content-Type", http.StatusBadRequest)
		return
	}

	var guess models.GuessAttempt
	var resp models.GuessResponse

	err := json.NewDecoder(r.Body).Decode(&guess)
	if err != nil {
		logger.Printf("failed to unmarshal guess: %v\n", err)
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
		logger.Printf("failed to marshal guess response: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(byteResp)
	return
}

func (h *HTTPHandler) HandleTemplate(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.gohtml").Funcs(template.FuncMap{
		"inc": func(x int) int {
			return x + 1
		},
	}).ParseFiles("templates/index.gohtml")
	if err != nil {
		logger.Printf("failed to parse template file: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	films := h.Films.All()
	err = t.Execute(w, films)
	if err != nil {
		logger.Printf("failed to execute template file: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}

type FilmStore struct {
	Films  []models.Film `json:"films"`
	source string
}

func NewStore(filename string) *FilmStore {
	return &FilmStore{
		Films:  make([]models.Film, 0),
		source: filename,
	}
}

func (s *FilmStore) Init() error {
	f, err := ioutil.ReadFile(s.source)
	if err != nil {
		return fmt.Errorf("failed to open store source: %v", err)
	}

	err = json.Unmarshal(f, &s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal source data: %v", err)
	}
	return nil
}

func (s *FilmStore) Random() models.Film {
	index := rand.Intn(len(s.Films))
	return s.Films[index]
}

func (s *FilmStore) All() []models.Film {
	return s.Films
}
