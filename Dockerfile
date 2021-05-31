FROM golang:1.16-alpine AS BASE

COPY . /code
WORKDIR /code
RUN go build -o /app/bot

FROM golang:1.16-alpine

RUN apk add --no-cache vim git bash gnupg
RUN go get github.com/github/hub

COPY --from=BASE /app/bot /bot
COPY /scripts /scripts
COPY ./entrypoint.sh /entrypoint.sh

CMD ["/entrypoint.sh"]
