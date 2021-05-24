FROM alpine:3.11 as builder

RUN apk add --no-cache git make go

WORKDIR /app

ARG HOST
ENV HOST $HOST

ARG PORT
ENV PORT $PORT

ARG GO_ENV
ENV GO_ENV $GO_ENV

ARG MASTER_KEY
ENV MASTER_KEY $MASTER_KEY

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
#ENV GO111MODULE=on

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

COPY go.mod go.sum ./
RUN go mod download
COPY . /app

RUN GOOS=linux go build -o /app/sv /app/cmd/spotify-views

EXPOSE 8080

CMD ["./sv"]