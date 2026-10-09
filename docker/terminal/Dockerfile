FROM debian:stable-slim

# Avoid prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages
RUN apt-get update && apt-get install -y \
    tigervnc-standalone-server \
    openbox \
    alacritty \
    sudo \
    --no-install-recommends \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -m -s /bin/bash app
RUN echo "app:app" | sudo chpasswd
RUN usermod -aG sudo app
WORKDIR /home/app

# Create the start script
RUN echo '#!/bin/bash' > /home/app/start.sh && \
    echo 'rm -f /tmp/.X1-lock /tmp/.X11-unix/X1' >> /home/app/start.sh && \
    echo 'Xvnc :1 -geometry 1280x720 -depth 24 -rfbport 5900 -SecurityTypes None &' >> /home/app/start.sh && \
    echo 'export DISPLAY=:1' >> /home/app/start.sh && \
    echo 'sleep 2' >> /home/app/start.sh && \
    echo 'openbox &' >> /home/app/start.sh && \
    echo 'exec alacritty' >> /home/app/start.sh && \
    chmod +x /home/app/start.sh && \
    chown app:app /home/app/start.sh

USER app

# VNC port
EXPOSE 5900

# Start the VNC server when container launches
CMD ["/home/app/start.sh"]