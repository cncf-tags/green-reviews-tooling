#!/bin/bash
set -e

IPTABLES_BACKUP="/tmp/iptables-rules.backup"
IPTABLES_MODIFIED="/tmp/iptables-rules.modified"

echo "removing iptables rules that block cluster networking"
sudo apt-get update -qq && sudo apt-get install -y iptables-persistent -qq
sudo iptables-save > "$IPTABLES_BACKUP"
grep -v "DROP" "$IPTABLES_BACKUP" | grep -v "REJECT" > "$IPTABLES_MODIFIED"
sudo iptables-restore < "$IPTABLES_MODIFIED"
echo "modified iptables rules applied successfully"
sudo iptables -L
sudo netfilter-persistent save
rm -f "$IPTABLES_BACKUP" "$IPTABLES_MODIFIED"
echo "iptables rules removed successfully"
