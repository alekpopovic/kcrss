FROM golang:1.21-alpine AS builder

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o controller main.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /workspace/controller .

USER 65532:65532

ENTRYPOINT ["/controller"]