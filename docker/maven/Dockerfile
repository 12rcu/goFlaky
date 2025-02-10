FROM ubuntu:22.04
RUN apt update
RUN apt install git openjdk-11-jdk maven -y
RUN export JAVA_HOME=/usr/lib/jvm/openjdk-11
RUN chmod +x /usr/share/maven/bin/*

WORKDIR /app
COPY ./goFlaky /app/goFlaky
CMD ["./goFlaky"]