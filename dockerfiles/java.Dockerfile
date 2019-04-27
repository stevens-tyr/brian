FROM openjdk:13-alpine
RUN apk add make
COPY brian /usr/local/bin/
RUN chmod +x /usr/local/bin/brian