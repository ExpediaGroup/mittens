FROM golang:1.14
# Create required dirs and copy files
RUN mkdir -p /mittens
COPY ./ /mittens/
WORKDIR /mittens
# Run unit tests & build app
RUN make unit-tests

FROM scratch
COPY --from=0 /mittens/mittens /app/mittens
ENTRYPOINT ["/app/mittens"]
