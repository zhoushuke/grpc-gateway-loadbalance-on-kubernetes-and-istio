FROM registry.sensetime.com/cloudos/golang:latest 

RUN apt-get update

RUN apt-get install -y ca-certificates
RUN apt-get install -y curl
RUN apt-get install -y telnet
RUN apt-get install -y dnsutils

ENV TZ=Asia/Shanghai

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN mkdir -p /src/bin

COPY process-service /src/bin

WORKDIR /src

ENV PATH "$PATH:/src/bin"

ENTRYPOINT ["process-service"]
