#!/bin/bash

set -euo pipefail

# --- HARDWARE PERMISSIONS ---
# If this service requires hardware access, you likely need a udev rule
# to assign ownership to the 'openmon' user.
# Example: /etc/udev/rules.d/99-openmon.rules
# SUBSYSTEM=="usb", ATTRS{idVendor}=="XXXX", OWNER="openmon"
# ----------------------------

echo "Linking sysusers config..."

mkdir -p /etc/sysusers.d

if [ ! -f /etc/sysusers.d/openmon.conf ]; then
    ln -s "/opt/openmon/conf/openmon.conf" /etc/sysusers.d/openmon.conf
fi

echo "Creating user..."
systemd-sysusers

echo "Linking unit..."
if [ -f /etc/systemd/system/openmon.service ]; then
    rm /etc/systemd/system/openmon.service
fi

systemctl link "/opt/openmon/conf/openmon.service"

echo "Reloading daemon..."
systemctl daemon-reload
systemctl enable openmon

echo "Fixing initial permissions..."
chown -R openmon:openmon "/opt/openmon"

find "/opt/openmon" -type d -exec chmod 755 {} +
find "/opt/openmon" -type f -exec chmod 644 {} +

chmod +x "/opt/openmon/openmon"

echo "Setup complete, starting service..."

service openmon restart

echo "Done."
