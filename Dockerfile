FROM alpine:3.11 as builder

RUN apk add --no-cache git make go

WORKDIR /app

ARG TOK
ARG HOST
ARG DB_HOST
ARG DB_USER
ARG DB_PASS
ARG DB_NAME 
ARG CLIENT_ID 
ARG CLIENT_SECRET 
ARG SEC_KEY
ARG LYRICS_KEY 

ENV TOK $TOK
ENV HOST $HOST
ENV DB_HOST $DB_HOST
ENV DB_USER $DB_USER
ENV DB_PASS $DB_PASS
ENV DB_NAME $DB_NAME
ENV CLIENT_ID $CLIENT_ID
ENV CLIENT_SECRET $CLIENT_SECRET
ENV SEC_KEY $SEC_KEY
ENV LYRICS_KEY $LYRICS_KEY

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
ENV GO111MODULE=on

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

COPY go.mod go.sum ./
RUN go mod download
COPY . /app

#RUN GOOS=linux go build -o /app/sv .

EXPOSE 8080

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build /app/." --command=./spotify-views

#CMD ["/app/sv"]