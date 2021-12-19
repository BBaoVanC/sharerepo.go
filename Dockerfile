FROM golang:1.17-alpine AS build

COPY . /app
WORKDIR /app

RUN go get github.com/GeertJohan/go.rice/rice
RUN go build -o sharerepo
RUN rice append --exec sharerepo


FROM alpine:latest
COPY --from=build /app/sharerepo /sharerepo
EXPOSE 8080
CMD [ "/sharerepo" ]
