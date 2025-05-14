FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /gin_executor ./cmd/gin_executor

CMD ["/gin_executor"]