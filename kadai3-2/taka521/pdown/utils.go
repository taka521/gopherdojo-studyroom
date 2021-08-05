package pdown

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
)

const (
	hAcceptRanges  = "Accept-Ranges"
	hContentLength = "Content-Length"
	hRange         = "Range"
)

// getSizeAndRangeSupport はダウンロード対象のファイルサイズ取得および、Range アクセス可能であるかを検証します。
// ファイルサイズの取得に失敗したり、Range アクセス不可の場合は error を返却します。
//
// なお、本処理は以下のコードを参考にしています。
//
// 	https://github.com/jacklin293/golang-parallel-download-with-accept-ranges/blob/688a62221cd0f0321754c12363c2ec562d8a63ee/main.go#L179
func getSizeAndRangeSupport(url string) (size int64, err error) {
	// ヘッダーだけ欲しいので HEAD アクセス
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("リクエストに失敗しました: %w", err)
	}

	acceptRanges, supported := res.Header[hAcceptRanges]
	if !supported || (supported && acceptRanges[0] != "bytes") {
		return 0, errors.New("doesn't support range access")
	}

	size, err = strconv.ParseInt(res.Header[hContentLength][0], 10, 64)
	return
}

func getFileName(path string) string {
	return filepath.Base(path)
}
