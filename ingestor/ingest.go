package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

type station struct {
	StationID       string `json:"station_id"`
	BikesAvailable  int    `json:"num_bikes_available"`
	EBikesAvailable int    `json:"num_ebikes_available"`
	BikesDisabled   int    `json:"num_bikes_disabled"`
	DocksAvailable  int    `json:"num_docks_available"`
	DocksDisabled   int    `json:"num_docks_disabled"`
	LastReported    int    `json:"last_reported"`
}

type stationData struct {
	Stations []station `json:"stations"`
}

// structures defined based on JSON response format from stationDataURL
type stationsDataResponse struct {
	LastUpdated int         `json:"last_updated"`
	TTL         int         `json:"ttl"`
	Data        stationData `json:"data"`
}

func getStationsDataResponse(body []byte) (*stationsDataResponse, error) {
	var response = new(stationsDataResponse)
	err := json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error while decoding stationDataResponse: %s\n", err)
	}

	return response, err
}

// private constants are defined in lowerCamelCase, use UpperCamelCase or ALL_CAPS to export the variables
const stationDataURL string = "https://gbfs.fordgobike.com/gbfs/en/station_status.json"

func main() {
	fmt.Println("Starting data ingestion module ...")

	for true {
		fmt.Println("Sending GET request")
		request := gorequest.New()
		_, body, errs := request.Get(stationDataURL).End()

		if errs != nil {
			fmt.Printf("Error while getting data from URL: %s\n", errs)
			fmt.Println("Will try again in few moments")
		} else {
			stationDataResponse, err := getStationsDataResponse([]byte(body))
			if err != nil {
				fmt.Printf("Error while converting json string into struct: %s\n", err)
				panic(err)
			}

			fmt.Printf("Got %d stations data\n", len(stationDataResponse.Data.Stations))

		}

		time.Sleep(10 * time.Second)
	}
}
