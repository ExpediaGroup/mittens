FROM golang:1.14
# Create required dirs and copy files
RUN mkdir -p /mittens
COPY ./ /mittens/
WORKDIR /mittens
# Build app
RUN CGO_ENABLED=0 go build

FROM alpine:3.12

# Create a group and user
RUN addgroup -g 1000 mittens && \
    adduser -D -u 1000 -G mittens mittens

# Layout folders
RUN mkdir /app && chown -R mittens:mittens /app

# Run as not root
USER mittens

# Set workdir
WORKDIR /app

COPY --from=0 /mittens/mittens /app/mittens
ENTRYPOINT ["/app/mittens"]
