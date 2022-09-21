FROM ubuntu:18.04

RUN apt update
RUN apt install sysstat tcpdump iproute2 net-tools -y