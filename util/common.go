package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func IpInfoLookup(token string, ip string) (IpInfoResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://ipinfo.io/%s", ip), nil)
	if err != nil {
		return IpInfoResponse{}, err
	}
	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	client := httpClient()

	resp, err := client.Do(req)
	if err != nil {
		return IpInfoResponse{}, err
	}

	defer resp.Body.Close()
	rBody, _ := ioutil.ReadAll(resp.Body)

	var v IpInfoResponse
	if err := json.Unmarshal(rBody, &v); err != nil {
		return IpInfoResponse{}, err
	}

	return v, nil
}

func httpClient() *http.Client {
	// reuse whenever possible
	client := &http.Client{Timeout: 15 * time.Second}
	return client
}
