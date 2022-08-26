FROM golang:1.18-alpine

WORKDIR /home/srcstats/fetcher

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /fetch

CMD ["/fetch"]