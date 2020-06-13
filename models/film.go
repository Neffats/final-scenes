package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
)

type Film struct {
	Name      string `json:"Name"`
	AudioFile string `json:"AudioFile"`
	Year      string `json:"Year"`
	ImageFile string `json:"ImageFile"`
	Hash      string `json:"Hash,omitempty"`
}

func (fs *Film) UnmarshalJSON(data []byte) error {
	var film map[string]string
	err := json.Unmarshal(data, &film)
	if err != nil {
		return err
	}

	name, ok := film["Name"]
	if !ok {
		return fmt.Errorf("failed to unmarshal data: missing field 'Name'")
	}

	audioFile, ok := film["AudioFile"]
	if !ok {
		return fmt.Errorf("failed to unmarshal data: missing field 'AudioFile'")
	}
	year, ok := film["Year"]
	if !ok {
		return fmt.Errorf("failed to unmarshal data: missing field 'Year'")
	}
	imageFile, ok := film["ImageFile"]
	if !ok {
		return fmt.Errorf("failed to unmarshal data: missing field 'ImageFile'")
	}
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.ToLower(name))))

	fs.Name = name
	fs.AudioFile = audioFile
	fs.Year = year
	fs.ImageFile = imageFile
	fs.Hash = hash

	return nil
}
