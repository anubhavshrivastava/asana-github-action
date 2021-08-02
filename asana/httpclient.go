package asana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Call(url string, data interface{}, responseHolder interface{}, headers map[string]string, method string) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		err := fmt.Errorf("unable to marshal the request: %v", data)
		return err
	}

	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(dataJson))
	if err != nil {
		err := fmt.Errorf("unable to initiate create task request: %v\n", err)
		return err
	}

	for key, value := range headers {
		httpRequest.Header.Set(key, value)
	}

	client := http.Client{}

	response, err := client.Do(httpRequest)
	if err != nil {
		err := fmt.Errorf("failed to make request: %v\n", err)
		return err
	}

	if response.StatusCode != 201 && response.StatusCode != 200 {
		return fmt.Errorf("non-OK status code: %v", response.Body)
	}

	err = json.NewDecoder(response.Body).Decode(responseHolder)
	if err != nil {
		err := fmt.Errorf("unable to decode the response: %v", err)
		return err
	}
	return nil
}
