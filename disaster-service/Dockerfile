FROM golang:1.21-alpine

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o disasterservice .

CMD ["/app/disasterservice", "--host=0.0.0.0", "--port=$DISASTER_PORT"]