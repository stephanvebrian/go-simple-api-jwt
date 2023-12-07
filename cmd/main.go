package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
)

func main() {
	viper.SetConfigFile("files/config.yaml")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err)
	}

	// Unmarshal the configuration into a map
	var configMap map[string]interface{}
	if err := viper.Unmarshal(&configMap); err != nil {
		log.Fatal().Err(err)
	}

	// Use mapstructure to decode the map into the Config struct
	var cfg config.Config
	if err := mapstructure.Decode(configMap, &cfg); err != nil {
		log.Fatal().Err(err)
	}

	go startScheduler(cfg)

	// DoRequest()
	log.Info().Msg("running application...")

	err := startServer(cfg)
	if err != nil {
		log.Fatal().Err(err)
	}
}

func DoRequest() {
	apiURL := "http://13.212.226.116:8000/api/register"

	// Create a new buffer to store the request body
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add fields to the form data
	fields := map[string]string{
		"username":   "admin123@gmail.com",
		"password":   "admin",
		"first_name": "Admin",
		"last_name":  "Last-Admin",
		"telephone":  "+628976833573",
		// "profile_image": "https://images.unsplash.com/photo-1648075082539-ca4a311d2afa?q=80&w=1170",
		"address":  "Jln Palmerah 53",
		"city":     "Jakarta Barata",
		"province": "DKI Jakarta",
		"country":  "Indonesia",
	}

	for key, value := range fields {
		writer.WriteField(key, value)
	}

	file := GetImageFileFromUrl("https://images.unsplash.com/photo-1648075082539-ca4a311d2afa?q=80&w=1170")

	// Add the file to the form data
	part, err := writer.CreateFormFile("profile_image", "image.jpg")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file data:", err)
		return
	}

	// Close the multipart writer to finalize the request body
	writer.Close()

	// Create a new POST request with the form data
	req, err := http.NewRequest("POST", apiURL, &requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the Content-Type header to indicate multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// // Read and print the response body
	// var result interface{}
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	fmt.Println("Error decoding JSON response:", err)
	// 	return
	// }
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(body))

	// // Check the response status
	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Printf("Error: %s\n", result)
	// 	return
	// }

	// fmt.Printf("API response: %+v\n", result)
}

// Download and get image from imageUrl
func GetImageFileFromUrl(imageUrl string) io.Reader {
	response, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return nil
	}
	defer response.Body.Close()

	// Read the entire response body into a buffer
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	return bytes.NewReader(body)
}

// Open the file for reading e.g "/path/to/your/image.jpg"
func GetImageFileFromPath(path string) io.Reader {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	return file
}
