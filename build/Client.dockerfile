FROM alpine:3.16.2

ENV CONFIG_FILE /etc/app/client.toml
ENV GRPC_ADDR "ubuntu-testing:9990"
ENV BIN_FILE "/opt/app/client"
ENV N_PARAM 5
ENV M_PARAM 15

WORKDIR /opt/app

COPY ./bin/client ${BIN_FILE}
COPY ./configs/client.toml ${CONFIG_FILE}

RUN mkdir logs

CMD ${BIN_FILE} --config=${CONFIG_FILE} --grpc-addr=${GRPC_ADDR} --m=${M_PARAM}  --n=${N_PARAM}