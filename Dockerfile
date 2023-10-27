################
# BUILD BINARY #
################
# golang:1.14.6-alpine3.12
FROM golang:1.14.6-alpine3.12 as builder


# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/airpaxpricetime
COPY . .

RUN echo $PWD && ls -lah

# Fetch dependencies.
# RUN go get -d -v
#RUN go mod download
#RUN go mod verify

# CMD go build -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o ./airpaxpricetime .

#####################
# MAKE SMALL BINARY #
#####################
FROM alpine:3.12

RUN apk update && apk add --no-cache tzdata
ENV TZ=UTC

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the executable.
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/airpaxpricetime/airpaxpricetime /app
