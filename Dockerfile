FROM golang as builder

RUN mkdir dumont

RUN cd dumont

COPY . .

RUN go build

RUN apt update
RUN apt install -y mariadb-client

#FROM mariadb

#FROM ubuntu


#RUN mkdir dumont
#COPY --from=builder / /dumont
#RUN cd dumont
#
CMD ./dumont