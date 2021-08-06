package pdown

import (
	"net/http"
	"strconv"
)

const (
	hAcceptRanges  = "Accept-Ranges"
	hContentLength = "Content-Length"
	hRange         = "Range"
)

func canRangeAccess(url string) (bool, int64, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, 0, err
	}

	if res.Header.Get(hAcceptRanges) != "bytes" {
		return false, 0, nil
	}

	size, err := strconv.ParseInt(res.Header.Get(hContentLength), 10, 64)
	return true, size, err
}
