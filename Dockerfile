# build stage
FROM golang:1.15.3-alpine3.12 AS build-dev
RUN mkdir /app                        
ADD . /app/                        
WORKDIR /app                         
RUN go mod download
ENV CGO_ENABLED=0
# RUN CGO_ENABLED=0 GOOS=linux go build
RUN go test -v
# CMD ["go", "test", "-v"]
