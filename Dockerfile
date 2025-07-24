FROM registry.erda.cloud/retag/golang:1.24-alpine AS BASE

COPY . /code
WORKDIR /code
RUN go build -o /app/bot

FROM registry.erda.cloud/retag/golang:1.24-alpine

RUN apk add --no-cache vim git bash gnupg docker-cli make grep
RUN go install github.com/github/hub@latest && go install github.com/mikefarah/yq/v4@latest

COPY --from=BASE /app/bot /bot
COPY /scripts /scripts
COPY ./entrypoint.sh /entrypoint.sh

CMD ["/entrypoint.sh"]
