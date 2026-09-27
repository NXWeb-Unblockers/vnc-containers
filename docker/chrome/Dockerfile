FROM debian:stable-slim

# Avoid prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages
RUN apt-get update && apt-get install -y \
    x11vnc \
    xvfb \
    chromium \
    --no-install-recommends \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -m -s /bin/bash app
WORKDIR /home/app

# Create directory for VNC (no password needed)
RUN mkdir -p /home/app/.vnc && \
    chown -R app:app /home/app/.vnc

# Create the start script
RUN echo '#!/bin/bash' > /home/app/start.sh && \
    echo 'Xvfb :1 -screen 0 1280x720x24 &' >> /home/app/start.sh && \
    echo 'export DISPLAY=:1' >> /home/app/start.sh && \
    echo 'sleep 2' >> /home/app/start.sh && \
    echo 'x11vnc -display :1 -forever -nopw &' >> /home/app/start.sh && \
    echo 'exec chromium --no-sandbox --window-position=0,0 --window-size=1280,720 https://google.com' >> /home/app/start.sh && \
    chmod +x /home/app/start.sh && \
    chown app:app /home/app/start.sh

USER app

# Expose VNC port
EXPOSE 5900

# Start the VNC server when container launches
CMD ["/home/app/start.sh"]