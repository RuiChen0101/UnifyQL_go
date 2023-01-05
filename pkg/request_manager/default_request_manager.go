package fetch_proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type DefaultRequestManager struct {
	c *http.Client
}

func (fp *DefaultRequestManager) Request(id string, url string, uqlPayload string) ([]byte, error) {
	resp, err := fp.c.Post(url, "text/plain", bytes.NewBuffer([]byte(uqlPayload)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []byte("[]"), nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RequestManager: %s response %d when executing %s", url, resp.StatusCode, uqlPayload)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return body, nil
}
