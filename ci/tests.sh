
sleep 10

WEBHOOK_URL="https://discord.com/api/webhooks/$WEBHOOKTOKEN"


MESSAGE="\$debug hello"

curl -H "Content-Type: application/json" \
     -d "{\"content\": \"$MESSAGE\"}" \
     "$WEBHOOK_URL"

sleep 5