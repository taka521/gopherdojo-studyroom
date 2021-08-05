package pdown

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/dustin/go-humanize"
)

type Input struct {
	URL     string // ダウンロード対象のURL
	DownDir string // ダウンロード先のディレクトリパス
}

type Downloader interface {
	Run(input Input) error
}

type downloader struct {
	url       string         // ダウンロードURL
	outDir    string         // 出力先ディレクトリ
	file      *os.File       // 出力ファイル
	totalSize int64          // ファイルサイズ
	procs     int64          // 並列数
	wg        sync.WaitGroup // 待ち合わせ用Group
}

// New は Downloader のインスタンスを生成し、返却します。
func New() Downloader {
	return &downloader{
		procs: int64(runtime.NumCPU()),
	}
}

func (d *downloader) Run(input Input) error {
	d.url = input.URL
	d.outDir = input.DownDir

	// Get file size and check Range access support.
	if size, err := getSizeAndRangeSupport(d.url); err != nil {
		return fmt.Errorf("%w", err)
	} else {
		log.Printf("download size: %s\n", humanize.Bytes(uint64(size)))
		d.totalSize = size
	}

	// Create output file.
	filePath := filepath.Join(d.outDir, getFileName(d.url))
	if f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return fmt.Errorf("faild to create download file: %w", err)
	} else {
		d.file = f
		log.Printf("download to: %q\n", filePath)
	}

	// Start parallels download.
	var start, end int64
	partialSize := d.totalSize / d.procs
	for part := int64(0); part < d.procs; part++ {
		end = start + partialSize
		d.wg.Add(1)
		go d.writeRange(part, start, end)
		start = end
	}

	d.wg.Wait()
	return nil
}

// writeRange は Range アクセスにより取得したデータを出力先へ書き込みます。
func (d *downloader) writeRange(part, start, end int64) {
	var written int64

	body, size, err := d.getRangeBody(start, end)
	if err != nil {
		log.Fatalf("Part %d request error: %s\n", part, err.Error())
	}
	defer body.Close()
	defer d.wg.Done()

	buf := make([]byte, 4*1024)
	for {
		nr, er := body.Read(buf)
		if nr > 0 {
			// 出力先ファイルへ書き込み
			nw, err := d.file.WriteAt(buf[0:nr], start)

			// 書き込みに失敗していないか
			if err != nil {
				log.Fatalf("Part %d occured error: %s.\n", part, err.Error())
			}

			// リクエストボティから読み込んだデータが、全てファイルに書き込まれているか
			if nr != nw {
				log.Fatalf("Part %d occured error of short writiing.\n", part)
			}

			// 次の書き込み位置を調整 & 総書き込み量を変更
			start = int64(nw) + start
			if nw > 0 {
				written += int64(nw)
			}
		}

		// EOF 判定
		if er != nil {
			if errors.Is(er, io.EOF) {
				if size != written {
					log.Printf("Part %d unfinished.\n", part)
				}
				break
			}
			log.Printf("Part %d occured error: %s\n", part, er.Error())
		}
	}
}

// getRangeBody は指定されたバイト数で Range アクセスを行い、レスポンスボディ と サイズ(byte数) を返却します。
func (d *downloader) getRangeBody(start int64, end int64) (io.ReadCloser, int64, error) {
	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("%w", err)
	}

	// Range ヘッダーを付与し、リクエスト送信
	req.Header.Add(hRange, fmt.Sprintf("bytes=%d-%d", start, end))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("%w", err)
	}

	size, err := strconv.ParseInt(res.Header[hContentLength][0], 10, 64)
	return res.Body, size, err
}
