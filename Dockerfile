FROM golang:1.11
# Create required dirs and copy files
RUN mkdir -p /mittens
COPY ./ /mittens/
WORKDIR /mittens
# Run tests & build app
RUN make dependencies
RUN make test
RUN make build

FROM alpine:3.7
COPY --from=0 /mittens/build /app
ENTRYPOINT ["/app/mittens"]
