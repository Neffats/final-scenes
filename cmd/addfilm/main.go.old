package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/Neffats/final-scenes/models"
)

func main() {
	name := flag.String("name", "", "Film name")
	image := flag.String("image", "", "Film image to add")
	audio := flag.String("audio", "", "Film audio to add")
	year := flag.String("year", "", "Film year of release")
	flag.Parse()

	if *name == "" {
		fmt.Println("missing arg: -name")
		return
	}
	if *image == "" {
		fmt.Println("missing arg: -image")
		return
	}
	if *audio == "" {
		fmt.Println("missing arg: -audio")
		return
	}
	if *year == "" {
		fmt.Println("missing arg: -year")
		return
	}

	dir := os.Getenv("FINAL_SCENES_DIR")
	if dir == "" {
		fmt.Println("missing environment variable: FINAL_SCENES_DIR")
		return
	}
	filmDir := os.Getenv("FINAL_SCENES_TARGET")
	if filmDir == "" {
		fmt.Println("missing environment variable: FINAL_SCENES_TARGET")
		return
	}

	db, err := ioutil.ReadFile(fmt.Sprintf("%s/films.json", dir))
	if err != nil {
		fmt.Printf("failed to load db file: %v\n", err)
		return
	}

	films := struct {
		Films []models.Film `json:"Films"`
	}{}

	err = json.Unmarshal(db, &films)
	if err != nil {
		fmt.Printf("failed to unmarshal film db: %v\n", err)
		return
	}

	for _, f := range films.Films {
		if f.Name == *name {
			fmt.Printf("Film: %s is already in the database\n", *name)
			return
		}
	}

	audioSrc, err := os.Open(fmt.Sprintf("%s/%s", filmDir, *audio))
	if err != nil {
		fmt.Printf("failed to open audio file: %v\n", err)
		return
	}
	defer audioSrc.Close()
	imageSrc, err := os.Open(fmt.Sprintf("%s/%s", filmDir, *image))
	if err != nil {
		fmt.Printf("failed to open image file: %v\n", err)
		return
	}
	defer imageSrc.Close()

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(*name)))

	audioDst, err := os.Create(fmt.Sprintf("%s/www/audio/%s.mp3", dir, hash))
	if err != nil {
		fmt.Printf("failed to create new audio file: %v\n", err)
		return
	}
	defer audioDst.Close()
	imageDst, err := os.Create(fmt.Sprintf("%s/www/images/%s.jpg", dir, hash))
	if err != nil {
		fmt.Printf("failed to create new image file: %v\n", err)
		return
	}
	defer imageDst.Close()
	_, err = io.Copy(audioDst, audioSrc)
	if err != nil {
		fmt.Printf("failed to copy source audio to destination: %v\n", err)
		return
	}
	_, err = io.Copy(imageDst, imageSrc)
	if err != nil {
		fmt.Printf("failed to copy source image to destination: %v\n", err)
		return
	}
	films.Films = append(films.Films, models.Film{

		Name:      *name,
		AudioFile: fmt.Sprintf("audio/%s.mp3", hash),
		ImageFile: fmt.Sprintf("images/%s.jpg", hash),
		Year:      *year,
	})

	filmsOut, err := json.Marshal(films)
	if err != nil {
		fmt.Printf("failed to marshal films: %v\n", err)
		return
	}
	outFile, err := os.OpenFile(fmt.Sprintf("%s/films.json", dir), os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("failed to open db file: %v\n", err)
		return
	}
	defer outFile.Close()
	_, err = outFile.Write(filmsOut)
	if err != nil {
		fmt.Printf("failed to write films to db file: %v\n", err)
		return
	}
}
