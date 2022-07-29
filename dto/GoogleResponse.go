package dto

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type GoogleResultItem struct {
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
	Title   string `json:"title"`
}
type GoogleSearchResult struct {
	Items []GoogleResultItem `json:"items"`
}

func MakeRequest(google_api_key, google_search_context string) GoogleSearchResult {

	currentTime := time.Now()

	timeString := strings.ReplaceAll(currentTime.Format("01-02"), "-", "%20")

	searchQuery := fmt.Sprintf("florida%20man%20%s", timeString)

	url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s", google_api_key, google_search_context, searchQuery)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var gResp GoogleSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		log.Fatal("Error parsing the response body")
	}
	return gResp
}
