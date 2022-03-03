#!/usr/bin/env bash

wget -q https://github.com/welcome-ad-infernum/mg/releases/latest/download/mg_linux_amd64 -O /usr/local/bin/mg
chmod +x /usr/local/bin/mg
cp -r mg.service /etc/systemd/system/mg.service
systemctl daemon-reload
systemctl start mg.service
systemctl status mg.service