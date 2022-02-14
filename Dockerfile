FROM registry.erda.cloud/retag/golang:1.16-alpine AS BASE

COPY . /code
WORKDIR /code
RUN go build -o /app/bot

FROM registry.erda.cloud/retag/golang:1.16-alpine

RUN apk add --no-cache vim git bash gnupg docker-cli make grep
RUN go get github.com/github/hub && go get github.com/mikefarah/yq/v4

COPY --from=BASE /app/bot /bot
COPY /scripts /scripts
COPY ./entrypoint.sh /entrypoint.sh

CMD ["/entrypoint.sh"]
