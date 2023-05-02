FROM golang:1.20

ENV GO111MODULE=on
WORKDIR /app/backend

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD ["go", "run", "./cmd/mockData/main.go"]
