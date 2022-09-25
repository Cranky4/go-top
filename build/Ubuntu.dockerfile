FROM ubuntu:18.04

ENV CONFIG_FILE /etc/app/config.toml
ENV GRPC_ADDR :9990
ENV BIN_FILE /opt/app/top

WORKDIR /opt/app

RUN apt update
RUN apt install sysstat tcpdump iproute2 net-tools -y

COPY ./bin/top ${BIN_FILE}
COPY ./configs/app.toml ${CONFIG_FILE}

RUN mkdir logs

CMD ${BIN_FILE} -config ${CONFIG_FILE} --grpc-addr=${GRPC_ADDR} > /opt/app/logs/app.log