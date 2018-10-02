qiita2crowi
===========

[:us:](../README.md) :jp:

[Qiita:Team](https://teams.qiita.com/) から [Crowi](http://site.crowi.wiki/) に移行するためのツールです

## 使い方

```console
$ ./qiita2crowi -h
Usage of ./qiita2crowi:
  -access-token string
        Crowi's access token
  -crowi-url string
        Your Crowi base URL
  -page-path string
        Default page path (default "/qiita")
  -qiita-access-token string
        Qiita's access token
```

実行するには `-access-token` と `-crowi-url` の指定が必須です。それぞれ、Crowi のアクセストークンとベース URL が必要です。

```console
$ cat exported-data-from-qiita.json \
    | qiita2crowi \
    -access-token="abcdefghijklmnopqrstuvwxyz=" \
    -crowi-url="http://your.crowi.url" \
    -page-path="/qiita/pages" \
    -qiita-access-token="abcdefghijklmnopqrstuvwxyz="
```

### 仕様

[こちら](spec_ja.md)

## インストール

```console
$ go get github.com/b4b4r07/qiita2crowi
```

## License

MIT

## Author

b4b4r07
