package controllers

import (
	"net/http"
	"github.com/beego/beego/v2/server/web"
	
)
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type CatController struct {
	web.Controller
	httpClient HTTPClient
}
// SetHTTPClient allows injection of mock client for testing
func (c *CatController) SetHTTPClient(client HTTPClient) {
    c.httpClient = client
}

// getHTTPClient returns the http client to use
func (c *CatController) getHTTPClient() HTTPClient {
    if c.httpClient != nil {
        return c.httpClient
    }
    return &http.Client{} // default client
}

type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}



type CatImage struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	MimeType  string   `json:"mime_type"`
	Breeds    []Breed  `json:"breeds"`
	Categories []string `json:"categories"`
}
type Vote struct {
	ImageID string `json:"image_id"`
	Value   int    `json:"value"`
}

type Breed struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Origin       string `json:"origin"`
	WikipediaURL string `json:"wikipedia_url"`
}



/*
// GetCatImage handles the web request for a cat image
func (c *CatController) GetCatImage() {
    if c.Data == nil {
        c.Data = make(map[interface{}]interface{}) // Use the correct map type
    }

    apiKey, _ := web.AppConfig.String("catapi.key")
    apiURL, _ := web.AppConfig.String("catapi.url")

    imageURL, err := c.FetchCatImage(apiURL, apiKey)
    if err != nil {
        c.Data["CatImage"] = "" // Use string key
    } else {
        c.Data["CatImage"] = imageURL
    }
    c.TplName = "index.tpl"
}

// Modify GetCatImagesAPI to use the new signature of FetchCatImages
func (c *CatController) GetCatImagesAPI() {
	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	client := &http.Client{} // Create an HTTP client

	// Create a channel to receive cat images
	imageChan := make(chan []CatImage)
	errorChan := make(chan error)

	// Use the client in the FetchCatImages call
	go FetchCatImages(client, apiURL, apiKey, imageChan, errorChan)

	select {
	case images := <-imageChan:
		c.Data["json"] = images
	case err := <-errorChan:
		c.Data["json"] = map[string]string{"error": err.Error()}
	}

	c.ServeJSON()
}
// FetchCatImage fetches a cat image from the API
func (c *CatController) FetchCatImage(apiURL, apiKey string) (string, error) {
	reqURL := apiURL + "/images/search?limit=1"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("x-api-key", apiKey)

	client := c.getHTTPClient()
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

// FetchCatImages fetches multiple cat images from the API
func FetchCatImages(client HTTPClient, apiURL, apiKey string, imageChan chan []CatImage, errorChan chan error) {
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
*/


