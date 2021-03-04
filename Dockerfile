FROM golang:1.14
# Create required dirs and copy files
RUN mkdir -p /mittens
COPY ./ /mittens/
WORKDIR /mittens
# Run unit tests & build app
RUN make unit-tests

FROM alpine:3.7

# Create a group and user
RUN addgroup -g 1000 mittens && \
    adduser -D -u 1000 -G mittens mittens

# Layout folders
RUN mkdir /app && chown -R mittens:mittens /app

# Run as not root
USER $APP_USER

# Set workdir
WORKDIR /app

COPY --from=0 /mittens/mittens /app/mittens
ENTRYPOINT ["/app/mittens"]
