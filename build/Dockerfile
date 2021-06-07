FROM golang:alpine as builder

ENV APP_USER rotator
ENV APP_HOME /rotator
ENV BIN rotator.app
ARG LDFLAGS

RUN addgroup -S ${APP_USER} && adduser -S ${APP_USER} -G ${APP_USER}
RUN mkdir -p ${APP_HOME} && chown -R ${APP_USER}:${APP_USER} ${APP_HOME}

WORKDIR ${APP_HOME}
USER $APP_USER

COPY go.mod go.sum ./
# RUN go mod download

COPY ./ ${APP_HOME}
RUN CGO_ENABLED=0 go build -v -o "${APP_HOME}/${BIN}" -ldflags "${LDFLAGS}" ./cmd/...

FROM alpine:latest

ENV APP_USER rotator
ENV APP_HOME /rotator
ENV BIN rotator.app
ENV CONFIG_PATH /etc/rotator.cfg

RUN addgroup -S ${APP_USER} && adduser -S ${APP_USER} -G ${APP_USER}
RUN mkdir -p ${APP_HOME} && chown -R ${APP_USER}:${APP_USER} ${APP_HOME}
WORKDIR $APP_HOME

COPY --chown=${APP_USER}:${APP_USER} --from=builder ${APP_HOME}/configs/config.toml /etc/rotator.cfg
COPY --chown=${APP_USER}:${APP_USER} --from=builder ${APP_HOME}/${BIN} ${APP_HOME}
# RUN chmod +x ${APP_HOME}/${BIN}

EXPOSE 8080
USER $APP_USER
CMD ["sh", "-c", "${APP_HOME}/${BIN} -config ${CONFIG_PATH}"]
