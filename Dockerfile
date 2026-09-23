# Stage 1: compile the Go programs. The golang image is big, so we only use
# it to build, then copy the binaries into a small image (stage 2).
FROM golang:1.26-alpine AS build
WORKDIR /src

# Download dependencies first so Docker can cache this layer between builds
# (it only re-runs when go.mod/go.sum change).
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 makes fully static binaries that run on any Linux image.
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server \
 && CGO_ENABLED=0 go build -o /out/migrate ./cmd/migrate

# Stage 2: the image that actually runs. Templates, static files, and
# migrations are embedded in the binaries, so we only need the binaries.
FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/server /out/migrate /app/
EXPOSE 8080
CMD ["/app/server"]
