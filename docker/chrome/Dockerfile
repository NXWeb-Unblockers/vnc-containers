FROM debian:stable-slim

# Avoid prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages
RUN apt-get update && apt-get install -y \
    tigervnc-standalone-server \
    openbox \
    chromium \
    --no-install-recommends \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -m -s /sbin/nologin app
WORKDIR /home/app

COPY --chown=app:app rc.xml /home/app/.config/openbox/rc.xml

# Hides Chromium's minimize/maximize/close buttons
RUN mkdir -p /home/app/.config/gtk-3.0 && \
    printf '[Settings]\ngtk-decoration-layout=:\n' > /home/app/.config/gtk-3.0/settings.ini && \
    chown -R app:app /home/app/.config

# Create the start script
RUN echo '#!/bin/bash' > /home/app/start.sh && \
    echo 'rm -f /tmp/.X1-lock /tmp/.X11-unix/X1' >> /home/app/start.sh && \
    echo 'Xvnc :1 -geometry 1280x720 -depth 24 -rfbport 5900 -SecurityTypes None &' >> /home/app/start.sh && \
    echo 'export DISPLAY=:1' >> /home/app/start.sh && \
    echo 'sleep 2' >> /home/app/start.sh && \
    echo 'openbox &' >> /home/app/start.sh && \
    echo 'exec chromium --no-sandbox --start-maximized https://google.com' >> /home/app/start.sh && \
    chmod +x /home/app/start.sh && \
    chown app:app /home/app/start.sh

USER app

# VNC port
EXPOSE 5900

# Start the VNC server when container launches
CMD ["/home/app/start.sh"]