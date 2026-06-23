# testdata

## custom golangci-lint実行方法

リポジトリ直下で custom binary をビルドします。

```sh
golangci-lint custom --destination tmp --name golangci-lint
```

fixture に対して golangci-lint を実行します。

```sh
cd noouterval/testdata
../../tmp/golangci-lint run .
```
