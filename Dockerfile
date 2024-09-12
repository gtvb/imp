FROM go:latest

WORKDIR /app

COPY go.mod go.sum .

RUN go mod tidy

COPY . .

RUN go build . -o imp

ENTRYPOINT [ "imp" ]
