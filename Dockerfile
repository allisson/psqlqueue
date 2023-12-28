#### development stage
FROM golang:1.21 AS build-env

# set envvar
ENV CGO_ENABLED=0
ENV GOOS=linux

# set workdir
WORKDIR /code

# get project dependencies
COPY go.mod /code/
RUN go mod download

# copy files
COPY . /code

# generate binary
RUN go build -ldflags="-s -w" -o ./psqlqueue ./cmd/psqlqueue

#### final stage
FROM gcr.io/distroless/base:nonroot
COPY --from=build-env /code/psqlqueue /
ENTRYPOINT ["/psqlqueue"]
