qiita2crowi
===========

:us: [:jp:](./docs/README_ja.md)

Migration tool from [Qiita:Team](https://teams.qiita.com/) to [Crowi](http://site.crowi.wiki/)

## Usage

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

The option of `-access-token` and `-crowi-url` must be specified to run.

```console
$ cat exported-data-from-qiita.json \
    | qiita2crowi \
    -access-token="abcdefghijklmnopqrstuvwxyz=" \
    -crowi-url="http://your.crowi.url" \
    -page-path="/qiita/pages" \
    -qiita-access-token="abcdefghijklmnopqrstuvwxyz="
```

## Installation

```console
$ go get github.com/b4b4r07/qiita2crowi
```

## License

MIT

## Author

b4b4r07
