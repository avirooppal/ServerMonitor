#!/bin/bash

# Server Monitor - One-Click Installer (FiveNines Style)
# Usage: curl ... | bash -s -- <TOKEN> [DOMAIN]

TOKEN=$1
DOMAIN=$2

if [ -z "$TOKEN" ]; then
    echo "Error: No token provided."
    echo "Usage: curl ... | bash -s -- <TOKEN> [DOMAIN]"
    echo "Example (HTTP):  curl ... | bash -s -- my-token"
    echo "Example (HTTPS): curl ... | bash -s -- my-token 1.2.3.4.nip.io"
    exit 1
fi

echo "🚀 Installing Server Monitor Agent..."

# 1. Install Docker if missing
if ! command -v docker &> /dev/null; then
    echo "Docker not found. Installing..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
fi

# 2. Cleanup old containers
echo "🧹 Cleaning up old containers..."
docker stop server-moni-agent server-moni-caddy 2>/dev/null || true
docker rm server-moni-agent server-moni-caddy 2>/dev/null || true
docker network rm server-moni-net 2>/dev/null || true

# 3. Pull Agent Image
echo "📦 Pulling agent image..."
docker pull avirooppal/server-moni-agent:latest

# 4. Determine Mode (HTTP vs HTTPS)
if [ -z "$DOMAIN" ]; then
    # --- HTTP MODE ---
    echo "🌐 Mode: HTTP (Port 8080)"
    echo "⚠️  Warning: Vercel (HTTPS) cannot connect to HTTP. Use a domain for SSL if needed."
    
    docker run -d \
      --name server-moni-agent \
      --restart unless-stopped \
      --network host \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -v /:/host:ro \
      -e AGENT_SECRET="$TOKEN" \
      -e API_PORT=8080 \
      avirooppal/server-moni-agent:latest

    IP=$(hostname -I | awk '{print $1}')
    URL="http://$IP:8080"

else
    # --- HTTPS MODE (Caddy) ---
    echo "🔒 Mode: HTTPS (Domain: $DOMAIN)"
    
    # Create Network
    docker network create server-moni-net

    # Run Agent (Internal Port 8080)
    docker run -d \
      --name server-moni-agent \
      --restart unless-stopped \
      --network server-moni-net \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -v /:/host:ro \
      -e AGENT_SECRET="$TOKEN" \
      -e API_PORT=8080 \
      avirooppal/server-moni-agent:latest

    # Create Caddyfile
    echo "$DOMAIN {
        reverse_proxy server-moni-agent:8080
    }" > Caddyfile

    # Run Caddy
    docker run -d \
      --name server-moni-caddy \
      --restart unless-stopped \
      --network server-moni-net \
      -p 80:80 -p 443:443 \
      -v $(pwd)/Caddyfile:/etc/caddy/Caddyfile \
      -v caddy_data:/data \
      -v caddy_config:/config \
      caddy:2

    URL="https://$DOMAIN"
fi

echo ""
echo "✅ Agent Installed Successfully!"
echo "------------------------------------------------"
echo "URL:   $URL"
echo "Token: $TOKEN"
echo "------------------------------------------------"
echo "Go back to your dashboard and enter these details to connect."
