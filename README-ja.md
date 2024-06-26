# atgo

`atgo` はAtCoderにGo言語で参加するときに利用するためのツールです。

- [English](README.md)

## 生成AIへのルールへの適合

[AtCoder生成AI対策ルール - 20240607版](https://info.atcoder.jp/entry/llm-abc-rules-ja)

AtCoderでは、ABCにおいて生成AIの利用を制限するため上記のドキュメントに記載のルールを設けています。  
`atgo` では v0.0.4 でこのルールに対応する修正を行なったため、ABCへの参加においては v0.0.4 以降を利用してください。

## インストール

`atgo` は以下のコマンドでインストールします。

```
curl -sSfL https://raw.githubusercontent.com/meian/atgo/main/install | bash
```

ビルドバージョンを指定する場合は、以下のようにします。

```
curl -sSfL https://raw.githubusercontent.com/meian/atgo/main/install | bash -s -- --tag v0.0.1
```

[リリース](https://github.com/meian/atgo/releases) にビルド済のバイナリが用意されている場合はバイナリをダウンロードします。  
ビルド済バイナリがないOS/アーキテクチャの場合は、インストーラー内部で `go install` によってビルドされますが、この場合は Go 1.22 以上が必要です。

## 動作環境

- Linux
  - 開発はDockerのDebian上で行っていて、そこでの動作を確認しています

MacとWindowsは未検証ですが、インストールが可能であれば動作すると期待されます。  
(もし変な挙動があればPRをお願いします)

## 利用方法

### ワークスペースの設定

`atgo` はコマンドを実行した時のカレントディレクトリをワークスペースとして使用します。  
そのため、作業しているディレクトリ上へ移動してからコマンドを操作してください。

### 認証

AtCoderのユーザー名とパスワードを使って、以下のコマンドで認証情報を登録します：

```bash
$ atgo auth
ユーザー名: kitamin
パスワード: *******
```

一旦 `atgo auth` で認証情報を保存すると、他のコマンドを実行する際に自動的にログインセッションが引き継がれます。

認証情報は簡易な暗号化とともにローカルに保存されます。
もしローカルに保存された認証情報を削除したい場合は、以下のコマンドを実行してください。

```bash
$ atgo auth clear
```

### 問題のロード

コンテスト情報とそれに紐づく問題のリストを表示するには、以下のコマンドを実行します。

```bash
$ atgo contest [コンテストID］
```

問題の詳細を表示するには、以下のコマンドを実行します。

```bash
$ atgo task [タスクID］
```

ローカル環境に問題の回答用ファイルを作成するしたい場合は、以下のコマンドを実行します。

```bash
$ atgo task local-init [タスクID]
```

このコマンドで以下のファイルが作成されます。

- `main.go`
- `main_test.go`
- `go.mod`
- `go.sum`

一度作成したファイルはキャッシュされ、次回以降に同じ問題に対して `atgo task local-init` を実行すると、以前に途中まで編集していたファイルがロードされます。

### 回答コードとテストの作成

回答用のコードは `main.go` に実装します。  
テストコードの動作確認は、サンプルの入出力に対するテストが `main_test.go` に実装されているため、 `go test` で行えます。

### コードの提出

コードを提出するには、以下のコマンドを実行します。  
提出先の問題は、`atgo task local-init` で最後にローカルに用意した問題になります。

```bash
$ atgo submit
```

### その他のコマンド

#### 過去のコンテスト一覧を出力する

```bash
# 過去のコンテスト情報を取得して管理用DBに登録する
$ atgo contest load [abc or arc or agc or ahc]

# 過去のコンテスト一覧を出力する
# 事前に atgo contest load でコンテスト情報を取得しておく必要がある
$ atgo contest list [abc or arc or agc or ahc]
```

#### 管理用DBをクリアする

```bash
$ atgo workspace clean
```

上記コマンドで削除された DB には以下の情報が含まれています。  
`atgo task local-init` で作成されてキャッシュされたファイルは削除されません。

- コンテスト情報
- 問題情報
- コンテストと問題の関連付け

### ライセンス

このツールはMITライセンスで提供されています。
