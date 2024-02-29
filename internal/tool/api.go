package tool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Tag struct {
	Name string `json:"name"`
}

func GetLatestTag() (string, error) {
	url := "https://api.github.com/repos/zaigie/palworld-server-tool/tags"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var tags []Tag
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return "", err
	}
	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}
	return tags[0].Name, nil
}
