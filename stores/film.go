package stores

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/Neffats/final-scenes/models"
)

type FilmNotFoundError struct {
	Message string
}

func (e *FilmNotFoundError) Error() string {
	return fmt.Sprintf("film not found: %s", e.Message)
}

type FilmStore struct {
	Films  []models.Film `json:"films"`
	source string
}

func NewFilmStore(filename string) *FilmStore {
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

type findFunc func(film models.Film) bool

func (s *FilmStore) Find(match findFunc) (models.Film, error) {
	for _, film := range s.Films {
		if match(film) {
			return film, nil
		}
	}
	return models.Film{}, &FilmNotFoundError{Message: "film not found in store"}
}

func ByHash(hash string) findFunc {
	return func(film models.Film) bool {
		return film.Hash == hash
	}
}

func ByName(name string) findFunc {
	return func(film models.Film) bool {
		return film.Name == name
	}
}
