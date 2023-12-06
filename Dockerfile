FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . /app/
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo -ldflags '-extldflags "-static"' ./...

FROM scratch
COPY --from=build /go/bin/* /bin/
ENTRYPOINT ["/bin/mleh"]
