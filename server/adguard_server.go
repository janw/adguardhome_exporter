package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AdGuardServer server information
type AdguardServer struct {
	Url        string
	Username   string
	Password   string
	HTTPClient http.Client
}

// SendRequest send requests to endpoints
func (ad AdguardServer) SendRequest(api string, jsonStruct interface{}) error {

	url := strings.TrimRight(ad.Url, "/") + "/" + api

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if ad.Username != "" || ad.Password != "" {
		request.SetBasicAuth(ad.Username, ad.Password)
	}
	request.Header.Add("Accept", "application/json")

	response, err := ad.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error: status code %d from server", response.StatusCode)
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonStruct)
	if err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}
	return nil

}
