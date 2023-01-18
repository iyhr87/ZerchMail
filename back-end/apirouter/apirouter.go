package apirouter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	zincSearchURL      = os.Getenv("ZINCSEARCH_URL")
	zincSearchPort     = os.Getenv("ZINCSEARCH_PORT")
	dataBaseURL        = os.Getenv("DATABASE_URL")
	zincSearchUsarName = os.Getenv("ZINCSEARCH_ADMIN")
	zincSearchPassword = os.Getenv("ZINCSEARCH_PASS")

	urlSearch = fmt.Sprintf("%s:%s/%s/_search", zincSearchURL, zincSearchPort, dataBaseURL)
)

// ShardsStruct An structure to save the Shard attribute returned by ZincSearch.
type ShardsStruct struct {
	Total      int
	Successful int
	Skipped    int
	Failed     int
}

// TotalStruct An structure to save the Total attribute returned by ZincSearch.
type TotalStruct struct {
	Value int
}

// SourceData An structure to save the Source attribute returned by ZincSearch.
type SourceData struct {
	Date    string `json:"Date"`
	From    string `json:"From"`
	Message string `json:"BodyMessage"`
	Subject string `json:"Subject"`
	To      string `json:"To"`
}

// TargetData An structure to save additional data, related to SourceData, returned by ZincSearch.
type TargetData struct {
	Index     string     `json:"_index"`
	Type      string     `json:"_type"`
	Id        string     `json:"_id"`
	Score     float64    `json:"_score"`
	Timestamp string     `json:"timestamp"`
	Source    SourceData `json:"_source"`
}

// HitsStruct An structure to save Hits attribute returned by ZincSearch.
type HitsStruct struct {
	Total     TotalStruct  `json:"total"`
	Max_score float64      `json:"max_score"`
	Hits      []TargetData `json:"hits"`
}

// ZincSearchData Main structure to save the data returned by ZincSearch.
type ZincSearchData struct {
	Took      int          `json:"took"`
	Timed_out bool         `json:"timed_out"`
	Shards    ShardsStruct `json:"_shards"`
	Hits      HitsStruct   `json:"hits"`
}

// FrontEndData Structure to send information into the front-end.
type FrontEndData struct {
	//Id          int
	From        string
	To          string
	Subject     string
	Date        string
	BodyMessage string
}

// CreateZincSearchQuery Returns a string with the query to the database (ZincSearch) for searching a word.
func CreateZincSearchQuery(word string) string {
	query := `{
        "search_type": "match",
        "query":
        {
            "term": "%s",
            "start_time": "2022-06-02T14:28:31.894Z",
            "end_time": "2024-12-02T15:28:31.894Z"
        },
        "from": 0,
        "max_results": 10,
        "_source": ["From", "To", "Subject", "Date","BodyMessage"]
    }`

	return fmt.Sprintf(query, word)
}

// FetchZincSearchApi Fecths ZincSarch API for searching a word.
func FetchZincSearchApi(queryString string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, urlSearch, strings.NewReader(queryString))
	if err != nil {
		log.Print("Error: Building New Request: ", err)
		return nil, err
	}

	req.SetBasicAuth(zincSearchUsarName, zincSearchPassword)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print("Error: Request failed: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	log.Println("Response status code: ", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Print("Error: Reading Body: ", err)
		return nil, err
	}

	return body, nil
}

// GetRequestedData Takes the essential information, requested from the front-end, from raw data returned by ZincSearch.
func GetRequestedData(RawBodyResp []byte) ([]byte, error) {
	var ZSdata ZincSearchData

	err := json.Unmarshal([]byte(RawBodyResp), &ZSdata)
	if err != nil {
		log.Print("error: unmarshal in RawBodyResp.", err)
	}

	data := make([]FrontEndData, len(ZSdata.Hits.Hits))
	for i, h := range ZSdata.Hits.Hits {
		data[i].From = h.Source.From
		data[i].Date = h.Source.Date
		data[i].Subject = h.Source.Subject
		data[i].To = h.Source.To
		data[i].BodyMessage = h.Source.Message
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Print("error: json.Marshal(data) in RawBodyResp.", err)
	}

	return b, err
}

// SearchInDataBase Main function for requesting a search to ZincSearch.
func SearchInDataBase(word string) []byte {
	queryString := CreateZincSearchQuery(word)
	bodyRaw, err := FetchZincSearchApi(queryString)
	if err != nil {
		log.Print("error: fetchZincSearchApi(queryString) in SearchInDataBase.", err)
	}

	respFrontEnd, err := GetRequestedData(bodyRaw)
	if err != nil {
		log.Print("error: GetRequestedData(bodyRaw) in SearchInDataBase.", err)
	}

	return respFrontEnd
}
