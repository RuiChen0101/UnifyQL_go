package fetch_proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DefaultFetchProxy struct {
	c *http.Client
}

func (fp *DefaultFetchProxy) Request(url string, uqlPayload string) ([]byte, error) {
	resp, err := fp.c.Post(url, "text/plain", bytes.NewBuffer([]byte(uqlPayload)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []byte("[]"), nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FetchProxy: %s response %d when executing %s", url, resp.StatusCode, uqlPayload)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return body, nil
}
