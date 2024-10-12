FROM golang

# NOTE: Expecting environment variables to be provided to the image, not defining here.

# Install PostgreSQL client (psql)
# This is needed for interactively working with the database (e.g. migrations, seeding, etc.)
RUN apt-get update && apt-get install -y postgresql-client

WORKDIR /ambrosia-server

# Copy in build deps first to leverage docker layer caching
COPY go.mod .
COPY go.sum .
RUN go mod tidy

# Copy rest of source code and build
COPY . .
RUN go build .

CMD ["./ambrosia-server"]
EXPOSE 8080
