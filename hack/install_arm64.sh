#!/bin/bash
set -e

sudo systemctl stop vcgencmd_exporter || true

sudo curl -sL -o /usr/local/bin/vcgencmd_exporter_linux_arm64 https://github.com/jsiebens/vcgencmd_exporter/releases/download/v0.0.1/vcgencmd_exporter_linux_arm64
sudo chmod 755 /usr/local/bin/vcgencmd_exporter_linux_arm64

sudo mkdir -p /etc/vcgencmd_exporter.d

sudo tee /etc/vcgencmd_exporter.d/env >/dev/null <<EOF
DELAY=10
PORT=2113
VCGENCMD_BINARY=/usr/bin/vcgencmd
EOF

sudo tee /etc/systemd/system/vcgencmd_exporter.service >/dev/null <<EOF
[Unit]
Description="vcgencmd exporter"

[Service]
Type=exec
EnvironmentFile=/etc/vcgencmd_exporter.d/env
ExecStart=/usr/local/bin/vcgencmd_exporter_linux_arm64
KillMode=process
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF
sudo chmod 0600 /etc/systemd/system/vcgencmd_exporter.service

sudo systemctl daemon-reload
sudo systemctl enable vcgencmd_exporter.service
sudo systemctl restart vcgencmd_exporter.service
