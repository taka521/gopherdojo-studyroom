package pdown

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cheggaaa/pb/v3"
	"golang.org/x/sync/errgroup"
)

type Downloader interface {
	Run(ctx context.Context) error
}

type downloader struct {
	URL      string // ダウンロードURL
	FileName string // DLファイル名
	FileSize int64  // ファイルサイズ
	OutDir   string // 出力先ディレクトリ
	Proc     int64  // 並列数

	tmpDir   string           // 一時ディレクトリ
	tmpFiles map[int64]string // 一時ファイルのリスト
}

// New は Downloader のインスタンスを生成し、返却します。
func New(url, dir string) Downloader {
	// create a temporary directory
	tmpDir, err := os.MkdirTemp("", "pdown")
	if err != nil {
		log.Fatalf("failed to create a temporary directory: %s\n", err.Error())
	}

	return &downloader{
		URL:      url,
		FileName: filepath.Base(url),
		FileSize: 0,
		OutDir:   dir,
		Proc:     int64(runtime.NumCPU()),
		tmpDir:   tmpDir,
		tmpFiles: make(map[int64]string, 0),
	}
}

func (d *downloader) Run(ctx context.Context) error {
	defer d.cleaning()

	if err := d.Before(); err != nil {
		return err
	}

	if err := d.Download(ctx); err != nil {
		return err
	}

	if err := d.After(); err != nil {
		return err
	}

	return nil
}

// Before は分割ダウンロード可能であるかを検証します。
func (d *downloader) Before() error {
	supported, size, err := canRangeAccess(d.URL)
	if err != nil {
		return err
	}

	if !supported {
		return fmt.Errorf("range access in %q not supported", d.URL)
	}

	d.FileSize = size // save file size
	return nil
}

// Download は対象ファイルの分割ダウンロードを行います。
func (d downloader) Download(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	partialSize := d.FileSize / d.Proc
	for i := int64(0); i < d.Proc; i++ {
		part := i
		start := i * partialSize
		end := start + partialSize - 1
		if i == d.Proc-1 {
			end = d.FileSize
		}

		wg.Go(func() error {
			return d.rangeDownload(ctx, part, start, end)
		})
	}

	return wg.Wait()
}

// After は複数の一時ファイルを、単一のファイルにマーシします。
func (d downloader) After() error {
	filename := filepath.Join(d.OutDir, d.FileName)
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create a file: %w", err)
	}
	defer out.Close()

	// progress bar
	bar := pb.New64(d.FileSize)
	bar.Start()

	for i := int64(0); i < d.Proc; i++ {
		f, err := os.Open(d.tmpFiles[i])
		if err != nil {
			return fmt.Errorf("failed to open %q: %w", d.tmpFiles[i], err)
		}

		proxy := bar.NewProxyReader(f)
		_, _ = io.Copy(out, proxy)

		_ = f.Close()

		if err := os.Remove(d.tmpFiles[i]); err != nil {
			return fmt.Errorf("failed to remove %q: %w", d.tmpFiles[i], err)
		}
	}

	bar.Finish()
	return nil
}

func (d downloader) cleaning() {
	_ = os.RemoveAll(d.tmpDir)
}
