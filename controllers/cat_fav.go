package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	
)

func (c *CatController) CreateFavorite() {
    apiKey, _ := web.AppConfig.String("catapi.key")
    apiURL, _ := web.AppConfig.String("catapi.url")

    var favorite struct {
        ImageID string `json:"image_id"`
    }
    if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&favorite); err != nil {
        c.Data["json"] = map[string]string{"error": "Invalid request body: " + err.Error()}
        c.Ctx.ResponseWriter.WriteHeader(400)
        c.ServeJSON()
        return
    }
	if favorite.ImageID == "" {
        c.Data["json"] = map[string]string{"error": "image_id is required"}
        c.Ctx.ResponseWriter.WriteHeader(400)
        c.ServeJSON()
        return
    }
    reqURL := fmt.Sprintf("%s/favourites", apiURL)
    jsonData, _ := json.Marshal(favorite)
    req, _ := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
    req.Header.Set("x-api-key", apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := c.httpClient
	if client == nil {
        client = &http.Client{}
    }
    resp, err := client.Do(req)
    if err != nil {
        errMsg := fmt.Sprintf("Failed to create favorite: %v", err)
        c.Data["json"] = map[string]string{"error": errMsg}
        c.Ctx.ResponseWriter.WriteHeader(500)
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    if resp.StatusCode != 200 && resp.StatusCode != 201 {
        errMsg := fmt.Sprintf("API returned status code %d: %s", resp.StatusCode, string(body))
        c.Data["json"] = map[string]string{"error": errMsg}
        c.Ctx.ResponseWriter.WriteHeader(resp.StatusCode)
        c.ServeJSON()
        return
    }

    c.Data["json"] = map[string]string{"status": "success"}
    c.ServeJSON()
}
// Handle fetching favorite cat images
func (c *CatController) GetFavorites() {
    apiKey, _ := web.AppConfig.String("catapi.key")
    apiURL, _ := web.AppConfig.String("catapi.url")

    reqURL := fmt.Sprintf("%s/favourites", apiURL)
    req, _ := http.NewRequest("GET", reqURL, nil)
    req.Header.Set("x-api-key", apiKey)

    client := c.httpClient
    if client == nil {
        client = &http.Client{}
    }
    resp, err := client.Do(req)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to fetch favorites"}
        c.Ctx.ResponseWriter.WriteHeader(500)
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        c.Data["json"] = map[string]string{"error": "Failed to fetch favorites"}
        c.Ctx.ResponseWriter.WriteHeader(resp.StatusCode)
        c.ServeJSON()
        return
    }

    var result interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to parse favorites"}
        c.Ctx.ResponseWriter.WriteHeader(500)
        c.ServeJSON()
        return
    }

    c.Data["json"] = result
    c.ServeJSON()
}
// Handle deleting a favorite cat image
func (c *CatController) DeleteFavorite() {
    // Get favorite ID from URL parameter
    favoriteId := c.Ctx.Input.Param(":id")
    
    // Get API configuration
    apiKey, _ := web.AppConfig.String("catapi.key")
    apiURL, _ := web.AppConfig.String("catapi.url")
    
    // Construct delete request
    reqURL := fmt.Sprintf("%s/favourites/%s", apiURL, favoriteId)
    req, err := http.NewRequest("DELETE", reqURL, nil)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to create delete request"}
        c.Ctx.ResponseWriter.WriteHeader(500)
        c.ServeJSON()
        return
    }
    
    // Set API key header
    req.Header.Set("x-api-key", apiKey)
    
    // Send request
	client := c.httpClient
    if client == nil {
        client = &http.Client{}
    }
    resp, err := client.Do(req)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to delete favorite"}
        c.Ctx.ResponseWriter.WriteHeader(500)
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()
    
    // Check response
    if resp.StatusCode != 200 {
        body, _ := io.ReadAll(resp.Body)
        var result map[string]interface{}
        json.Unmarshal(body, &result)
        c.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to delete favorite: %v", result)}
        c.Ctx.ResponseWriter.WriteHeader(resp.StatusCode)
        c.ServeJSON()
        return
    }
    
    // Return success response
    c.Data["json"] = map[string]string{"message": "Favorite deleted successfully"}
    c.ServeJSON()
}

