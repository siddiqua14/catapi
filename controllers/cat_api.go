package controllers

import (
    "encoding/json"
	"log"
	"fmt"
	"io"
	"net/http"
	"github.com/beego/beego/v2/server/web"
)


// ServeSingleCatImage fetches and serves a single cat image
func (c *CatController) GetCatImage() {
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{}) // Use the correct map type
	}

	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	imageURL, err := FetchCatImage(apiURL, apiKey, c.getHTTPClient())
	if err != nil {
		c.Data["CatImage"] = "" // Use string key
	} else {
		c.Data["CatImage"] = imageURL
	}
	c.TplName = "index.tpl"
}

func (c *CatController) GetCatImagesAPI() {
	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	imageChan := make(chan []CatImage)
	errorChan := make(chan error)

	// Pass the correct arguments to FetchMultipleImageURLs
	go FetchCatImages(apiURL, apiKey, &http.Client{}, imageChan, errorChan)

	select {
	case images := <-imageChan:
		c.Data["json"] = images
	case err := <-errorChan:
		c.Data["json"] = map[string]string{"error": err.Error()}
	}

	c.ServeJSON()
}
func FetchCatImage(apiURL, apiKey string, client HTTPClient) (string, error) {
    reqURL := apiURL + "/images/search?limit=1"
    req, err := http.NewRequest("GET", reqURL, nil)
    if err != nil {
        log.Printf("Error creating request: %v", err)
        return "", err
    }

    req.Header.Add("x-api-key", apiKey)

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making request: %v", err)
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response: %v", err)
        return "", err
    }

    if resp.StatusCode != http.StatusOK {
        log.Printf("API returned non-200 status code: %d", resp.StatusCode)
        return "", fmt.Errorf("API returned status code %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
    }

    var result []CatImage
    if err := json.Unmarshal(body, &result); err != nil {
        log.Printf("Error parsing JSON response: %v", err)
        return "", err
    }

    if len(result) == 0 {
        log.Print("No images returned from API")
        return "", fmt.Errorf("no images returned from API")
    }

    return result[0].URL, nil
}

func FetchCatImages(apiURL, apiKey string, client HTTPClient, imageChan chan []CatImage, errorChan chan error) {
    defer close(imageChan)
    defer close(errorChan)

    reqURL := apiURL + "/images/search?limit=10"
    req, err := http.NewRequest("GET", reqURL, nil)
    if err != nil {
        log.Printf("Error creating request: %v", err)
        errorChan <- err
        return
    }

    req.Header.Add("x-api-key", apiKey)

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making request: %v", err)
        errorChan <- err
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response: %v", err)
        errorChan <- err
        return
    }

    if resp.StatusCode != http.StatusOK {
        log.Printf("API returned non-200 status code: %d", resp.StatusCode)
        errorChan <- fmt.Errorf("API returned status code %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
        return
    }

    var result []CatImage
    if err := json.Unmarshal(body, &result); err != nil {
        log.Printf("Error parsing JSON response: %v", err)
        errorChan <- err
        return
    }

    imageChan <- result
}