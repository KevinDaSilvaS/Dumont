FROM mariadb as builder

RUN apt-get update

RUN apt-get install golang -y

RUN mkdir dumont

WORKDIR /dumont

COPY . .

RUN go build -o dumont .

FROM mariadb

COPY --from=builder dumont/dumont .

CMD ./dumont