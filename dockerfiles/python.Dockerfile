FROM python:3-alpine
RUN apk add make
ARG VERSION
RUN wget -O /usr/local/bin/brian https://github.com/stevens-tyr/brian/releases/download/$VERSION/brian
RUN chmod +x /usr/local/bin/brian