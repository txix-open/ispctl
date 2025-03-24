FROM golang:1.24-alpine3.21 as builder
WORKDIR /build
ARG version
ENV version_env=$version
ARG app_name
ENV app_name_env=$app_name
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -ldflags="-X 'main.version=$version_env'" -o /main .

FROM alpine:3.21

RUN apk add --no-cache tzdata bash-completion jq
RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN echo "Europe/Moscow" > /etc/timezone

ARG app_name
ENV app_name_env=$app_name
COPY --from=builder main /usr/bin/$app_name_env
COPY /conf/config.yml /etc/$app_name_env/config.yml
COPY /static/autocomplete /etc/bash_completion.d/$app_name_env

CMD ["exec $app_name_env"]
