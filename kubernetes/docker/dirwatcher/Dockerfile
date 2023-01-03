FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# USER 405 in /etc/passwd is
#    guest:x:405:100:guest:/dev/null:/sbin/nologin
# So container runtime will assign GID 100 for the process
USER 405
