# STAGE 1: build
FROM golang:1.16-alpine as build
 
WORKDIR /source

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ticket /source/ticket
COPY domain /source/domain
COPY mail /source/mail
COPY space /source/space
COPY webserver /source/webserver

RUN CGO_ENABLED=0 GOOS=linux go build -o /source/webserver/.temp/webserver /source/webserver/server/webserver
 
# STAGE 2: run
FROM scratch as run
 
ENV FTA_SERVICE_NAME webserver

COPY --from=build /source/webserver/.temp/webserver /webserver
 
CMD ["/webserver"]
