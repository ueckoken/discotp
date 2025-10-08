FROM golang:1.25.2 AS builder

ENV CGO_ENABLED=0

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN make build

FROM gcr.io/distroless/static-debian11:nonroot AS runner

COPY --from=builder --chown=nonroot:nonroot /app/discotp /app/discotp

ENTRYPOINT [ "/app/discotp" ]
