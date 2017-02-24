Spec
====

## 仕様的な部分

- このツールを使うにあたり、Qiita:Team にある記事を Dumping した JSON ファイルが必要です
  - Qiita:Team では記事のエクスポート機能があるのでそれを使ってください
- その JSON ファイルを標準出力から食わせると、パースしてよしなに Crowi API に投げてくれます
  - 実行すると、
  - まずは記事を作成します
    - そのとき JSON にある title 要素が Crowi の記事のパスとして使われます
    - `qiita2crowi` の `-page-path` オプションを使えばその記事の root path を設定できます (default: `/qiita`)
    - title にある `^`, `$`, `*`, `%`, `?`, `/` は全角に置換されます
    - Qiita 記事についているコメントは本文末尾に追記される形になります
  - 次に、画像をアップロードします
    - Qiita からアップロードされた画像ファイルは `qiita-image-store.s3.amazonaws.com` というホスト名の URL を持ちます
    - Qiita:Team 解約時に削除される可能性もあるので、これらをダウンロードして Crowi にアップロードします
  - 最後にアップロードした画像パスで Crowi 記事を書き換えます

## 実装的な部分

- 4 つの goroutine を起動して並列に処理します
  - 最適化用のオプションなどは提供していません

## 参考

- [Qiita::Team やめた - @kyanny's blog](http://blog.kyanny.me/entry/2015/07/30/020046)
