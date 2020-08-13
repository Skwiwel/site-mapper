package mapping

import (
	"io"
	"net/http"
)

func getHTTPResponseFast(address, header string) (io.Reader, int, error) {
	client := &http.Client{}

	request, err := createRequest(address, header)
	if err != nil {
		return nil, 0, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	return response.Body, response.StatusCode, nil
}

func createRequest(address, header string) (*http.Request, error) {
	request, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", header)
	return request, nil
}
