#!/usr/bin/bash

HOOK_FILE=".git/hooks/commit-msg"

cat > "$HOOK_FILE" << 'EOF'
#!/usr/bin/bash

$(git rev-parse --show-toplevel)/scripts/commit-msg.sh "$1"
EOF

chmod +x $HOOK_FILE

echo "'commit-msg' hook set up successfully"
