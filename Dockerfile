FROM golang:1.24.4 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/myapp ./src/cmd/*

FROM scratch

COPY --from=build /out/myapp /myapp

COPY migrations /migrations

EXPOSE 8080
ENTRYPOINT ["/myapp"]