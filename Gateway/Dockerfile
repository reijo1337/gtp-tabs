FROM golang:1.12 AS stage
ENV CGO_ENABLED 0
WORKDIR /Gateway
COPY . .
RUN go build -mod=vendor -o ./application ./*.go

FROM alpine:3.7
WORKDIR /Gateway
COPY --from=stage Gateway/application application
ENV PORT ":9090"
EXPOSE 9090
ENTRYPOINT [ "/Gateway/application" ]