FROM golang

# NOTE: Expecting environment variables to be provided to the image, not defining here.

WORKDIR /ambrosia-server
COPY . .
RUN go mod tidy
RUN go build .
CMD ["./ambrosia-server"]
EXPOSE 8080
