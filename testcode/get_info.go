package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	//Get information for connecting Wunderlist
	accessToken := os.Getenv("WUND_ACTOKEN")
	clientID := os.Getenv("WUND_CLIENT")

	wunderlistURL := "https://a.wunderlist.com/api/v1/"
	endpoint := "lists"

	//create request
	req, _ := http.NewRequest("GET", wunderlistURL+endpoint, nil)
	req.Header.Set("X-Access-Token", accessToken)
	req.Header.Set("X-Client-ID", clientID)

	//send http GET request to Wunderlist
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Access Key Problem?")
	}

	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	//parse JSON by official way
	type List struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		CreateAt string `json:"create_at"`
		ListType string `json:"list_type"`
		Revision int    `json:"revision"`
	}
	var lists []List
	err = json.Unmarshal(byteArray, &lists)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf(lists[1].Title)

}
