package main

import (
	/*
	"crypto/sha256"
	"encoding/json"
	
	"io"
	"io/ioutil"*/
	"os"

	"flag"
	"fmt"
	"os/exec"
	"strings"

	//"github.com/Neffats/final-scenes/models"
)

func main() {
	url := flag.String("url", "", "Youtube url")
	name := flag.String("name", "", "Film name")
	startTime := flag.String("start", "", "Clip start timestamp. Format: hh:mm:ss")
	endTime := flag.String("end", "", "Clip end timestamp. Format: hh:mm:ss")
	year := flag.String("year", "", "Film year of release")
	imageTime := flag.String("image", "", "Timestamp for image. Format: hh:mm:ss")
	flag.Parse()

	if *url == "" {
		fmt.Println("missing arg: -url")
		return
	}
	if *name == "" {
		fmt.Println("missing arg: -name")
		return
	}
	if *startTime == "" {
		fmt.Println("missing arg: -start")
		return
	}
	if *endTime == "" {
		fmt.Println("missing arg: -end")
		return
	}
	if *year == "" {
		fmt.Println("missing arg: -year")
		return
	}
	if *imageTime == "" {
		fmt.Println("missing arg: -year")
		return
	}

	workingDir := os.Getenv("FINAL_SCENES_ADD_FILM_DIR")
	if workingDir == "" {
		fmt.Printf("missing environment variable: FINAL_SCENES_ADD_FILM_DIR")
		return
	}

	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		fmt.Printf("directory does not exist: %s", workingDir)
		return
	}

	videoNameFull := fmt.Sprintf("%s/%s-full.mp4", workingDir, *name)
	videoNameTrim := fmt.Sprintf("%s/%s-trim.mp4", workingDir, *name)
	audioName := fmt.Sprintf("%s/%s.mp3", workingDir, *name)
	imageName := fmt.Sprintf("%s/%s.jpg", workingDir, *name)

	
	err := DownloadVideo(*url, videoNameFull)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	err = TrimVideo(videoNameFull, videoNameTrim, *startTime, *endTime)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	err = ExtractAudio(videoNameTrim, audioName)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	err = ExtractFrames(videoNameFull, imageName, *imageTime, "1")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	return
}

func DownloadVideo(url string, outputFilename string) error {
	err := runCommand("youtube-dl", "-o", outputFilename, url)
	if err != nil {
		return fmt.Errorf("error downloading video: %v", err)
	}
	return nil
}

func TrimVideo(inputFilename string, outputFilename string, startTime string, endTime string) error {
	err := runCommand("ffmpeg", "-ss", startTime, "-i", inputFilename, "-to", endTime, "-c", "copy", outputFilename)
	if err != nil {
		return fmt.Errorf("error trimming video %s: %v", inputFilename, err)
	}
	return nil
}

func ExtractAudio(inputFilename string, outputFilename string) error {
	err := runCommand("ffmpeg", "-i", inputFilename, "-f", "mp3", "-ab", "192000", "-vn", outputFilename)
	if err != nil {
		return fmt.Errorf("error extracting audio from %s: %v", inputFilename, err)
	}
	return nil
}

func ExtractFrames(inputFilename string, outputFilename string, timestamp string, numFrames string) error {
	err := runCommand("ffmpeg", "-ss", timestamp, "-i", inputFilename, "-vframes", numFrames, "-q:v", "2", outputFilename)
	if err != nil {
		return fmt.Errorf("error extracting frames from %s: %v", inputFilename, err)
	}
	return nil
}

func MoveFile(source string, destination string) error {
	err := runCommand("mv", source, destination)
	if err != nil {
		return fmt.Errorf("error moving file: %v", err)
	}

	return nil
}

func DeleteFile(filename string) error {
	if !strings.HasPrefix(filename, ".") {
		return fmt.Errorf("error deleting file: only relative paths are allowed")
	}
	err := runCommand("rm", filename)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	var errbuf strings.Builder
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = &errbuf

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("error running command: %s", errbuf.String())
	}

	return nil
	
}
