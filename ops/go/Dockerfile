FROM golang:1.24.1-alpine as builder

ENV PATH /usr/local/go/bin:$PATH
ENV GOLANG_VERSION 1.24.1

RUN apk --no-cache add gcc musl-dev

WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

COPY . .
# COPY .env.example .env

RUN if [ ! -f go.mod ]; then \
        go mod init github.com/vnlab/makeshop-payment && \
        go mod tidy; \
    else \
        go mod download; \
    fi


RUN go build -trimpath -ldflags "-w -s" -o ./main ./src/main.go

FROM alpine3.21 as runner

RUN apk update && \
    apk add --no-cache shadow && \
    useradd -m appuser && \
    rm -f /usr/bin/gpasswd /usr/bin/passwd /usr/bin/chfn /sbin/unix_chkpwd /usr/bin/expiry /usr/bin/chage /usr/bin/chsh && \
    rm -rf /var/cache/apk/*

USER appuser

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]

FROM golang:1.24-alpine3.21 as dev

ENV CGO_ENABLED 0
ENV GO111MODULE auto
ENV GOLANG_VERSION 1.24.1
# don't auto-upgrade the gotoolchain
# https://github.com/docker-library/golang/issues/472
ENV GOTOOLCHAIN=local

RUN apk --no-cache add git

WORKDIR /app
COPY . /app

RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.1 && \
    go install github.com/pressly/goose/v3/cmd/goose@v3.7.0 && \
    go install github.com/golang/mock/mockgen@v1.6.0 && \
    go install github.com/cosmtrek/air@v1.27.3 && \
    go install github.com/99designs/gqlgen@v0.17.24 && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.3

CMD ["sh", "-c", "air -c /app/.air.toml"]