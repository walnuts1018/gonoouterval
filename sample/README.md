# sample

## custom golangci-lint実行方法

custom binaryをビルドします。

```sh
golangci-lint custom --destination tmp --name golangci-lint
```

customのgolangci-lintを実行します。

```sh
tmp/golangci-lint run .
```
