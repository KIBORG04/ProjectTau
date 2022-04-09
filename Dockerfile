FROM golang:1.18-alpine
WORKDIR /usr/src/taustats

# If you enable this, then gcc is needed to debug your app
ENV CGO_ENABLED 0

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/ ./...

EXPOSE 8080

CMD ["ProjectTau"]