FROM golang:alpine

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
COPY ./main .

# Command to run when starting the container
ENTRYPOINT [ "/dist/main" ]