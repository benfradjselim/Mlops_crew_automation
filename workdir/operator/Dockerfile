FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /ruptura-operator .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /ruptura-operator /ruptura-operator
ENTRYPOINT ["/ruptura-operator"]
