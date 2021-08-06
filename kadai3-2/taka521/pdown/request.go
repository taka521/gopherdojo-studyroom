package pdown

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// rangeDownload は Range アクセスにより取得したデータを一時ファイルへ書き込みます。
func (d *downloader) rangeDownload(ctx context.Context, part, start, end int64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// nothing
	}

	// get range body
	body, err := d.getRangeBody(d.URL, start, end)
	if err != nil {
		return fmt.Errorf("part %d request error: %w", part, err)
	}
	defer body.Close()

	// temporary file name
	dirname := d.tmpDir
	filename := fmt.Sprintf("%s.%d", d.FileName, part)
	path := filepath.Join(dirname, filename)

	// create a temporary file
	output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to create %s in %s: %w", filename, dirname, err)
	}
	defer output.Close()

	if _, err = io.Copy(output, body); err != nil {
		return fmt.Errorf("failed to copy %s in %q: %w", filename, dirname, err)
	}
	d.tmpFiles[part] = path

	return nil
}

func (d downloader) getRangeBody(url string, start, end int64) (io.ReadCloser, error) {
	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// set range header
	req.Header.Add(hRange, fmt.Sprintf("bytes=%d-%d", start, end))

	// execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}

	return res.Body, err
}
