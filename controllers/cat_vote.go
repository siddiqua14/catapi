package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	
)
// Handle voting on a cat image
func (c *CatController) CreateVote() {
	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	var vote Vote
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&vote); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid request body: " + err.Error()}
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}

	if vote.Value != 1 && vote.Value != -1 {
		c.Data["json"] = map[string]string{"error": "Vote value must be 1 or -1"}
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}

	// Log the vote being created
	fmt.Printf("Creating vote for image %s with value %d\n", vote.ImageID, vote.Value)

	reqURL := fmt.Sprintf("%s/votes", apiURL)
	jsonData, _ := json.Marshal(vote)
	req, _ := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create vote: %v", err)
		fmt.Println(errMsg)
		c.Data["json"] = map[string]string{"error": errMsg}
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("API Response: %s\n", string(body))

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		errMsg := fmt.Sprintf("API returned status code %d: %s", resp.StatusCode, string(body))
		fmt.Println(errMsg)
		c.Data["json"] = map[string]string{"error": errMsg}
		c.Ctx.ResponseWriter.WriteHeader(resp.StatusCode)
		c.ServeJSON()
		return
	}

	fmt.Println("Vote created successfully")
	c.Data["json"] = map[string]string{"status": "success"}
	c.ServeJSON()
}

func (c *CatController) GetVotes() {
	apiKey, _ := web.AppConfig.String("catapi.key")
	apiURL, _ := web.AppConfig.String("catapi.url")

	reqURL := fmt.Sprintf("%s/votes", apiURL)
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("x-api-key", apiKey)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch votes"}
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		c.Data["json"] = map[string]string{"error": "Failed to fetch votes"}
		c.Ctx.ResponseWriter.WriteHeader(resp.StatusCode)
		c.ServeJSON()
		return
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse votes"}
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.ServeJSON()
		return
	}

	c.Data["json"] = result
	c.ServeJSON()
}
