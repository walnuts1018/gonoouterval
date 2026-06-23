# gonoouterval

※ Linterの実装自体は [motemen](https://github.com/motemen)さんの[go-statictools](https://github.com/motemen/go-statictools)をそのまま利用しています。

## golangci-lint v2

### ビルド

custom golangci-lint binary をビルドする必要があります。

まず、`.custom-gcl.yml`ファイルを以下の内容で作成します。

```yaml
version: v2.12.2
plugins:
  - module: github.com/walnuts1018/gonoouterval
    import: "github.com/walnuts1018/gonoouterval"
    version: latest
```

その後、`golangci-lint custom`を使って、plugin入りのgolangci-lintをビルドします。

```bash
golangci-lint custom --name custom-gcl
```

### 利用

`.golangci.yml` で `gonoouterval` を有効化します。
`settings.type` には、検査対象の型を `path/to/pkg.Type` 形式で指定します。

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

custom golangci-lint を実行します。

```sh
./custom-gcl run ./...
```

## 単体 analyzer として使う

golangci-lint を経由せず、単体 command としても実行できます。

```sh
go build ./cmd/gonoouterval
gonoouterval -type path/to/pkg.Type ./...
```
