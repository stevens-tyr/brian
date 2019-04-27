FROM python:3-alpine
RUN apk add make
COPY brian /usr/local/bin/
RUN chmod +x /usr/local/bin/brian