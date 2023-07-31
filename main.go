package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	// "strconv"
)

const googleMapsAPIKey = "AIzaSyBhgBfG2YQsF_CivgkwKP39AP_d-Q-2aEU"

type GeocodeResult struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func resolveGeocodingInfo(input string) {

	requestType := identifyInputType(input)

	if requestType == "UNKNOWN" {
		fmt.Println("Could not resolve request type \"%s\"", input)
	}

	geocodeURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?key=%s&%s=%s", googleMapsAPIKey, requestType, url.QueryEscape(input))

	resp, err := http.Get(geocodeURL)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
		return
	}

	var result GeocodeResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error decoding the response:", err)
		return
	}

	if len(result.Results) > 0 {
		location := result.Results[0].Geometry.Location
		fmt.Printf("Formatted Address: %s\n", result.Results[0].FormattedAddress)
		fmt.Printf("Latitude: %f\n", location.Lat)
		fmt.Printf("Longitude: %f\n", location.Lng)
	} else {
		fmt.Println("Location not found")
	}
}

func identifyInputType(input string) string {
	// Regular expressions for matching address, coordinates, and Google Place ID patterns
	addressPattern := `^[A-Za-z0-9\s.,-]+$`
	coordinatesPattern := `^-?\d+\.\d+,\s*-?\d+\.\d+$`
	placeIDPattern := `^ChIJ[0-9A-Za-z_-]+$`

	if matched, _ := regexp.MatchString(addressPattern, input); matched {
		return "address"
	} else if matched, _ := regexp.MatchString(coordinatesPattern, input); matched {
		return "latlng"
	} else if matched, _ := regexp.MatchString(placeIDPattern, input); matched {
		return "place_id"
	}

	return "UNKNOWN"
}

func main() {
	// Test cases
	resolveGeocodingInfo("Divinics Electrical Shop")    // Address
	resolveGeocodingInfo("6.024519, 7.084139")          // Coordinates (latitude, longitude)
	resolveGeocodingInfo("ChIJ2eUgeAK6j4ARbn5u_wAGqWA") // Google Place ID
}
