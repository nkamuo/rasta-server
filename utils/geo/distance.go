package geo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/nkamuo/rasta-server/model"
)

type DistanceMatrixElementEntry struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type DistanceMatrixElement struct {
	Duration DistanceMatrixElementEntry `json:"duration"`
	Distance DistanceMatrixElementEntry `json:"distance"`
}

type DistanceMatrixResponse struct {
	Rows []struct {
		Elements []DistanceMatrixElement `json:"elements"`
	} `json:"rows"`
}

func ResolveDistanceMatrix(origin *model.Location, destination *model.Location) (response *DistanceMatrixResponse, err error) {
	// origin := "37.7749,-122.4194"      // Latitude and longitude of origin
	// destination := "34.0522,-118.2437" // Latitude and longitude of destination

	iOrigin := origin.GetReference()
	iDestination := destination.GetReference()

	apiURL := "https://maps.googleapis.com/maps/api/distancematrix/json"
	params := url.Values{}
	params.Add("origins", iOrigin)
	params.Add("destinations", iDestination)
	params.Add("key", googleMapsAPIKey)

	fullURL := apiURL + "?" + params.Encode()

	data, err := http.Get(fullURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer data.Body.Close()

	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Rows) > 0 && len(response.Rows[0].Elements) > 0 {

	} else {
		return nil, errors.New("No data available.")
	}

	return response, nil
}
