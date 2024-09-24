# ベースイメージとしてGoを使用
FROM golang:1.23-alpine

# 作業ディレクトリを作成
WORKDIR /app

# Go Modulesの依存関係をキャッシュするためにまずgo.modとgo.sumをコピー
COPY go.mod .
COPY go.sum .
RUN go install github.com/air-verse/air@latest
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# ポートをエクスポート
EXPOSE 8080

# ホットリロード用のエントリーポイント
CMD ["air", "-c", ".air.toml"]

