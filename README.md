# envfile-to-github-secrets

## Getting Started

1. 必要なライブラリをインストール

Cバインディングしているので以下のライブラリをインストールする必要があります。

```shell script
brew install libsodium
```

2. `.env` に投入するシークレットを用意

3. 実行

```shell script
 go run main.go -owner=p1ass -repo=envfile-to-github-secrets
```