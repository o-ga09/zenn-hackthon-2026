#API用コンテナに含めるバイナリを作成するコンテナ
FROM golang:1.25.5-alpine3.23 AS deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -trimpath -ldflags "-w -s" -o main ./cmd/api/main.go

#-----------------------------------------------
#API デプロイ用コンアテナ
FROM ubuntu:22.04 AS deploy

RUN apt update
RUN apt-get install -y ca-certificates openssl

EXPOSE "8080"

COPY --from=deploy-builder /app/main .

CMD ["./main"]

#-----------------------------------------------
#ローカル開発環境で利用するホットリロード環境
FROM golang:1.25 AS dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go install github.com/air-verse/air@latest
CMD ["air"]
