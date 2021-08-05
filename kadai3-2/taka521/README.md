# 課題3-2

分割ダウンローダを作ろう

## 仕様

* Rangeアクセスを用いる
* いくつかのゴルーチンでダウンロードしてマージする
* エラー処理を工夫する
    * `golang.org/x/sync/errgourp` パッケージなどを使ってみる
* キャンセルが発生した場合の実装を行う

## ビルド

```shell
$ make 
```

## 使い方

```shell
$ ./pdown.exe <DLディレクトリ> <URL>
```

## 参考にした実装

* https://github.com/jacklin293/golang-parallel-download-with-accept-ranges
