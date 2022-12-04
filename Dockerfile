# https://docs.docker.com/language/golang/build-images/

FROM golang:1.19.3-alpine3.17

# Working directory
WORKDIR /usr/src/app

# Copy package to Working directory
COPY go.mod ./
COPY go.sum ./

# Download pakcage
RUN go mod download

# Copy all source code to Working directory
COPY . ./

# Build main.exe
RUN go build main.go

# Expose port
EXPOSE 5002

# Execute program
CMD [ "/main" ]
