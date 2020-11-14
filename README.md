# envfile-to-github-secrets

## Getting Started

1. 必要なライブラリをインストール

```shell script
# Cバインディングしているので以下のライブラリをインストールする必要がある
brew install libsodium
```

2. `.env` に投入するシークレットを用意

3. 実行

```shell script
 go run main.go -owner=p1ass -repo=envfile-to-github-secrets
```

## License

以下のコードを参考にさせていただきました。
- https://qiita.com/kazz187/items/aa9885bb968722ac9b1d

