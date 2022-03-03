#!/usr/bin/env bash

systemctl stop mg.service
rm -rf /etc/systemd/system/mg.service
systemctl daemon-reload
rm -rf /usr/local/bin/mg