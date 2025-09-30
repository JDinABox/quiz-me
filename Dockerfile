FROM oven/bun:latest AS builder-bun
WORKDIR /quiz-me
COPY ./package.json ./package.json
COPY ./bun.lock ./bun.lock
RUN bun install
COPY ./tsconfig.json ./tsconfig.json
COPY ./vite.config.js ./vite.config.js
COPY ./web ./web
RUN bun run build

FROM golang:alpine AS builder-go
WORKDIR /go/src/github.com/JDinABox/quiz-me
COPY go.* ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./web ./web
COPY ./internal ./internal
COPY --from=builder-bun /quiz-me/web/dist ./web/dist
COPY *.go ./
RUN go tool templ generate
ENV GOEXPERIMENT=jsonv2
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags="-s -w" -o ./cmd/quiz-me/quiz-me.so ./cmd/quiz-me


FROM alpine:latest

RUN apk --no-cache -U upgrade \
    && apk --no-cache add --upgrade ca-certificates \
    && ARCH=$(uname -m) && wget -O /bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_${ARCH} \
    && chmod +x /bin/dumb-init

COPY --from=builder-go /go/src/github.com/JDinABox/quiz-me/cmd/quiz-me/quiz-me.so /bin/quiz-me.so

ENTRYPOINT ["/bin/dumb-init", "--"]
CMD ["/bin/quiz-me.so"]