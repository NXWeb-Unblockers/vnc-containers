FROM archlinux:latest

RUN pacman -Syu --noconfirm --needed \
    tigervnc \
    openbox \
    alacritty \
    sudo \
    noto-fonts \
    && rm -rf /var/cache/pacman/pkg/*

# Create a non-root user
RUN useradd -m -s /bin/bash app
RUN echo "app:app" | chpasswd && \
    usermod -aG wheel app && \
    printf '%s\n' '%wheel ALL=(ALL:ALL) ALL' > /etc/sudoers.d/wheel && \
    chmod 0440 /etc/sudoers.d/wheel
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