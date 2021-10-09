FROM golang:1.17.1-alpine AS build

COPY . /go/src/time_app
WORKDIR /go/src/time_app

RUN go mod download
RUN go build -o /time_app


FROM alpine

WORKDIR /

COPY --from=build /time_app /time_app

ENTRYPOINT ["./time_app"]
CMD ["--help"]
