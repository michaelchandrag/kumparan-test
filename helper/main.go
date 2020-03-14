package helper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"bytes"
	"encoding/json"
	"strconv"
	"os"
	"time"

	model "bitbucket.org/michaelchandrag/kumparan-test/model"
)

var ES_HOST = os.Getenv("ES_HOST")

type (
	EsResult struct {
		Took 			int 			`json:"took"`
		HitsResult 		HitsResponse 	`json:"hits"`
	}

	HitsResponse struct {
		HitsHits		[]Detail 		`json:"hits"`
		HitsTotal 		Total 			`json:"total"`
		HitsMaxScore 	float64 		`json:"max_score"`
	}

	Total struct {
		Value 		int 	`json:"value"`
		Relation 	string 	`json:"relation"`
	}

	Detail struct {
		Index 		string 			`json:"_index"`
		Type 		string 			`json:"_type"`
		ID 			string 			`json:"_id"`
		Source 		Data 			`json:"_source"`
		News 		model.News 		`json:"news",omitempty`
	}

	Data struct {
		ID 			int 		`json:"id"`
		Created 	string 		`json:"created"`
	}

	EsRequest struct {
		Sort 	[]SortBody 		`json:"sort"`
		From 	int 			`json:"from"`
		Size 	int 			`json:"size"`
	}

	SortBody struct {
		Created 	SortField 		`json:"created"`
	}

	SortField struct {
		Order 		string 		`json:"order"`
	}

)

func EsPost(news model.News) error {
	type customModel struct {
		ID 		int `json:"id"`
		Created string `json:"created"`
	}
	reqBody := customModel{
		ID: news.ID,
		Created: news.EsCreated,
	}

	reqJson, err := json.Marshal(&reqBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/kumparan/news", ES_HOST), "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func EsGet(query url.Values) (result EsResult) {
	offset := 0
	limit := 10
	if _, ok := query["page"]; ok {
		page, _ := strconv.Atoi(query["page"][0])
		if (page < 1) {
			page = 1
		}
		offset = (page - 1) * limit
	}

	sortFieldCreated := SortField{
		Order: "desc",
	}

	sortBody := []SortBody {
		SortBody {
			Created: sortFieldCreated,
		},
	}

	reqBody := EsRequest{
		Sort: sortBody,
		From: offset,
		Size: limit,
	}

	reqJson, err := json.Marshal(&reqBody)	
	
	client := http.Client{
		Timeout: time.Duration(5*time.Second),
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/kumparan/_search", ES_HOST), bytes.NewBuffer(reqJson))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
		return result
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return result
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&result)

	return result
}