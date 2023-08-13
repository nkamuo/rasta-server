package geo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DirectionsResponse struct {
	Routes []struct {
		OverviewPolyline struct {
			Points string `json:"points"`
		} `json:"overview_polyline"`
		Legs []struct {
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
		} `json:"legs"`
	} `json:"routes"`
}

func main() {
	origin := "37.7749,-122.4194"      // Latitude and longitude of origin
	destination := "34.0522,-118.2437" // Latitude and longitude of destination

	apiURL := "https://maps.googleapis.com/maps/api/directions/json"
	params := fmt.Sprintf("origin=%s&destination=%s&key=%s", origin, destination, googleMapsAPIKey)
	fullURL := apiURL + "?" + params

	response, err := http.Get(fullURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var data DirectionsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(data.Routes) > 0 {
		polyline := data.Routes[0].OverviewPolyline.Points
		distance := data.Routes[0].Legs[0].Distance
		duration := data.Routes[0].Legs[0].Duration

		fmt.Println("Polyline:", polyline)
		fmt.Println("Distance Text:", distance.Text)
		fmt.Println("Distance Value:", distance.Value)

		fmt.Println("Duration Text:", duration.Text)
		fmt.Println("Duration Value:", duration.Value)
	} else {
		fmt.Println("No route data available.")
	}
}
