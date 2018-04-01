package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// URLParam represents a single key value param to add to a request.
type URLParam struct {
	Key string
	Val string
}

func createGetRequest(baseURL string, isJSON bool, params []URLParam) (*http.Request, error) {
	URL, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	queryParams := URL.Query()

	for _, param := range params {
		queryParams.Set(param.Key, param.Val)
	}

	URL.RawQuery = queryParams.Encode()
	urlString := URL.String()

	req, err := http.NewRequest("GET", urlString, nil)

	if err != nil {
		return nil, err
	}

	if isJSON {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// Response represents a response from a GetRequest or a PostRequest.
type Response struct {
	Response *http.Response
	Error    error
}

// GetRequest performs a GetRequest, returning an error and
func GetRequest(url string, isJSON bool, params []URLParam, timeout time.Duration) chan Response {
	c := make(chan Response)
	request, err := createGetRequest(url, isJSON, params)
	if err != nil {
		// TODO: Log that shit
		c <- Response{nil, err}
		return c
	}

	client := http.Client{Timeout: timeout}
	go func() {
		res, err := client.Do(request)
		if err != nil {
			c <- Response{nil, err}
		} else {
			c <- Response{res, nil}
		}
	}()

	return c
}

// ParseResponseIntoJSON attemps to marshall and http.Response into obj.
func ParseResponseIntoJSON(res *http.Response, obj interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return err
	}
	// TODO: What if this throws an error?
	jsonErr := json.Unmarshal(body, obj)

	if jsonErr != nil {
		return jsonErr
	}

	return nil
}
