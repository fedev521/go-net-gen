# build stage
FROM golang:1.19 AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /go/bin/app /go/src/app/cmd/app
RUN --mount=type=cache,target=/root/.cache/go-build \
    go test -v
RUN --mount=type=cache,target=/root/.cache/go-build \
    go vet -v

# package stage
FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=builder /go/bin/app .
COPY configs .
CMD ["/app", "-p", "8000"]
