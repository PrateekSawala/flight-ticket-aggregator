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

RUN CGO_ENABLED=0 GOOS=linux go build -o /source/mail/.temp/mail /source/mail/server/mail
 
# STAGE 2: run
FROM scratch as run
 
ENV FTA_SERVICE_NAME space

COPY --from=build /source/mail/.temp/mail /mail
 
CMD ["/mail"]