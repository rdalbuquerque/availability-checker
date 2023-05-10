package checker

import (
	"net/http"
)

type HttpChecker struct {
	URL string
}

func (c *HttpChecker) Check() (bool, error) {
	resp, err := http.Get(c.URL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

func (c *HttpChecker) Name() string {
	return c.URL
}
