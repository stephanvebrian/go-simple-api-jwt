package helper

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

// UT injection purpose
var (
	fnOSOpen = os.Open
	fnIOCopy = io.Copy
)

// Download and get image from imageUrl
func GetImageFileFromUrl(imageUrl string) io.Reader {
	response, err := http.Get(imageUrl)
	if err != nil {
		log.Error().Err(err).Str("message", "error downloading image")
		return nil
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Str("message", "error reading response body")
		return nil
	}

	return bytes.NewReader(body)
}

// Open the file for reading e.g "/path/to/your/image.jpg"
func GetImageFileFromPath(path string) io.Reader {
	file, err := fnOSOpen(path)
	if err != nil {
		log.Error().Err(err).Str("message", "error opening file")
		return nil
	}
	defer file.Close()

	// Read the content of the file into a buffer
	var buffer bytes.Buffer
	_, err = fnIOCopy(&buffer, file)
	if err != nil {
		fmt.Println("error io")
		log.Error().Err(err).Str("message", "error reading file content")
		return nil
	}

	// Return a reader for the buffer
	return bytes.NewReader(buffer.Bytes())
}
