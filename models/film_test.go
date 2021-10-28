package models

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct{
		name string
		input string
		expected Film
		err bool
	}{
		{
			name: "Valid film json format",
			input: "{\"Name\": \"TestFilm\", \"AudioFile\": \"TestAudioFile\", \"Year\": \"2000\", \"ImageFile\": \"TestImageFile\"}",
			expected: Film{
				Name: "TestFilm",
				AudioFile: "TestAudioFile",
				Year: "2000",
				ImageFile: "TestImageFile",
				Hash: fmt.Sprintf("%x", sha256.Sum256([]byte("testfilm"))),
			},
			err: false,
		},
		{
			name: "Missing value",
			input: "{\"Name\":, \"AudioFile\": \"TestAudioFile\", \"Year\": \"2000\", \"ImageFile\": \"TestImageFile\"}",
			expected: Film{
				Name: "TestFilm",
				AudioFile: "TestAudioFile",
				Year: "2000",
				ImageFile: "TestImageFile",
				Hash: fmt.Sprintf("%x", sha256.Sum256([]byte("testfilm"))),
			},
			err: true,
		},
		{
			name: "Missing field",
			input: "{\"AudioFile\": \"TestAudioFile\", \"Year\": \"2000\", \"ImageFile\": \"TestImageFile\"}",
			expected: Film{
				Name: "TestFilm",
				AudioFile: "TestAudioFile",
				Year: "2000",
				ImageFile: "TestImageFile",
				Hash: fmt.Sprintf("%x", sha256.Sum256([]byte("testfilm"))),
			},
			err: true,
		},
		{
			name: "Blank field value",
			input: "{\"Name\": \"\", \"AudioFile\": \"TestAudioFile\", \"Year\": \"2000\", \"ImageFile\": \"TestImageFile\"}",
			expected: Film{
				Name: "",
				AudioFile: "TestAudioFile",
				Year: "2000",
				ImageFile: "TestImageFile",
				Hash: fmt.Sprintf("%x", sha256.Sum256([]byte(""))),
			},
			err: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Film{}
			err := got.UnmarshalJSON([]byte(tc.input))
			if err != nil {
				if tc.err {
					return
				}
				t.Fatalf("received error when not expected: %v", err)
			}
			if tc.err {
				t.Fatalf("expected an error but didn't get one")
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected:\n%+v\nGot:\n%+v\n", tc.expected, got)
			}
		})
	}
}
