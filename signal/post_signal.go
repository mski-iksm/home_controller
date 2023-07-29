package signal

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func post_signal(api_secret string, fullurl string, post_body url.Values) int {
	req, err := http.NewRequest("POST", fullurl, strings.NewReader(post_body.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+api_secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		errLog.Println(err)
		return 500
	}
	defer res.Body.Close()

	// -------------------------------------------
	// Check http status code
	// -------------------------------------------

	if res.StatusCode != http.StatusOK {
		errLog.Printf("http status code: %d", res.StatusCode)
	}

	fmt.Printf("Sent order")
	return res.StatusCode
}
