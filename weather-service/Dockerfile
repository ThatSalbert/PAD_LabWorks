FROM golang:1.21-alpine

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o weatherservice .

CMD ["/app/weatherservice", "--host=0.0.0.0", "--port=$WEATHER_PORT"]