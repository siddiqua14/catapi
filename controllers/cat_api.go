package controllers

import (
    "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/beego/beego/v2/server/web"
)


// ServeSingleCatImage fetches and serves a single cat image
func (c *CatController) ServeSingleCatImage() {
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{}) // Use the correct map type
	}

	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	imageURL, err := FetchSingleCatImage(apiURL, apiKey, c.getHTTPClient())
	if err != nil {
		c.Data["CatImage"] = "" // Use string key
	} else {
		c.Data["CatImage"] = imageURL
	}
	c.TplName = "index.tpl"
}

func (c *CatController) ServeMultipleCatImages() {
	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	imageChan := make(chan []CatImage)
	errorChan := make(chan error)

	// Pass the correct arguments to FetchMultipleImageURLs
	go FetchMultipleImageURLs(apiURL, apiKey, &http.Client{}, imageChan, errorChan)

	select {
	case images := <-imageChan:
		c.Data["json"] = images
	case err := <-errorChan:
		c.Data["json"] = map[string]string{"error": err.Error()}
	}

	c.ServeJSON()
}

// FetchSingleImageURL fetches the URL of a single cat image from the API
func FetchSingleCatImage(apiURL, apiKey string, client HTTPClient) (string, error) {
	reqURL := apiURL + "/images/search?limit=1"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var result []CatImage
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("no images returned from API")
	}

	return result[0].URL, nil
}

// FetchMultipleImageURLs fetches the URLs of multiple cat images from the API
func FetchMultipleImageURLs(apiURL, apiKey string, client HTTPClient, imageChan chan []CatImage, errorChan chan error) {
	defer close(imageChan)
	defer close(errorChan)

	reqURL := apiURL + "/images/search?limit=10"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		errorChan <- fmt.Errorf("error creating request: %v", err)
		return
	}

	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		errorChan <- fmt.Errorf("error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorChan <- fmt.Errorf("error reading response: %v", err)
		return
	}

	if resp.StatusCode != 200 {
		errorChan <- fmt.Errorf("API returned status code %d", resp.StatusCode)
		return
	}

	var result []CatImage
	err = json.Unmarshal(body, &result)
	if err != nil {
		errorChan <- fmt.Errorf("error parsing response: %v", err)
		return
	}

	imageChan <- result
}
