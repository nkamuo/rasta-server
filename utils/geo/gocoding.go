package geo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/tidwall/gjson"
	// "strconv"
)

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type GeocodeResult struct {
	Results []struct {
		AddressComponents []AddressComponent `json:"address_components"`
		PlaceID           string             `json:"place_id"`
		FormattedAddress  string             `json:"formatted_address"`
		Geometry          struct {
			Location struct {
				Lat float32 `json:"lat"`
				Lng float32 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func ResolveGeocodingInfo(googleMapsAPIKey string, input string) (result *GeocodeResult, err error) {

	requestType := IdentifyInputType(input)

	if requestType == "UNKNOWN" {
		message := fmt.Sprintf("Could not resolve request type \"%s\"", input)
		return nil, errors.New(message)
	}

	geocodeURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?key=%s&%s=%s", googleMapsAPIKey, requestType, url.QueryEscape(input))

	resp, err := http.Get(geocodeURL)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
		return
	}

	// var result GeocodeResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error decoding the response:", err)
		return
	}

	if len(result.Results) > 0 {
		// fmt.Printf("RESULT: %#v", result)
		// location := result.Results[0].Geometry.Location
		// fmt.Printf("Formatted Address: %s\n", result.Results[0].FormattedAddress)
		// fmt.Printf("Latitude: %f\n", location.Lat)
		// fmt.Printf("Longitude: %f\n", location.Lng)
	} else {
		return nil, errors.New("Location not found")
	}

	return result, nil
}

func GetPlaceName(googleMapApiKey string, placeId string) (name string, err error) {

	var requestUrl = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?fields=name&place_id=%s&key=%s", placeId, googleMapApiKey)

	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		message := fmt.Sprintf("Error reading the response: %s", err)
		return "", errors.New(message)
	}

	bodyString := string(body)

	jsonData := gjson.Parse(bodyString)
	name = jsonData.Get("result").Get("name").String()
	return name, nil
}

func IdentifyInputType(input string) string {
	// Regular expressions for matching address, coordinates, and Google Place ID patterns
	addressPattern := `^[A-Za-z0-9\s.,-]+$`
	coordinatesPattern := `^-?\d+\.\d+,\s*-?\d+\.\d+$`
	placeIDPattern := `^ChIJ[0-9A-Za-z_-]+$`

	if matched, _ := regexp.MatchString(placeIDPattern, input); matched {
		return "place_id"
	} else if matched, _ := regexp.MatchString(coordinatesPattern, input); matched {
		return "latlng"
	} else if matched, _ := regexp.MatchString(addressPattern, input); matched {
		return "address"
	}

	return "UNKNOWN"
}

// func main() {
// 	// Test cases
// 	ResolveGeocodingInfo("Divinics Electrical Shop")    // Address
// 	ResolveGeocodingInfo("6.024519,7.084139")           // Coordinates (latitude, longitude)
// 	ResolveGeocodingInfo("ChIJ2eUgeAK6j4ARbn5u_wAGqWA") // Google Place ID
// }
