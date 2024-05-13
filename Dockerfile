FROM golang:1.21.3-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY main.go .

RUN go build -o /simple-go-app main.go

FROM gcr.io/distroless/base-debian12

COPY --from=build /simple-go-app /simple-go-app

ENTRYPOINT ["/simple-go-app"]