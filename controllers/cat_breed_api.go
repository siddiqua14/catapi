package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/beego/beego/v2/server/web"
)

// Update the caller to pass an HTTP client
func (c *CatController) GetBreeds() {
	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	// Create a channel to receive breeds
	breedChan := make(chan []Breed)
	errorChan := make(chan error)

	go FetchBreeds(apiURL, apiKey, &http.Client{}, breedChan, errorChan) // Use FetchBreeds here

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

	// Create a channel to receive breed images
	imageChan := make(chan []CatImage)
	errorChan := make(chan error)

	// Use a real HTTP client in production
	client := &http.Client{}

	// Use the exported FetchBreedImages function
	go FetchBreedImages(apiURL, apiKey, breedID, client, imageChan, errorChan)

	select {
	case images := <-imageChan:
		c.Data["json"] = images
	case err := <-errorChan:
		c.Data["json"] = map[string]string{"error": err.Error()}
	}

	c.ServeJSON()
}

// Fetch Breeds Concurrently
func FetchBreeds(apiURL, apiKey string, client HTTPClient, breedChan chan []Breed, errorChan chan error) {
    reqURL := apiURL + "/breeds"
    req, _ := http.NewRequest("GET", reqURL, nil)
    req.Header.Add("x-api-key", apiKey)

    resp, err := client.Do(req)
    if err != nil {
        errorChan <- err
        close(breedChan)
        close(errorChan)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        errorChan <- fmt.Errorf("API returned status code %d", resp.StatusCode)
        close(breedChan)
        close(errorChan)
        return
    }

    var result []Breed
    err = json.Unmarshal(body, &result)
    if err != nil {
        errorChan <- err
        close(breedChan)
        close(errorChan)
        return
    }

    breedChan <- result
    close(breedChan)
    close(errorChan)
}

// Fetch Breed Images Concurrently

func FetchBreedImages(apiURL, apiKey, breedID string, client HTTPClient, imageChan chan []CatImage, errorChan chan error) {
	reqURL := apiURL + "/images/search?breed_ids=" + breedID + "&limit=10"
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		errorChan <- err
		close(imageChan)
		close(errorChan)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		errorChan <- fmt.Errorf("API returned status code %d", resp.StatusCode)
		close(imageChan)
		close(errorChan)
		return
	}

	var result []CatImage
	err = json.Unmarshal(body, &result)
	if err != nil {
		errorChan <- err
		close(imageChan)
		close(errorChan)
		return
	}

	imageChan <- result
	close(imageChan)
	close(errorChan)
}