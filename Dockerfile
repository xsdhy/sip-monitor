FROM node:20.0 as web

ARG npm_registry

WORKDIR /app
COPY ./web .
RUN npm install --registry=${npm_registry} && npm run build

FROM golang:1.23 as golang

WORKDIR /app
COPY . .
COPY --from=web /app/build /app/web/build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o sbc

FROM alpine:latest

WORKDIR /app
COPY --from=golang /app/sbc .


EXPOSE 9060
EXPOSE 9059

CMD ["./sbc"]