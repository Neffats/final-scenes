package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	logger = log.New(os.Stdout, "logger: ", log.Ldate | log.Ltime | log.Lshortfile)
)


func main() {
	port := os.Getenv("FINAL_SCENES_PORT")
	if port == "" {
		logger.Fatal("FINAL_SCENES_PORT environment variable missing.")
	}
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./www"))
	mux.HandleFunc("/templates/", LogWrapperHF(HandleTemplate))
	mux.Handle("/", LogWrapper(fs))

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
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func LogWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}

func LogWrapperHF(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}
type GuessAttempt struct {
	Question string `json:"question"`
	Guess    string `json:"guess"`
}

func HandleGuess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		return
	}

	var guess GuessAttempt
	err := json.NewDecoder(r.Body).Decode(&guess)
	if err != nil {
		logger.Printf("failed to unmarshal guess: %v\n", err)
		http.Error(w, "something went wring", http.StatusInternalServerError)
	}
}

type FinalScene struct {
	Name string
	AudioFile string
	Year string
	ImageFile string
	Hash string
}

func HandleTemplate(w http.ResponseWriter, r *http.Request) {
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

	scenes := make([]FinalScene, 0)
	casablanca := FinalScene{
		Name: "Casablanca",
		AudioFile: "audio/sound1.wav",
		Year: "1942",
		ImageFile: "images/picture1.png",
		Hash: "Casablanca",
	}
	psycho := FinalScene{
		Name: "Psycho",
		AudioFile: "audio/sound2.wav",
		Year: "1960",
		ImageFile: "images/picture2.png",
		Hash: "Psycho",
	}

	scenes = append(scenes, casablanca)
	scenes = append(scenes, psycho)

	err = t.Execute(w, scenes)
	if err != nil {
		logger.Printf("failed to execute template file: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
