FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Create a simple health check script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'echo "MCP Service starting..."' >> /app/start.sh && \
    echo 'while true; do' >> /app/start.sh && \
    echo '  echo "MCP Service running at $(date)"' >> /app/start.sh && \
    echo '  sleep 60' >> /app/start.sh && \
    echo 'done' >> /app/start.sh && \
    chmod +x /app/start.sh

EXPOSE 8080

USER 1001:1001

CMD ["/app/start.sh"]
