FROM golang:latest

WORKDIR /app

COPY ./ ./

EXPOSE 8080

RUN go mod download
RUN go build -o weather-averager ./cmd/open-weather-map/main.go

CMD [ "./weather-averager" ]