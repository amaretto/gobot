package wunderlist

import (
    "fmt"    
)

func getListInfo() {

}

func getTaskInfo() {

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

	//TODO implements
	fmt.Println(byteArray)
}

func sendTaskInfo() {

}
