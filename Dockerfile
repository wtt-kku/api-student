## Build
FROM golang:bullseye AS build


ADD . /go/src/run_service
WORKDIR /go/src/run_service

# RUN go get -u
# RUN go mod vendor
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /run_service

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /run_service /run_service

ENV TZ="Asia/Bangkok"

ENV ENV=prod
ENV HTTP_PORT=6000
ENV DB_HOST=thaicharoendb-do-user-8957810-0.b.db.ondigitalocean.com
ENV DB_PORT=25060
ENV DB_USER=doadmin
ENV DB_PASS=AVNS_qSdPKF_oArfQKvhn2oM
ENV DB_NAME=student_score
ENV SECRET_KEY=WTTxPetch
ENV ACCESS_TOKEN_EXP_TIME_MIN=6000

EXPOSE 6000


ENTRYPOINT ["/run_service"]
