FROM gcc:8
COPY brian /usr/local/bin/
RUN chmod +x /usr/local/bin/brian