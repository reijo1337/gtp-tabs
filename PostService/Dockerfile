FROM golang:1.12 AS stage
ENV CGO_ENABLED 0
WORKDIR /PostService
COPY . .
RUN go build -mod=vendor -o ./application ./*.go

FROM alpine:3.7
WORKDIR /PostService
COPY --from=stage PostService/application application
ENV PORT ":9090"
EXPOSE 9090
ENTRYPOINT [ "/PostService/application" ]