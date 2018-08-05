package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const FITBIT_API_URL = "https://api.fitbit.com/0/user/"
const FITBIT_WEIGHT_ENDPOINT = "/body/log/weight/date/"
const FITBIT_FAT_ENDPOINT = "/body/log/fat/date/"

type Outline struct {
	Body []Body `json:"weight"`
}

type Body struct {
	Weight float64 `json:"weight"`
	Bmi    float64 `json:"bmi"`
	Fat    float64 `json:"fat"`
}

func main() {
	uid := os.Getenv("FITBIT_USER")
	at := "Bearer " + os.Getenv("FITBIT_ACCESS_TOKEN")

	fmt.Println(uid)
	fmt.Println(at)

	// for date
	now := time.Now()
	layout := "2006-01-02"
	todayString := now.Format(layout)

	// Get body weight
	requestString := FITBIT_API_URL + uid + FITBIT_WEIGHT_ENDPOINT + todayString + ".json"
	req, _ := http.NewRequest("GET", requestString, nil)
	req.Header.Set("Authorization", at)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err!:", err)
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)

	// parse JSON
	outline := new(Outline)
	err = json.Unmarshal(byteArray, &outline)
	if err != nil {
		fmt.Println("unmarshal problem occured : ", err)
	}
	fmt.Println(outline.Body[0].Weight)
	fmt.Println(outline.Body[0].Bmi)
	fmt.Println(outline.Body[0].Fat)

	// ToDo : Alert there are no data if record doesn't exist
	// ToDo : Show today's weight, fat, and Bmi
	// ToDo : Save today's data for Amazon Dynamo DB
	// ToDo : Compare today and yesterday weight
	// ToDo : Remind goal and diplay calculated pace I need to keep for achieving the goal
	// ToDo : Make graph and upload it
	// ToDo : Implement using reflesh token for OAuth2

}
