# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.13 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o twitter-fn ./funcs/twitter/...

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/twitter-fn /twitter-fn

# Run the web service on container startup.
CMD ["/twitter-fn", "search", "NBA", "-o", "json", "-c", "10", "-S", "-p", "8080", \
	                "--twitter-api-key", "$TWITTER_API_KEY", \
	                "--twitter-api-secret-key", "$TWITTER_API_SECRET_KEY", \
	                "--twitter-access-token", "$TWITTER_ACCESS_TOKEN", \
	                "--twitter-access-token-secret", "$TWITTER_ACCESS_TOKEN_SECRET" \
]