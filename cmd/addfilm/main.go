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
	flag.Parse()

	if *url == "" {
		fmt.Println("missing arg: -url")
		return
	}

	err := DownloadVideo(*url, "video.mp4")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	err = TrimVideo("video.mp4", "trimmed.mp4", "00:01:00", "00:01:55")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	err = ExtractAudio("trimmed.mp4", "audio.mp3")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	err = ExtractFrames("video.mp4", "frame.jpg", "00:01:52", "1")
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
