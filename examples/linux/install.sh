#!/usr/bin/env bash

wget -q https://github.com/welcome-ad-infernum/mg/releases/latest/download/mg_linux_amd64 -O /usr/local/bin/mg
chmod +x /usr/local/bin/mg

cat << EOF > /etc/systemd/system/mg.service
[Unit]
Description=MG
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
ExecStart=/usr/local/bin/mg

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl start mg.service
echo "Install complete!"