# STAGE 1: build
FROM golang:1.16-alpine as build
 
WORKDIR /source

COPY go.mod ./
COPY go.sum ./
COPY vendor ./

COPY ticket /source/ticket
COPY domain /source/domain
COPY mail /source/mail
COPY space /source/space

RUN CGO_ENABLED=0 GOOS=linux go build -o /source/ticket/.temp/ticket /source/ticket/server/ticket
 
# STAGE 2: run
FROM scratch as run
 
ENV FTA_SERVICE_NAME ticket

COPY --from=build /source/ticket/.temp/ticket /ticket
 
CMD ["/ticket"]