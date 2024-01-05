FROM golang:1.21.5-bullseye AS build

RUN apt-get update && apt-get install -y git

WORKDIR /app

RUN echo cart-service

RUN git clone https://github.com/akshay0074700747/cart-service-grpc.git .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o bin/cart-service

COPY /cmd/.env /app/cmd/bin/

FROM busybox:latest

WORKDIR /cart-service

COPY --from=build /app/cmd/bin/cart-service .

COPY --from=build /app/cmd/bin/.env .

EXPOSE 50006

CMD ["./cart-service"]