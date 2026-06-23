# gonoouterval

## golangci-lint v2

### ビルド

まずは、以下のコマンドをリポジトリ直下で実行し、custom golangci-lint binary をビルドします。

```sh
golangci-lint custom --name custom-gcl
```

### 利用

`.golangci.yml` で `gonoouterval` を有効化します。

```yaml
version: "2"
linters:
  default: none
  enable:
    - gonoouterval
  settings:
    custom:
      gonoouterval:
        type: module
        description: Check for outer scoped values when same-typed inner values exist.
        settings:
          type: path/to/pkg.Type
```

`settings.type` には、検査対象の型を `path/to/pkg.Type` 形式で指定します。

custom binary を実行します。

```sh
./custom-gcl run ./...
```

## 単体 analyzer として使う

golangci-lint を経由せず、単体 command としても実行できます。

```sh
go build ./cmd/gonoouterval
gonoouterval -type path/to/pkg.Type ./...
```
