# ベースイメージの指定
FROM golang:1.23-alpine


# 作業ディレクトリの設定
WORKDIR /app

# Go modulesのキャッシュを活用するために必要なファイルをコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# ビルドコマンド
RUN go build -o main .

# 実行コマンド
CMD ["/app/main"]
