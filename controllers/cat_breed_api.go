package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/beego/beego/v2/server/web"
)

// GetBreeds retrieves all cat breeds
func (c *CatController) GetBreeds() {
	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	breedChan := make(chan []Breed)
	errorChan := make(chan error)

	go FetchBreeds(apiURL, apiKey, &http.Client{}, breedChan, errorChan)

	select {
	case breeds := <-breedChan:
		c.Data["json"] = breeds
	case err := <-errorChan:
		c.Data["json"] = map[string]string{"error": err.Error()}
	}

	c.ServeJSON()
}

// GetBreedImages retrieves images for a specific breed
func (c *CatController) GetBreedImages() {
	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	breedID := c.GetString("breed_id")

	imageChan := make(chan []CatImage)
	errorChan := make(chan error)

	go FetchBreedImages(apiURL, apiKey, breedID, &http.Client{}, imageChan, errorChan)

	select {
	case images := <-imageChan:
		c.Data["json"] = images
	case err := <-errorChan:
		c.Data["json"] = map[string]string{"error": err.Error()}
	}

	c.ServeJSON()
}

// FetchBreeds retrieves cat breeds from the API
func FetchBreeds(apiURL, apiKey string, client HTTPClient, breedChan chan []Breed, errorChan chan error) {
	defer close(breedChan)
	defer close(errorChan)

	reqURL := apiURL + "/breeds"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("error creating request at %s:%d: %v", file, line, err)
		return
	}
	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("error making request at %s:%d: %v", file, line, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("error reading response at %s:%d: %v", file, line, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("API returned status code %d at %s:%d: %s", resp.StatusCode, file, line, http.StatusText(resp.StatusCode))
		return
	}

	var result []Breed
	if err := json.Unmarshal(body, &result); err != nil {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("error parsing response at %s:%d: %v", file, line, err)
		return
	}

	breedChan <- result
}

// FetchBreedImages retrieves cat images for a specific breed
func FetchBreedImages(apiURL, apiKey, breedID string, client HTTPClient, imageChan chan []CatImage, errorChan chan error) {
	defer close(imageChan)
	defer close(errorChan)

	reqURL := fmt.Sprintf("%s/images/search?breed_ids=%s&limit=10", apiURL, breedID)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("error creating request at %s:%d: %v", file, line, err)
		return
	}
	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("error making request at %s:%d: %v", file, line, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("error reading response at %s:%d: %v", file, line, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("API returned status code %d at %s:%d: %s", resp.StatusCode, file, line, http.StatusText(resp.StatusCode))
		return
	}

	var result []CatImage
	if err := json.Unmarshal(body, &result); err != nil {
		file, line := getCallerInfo()
		errorChan <- fmt.Errorf("error parsing response at %s:%d: %v", file, line, err)
		return
	}

	imageChan <- result
}