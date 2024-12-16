WEBHOOK_URL="https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN"


MESSAGE="\$debug hello"

curl -H "Content-Type: application/json" \
     -d "{\"content\": \"$MESSAGE\"}" \
     "$WEBHOOK_URL"