FROM golang:1.12 AS stage
ENV CGO_ENABLED 0
WORKDIR /GtpStorage
COPY . .
RUN go build -mod=vendor -o ./application ./*.go

FROM alpine:3.7
WORKDIR /GtpStorage
COPY --from=stage GtpStorage/application application
ENV PORT ":9090"
EXPOSE 9090
ENTRYPOINT [ "/GtpStorage/application" ]