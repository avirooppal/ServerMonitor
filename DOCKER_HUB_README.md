# Server Monitor Agent 🕵️‍♂️

A lightweight, secure metrics collector for the **Server Monitor** SaaS platform.

## 🚀 Quick Start

The easiest way to install the agent is using our one-click installer script.

### 1. Get Your Dashboard
Ensure you have a running instance of the **Server Monitor Dashboard**.

### 2. Install Agent
Run this command on your Linux server:

```bash
curl -sL https://raw.githubusercontent.com/avirooppal/ServerMonitor/main/web/public/get-key.sh | bash
```

This script will:
- ✅ Check for Docker
- ⬇️ Pull this image (`avirooppal/server-moni-agent`)
- 🏃‍♂️ Run the agent container
- 🔑 Output your **API Key** and **Connection URL**

### 3. Connect
Copy the API Key and URL into your Dashboard's "Add System" settings.

---

## 🐳 Manual Run

If you prefer to run the container manually without the helper script:

```bash
docker run -d \
  --name server-moni-agent \
  --network host \
  --restart unless-stopped \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v server-moni-data:/app/data \
  -e API_PORT=8080 \
  avirooppal/server-moni-agent:latest
```

*Note: `--network host` is recommended to allow the agent to accurately report host network metrics.*

## 🔧 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `API_PORT` | Port for the Agent API | `8080` |
| `COLLECTION_INTERVAL_SECONDS` | How often to collect metrics | `5` |
| `DB_NAME` | Internal SQLite DB name | `metrics.db` |
| `RETENTION_HOURS` | Local data retention | `24` |

## 📦 Volume Mounts

- `/var/run/docker.sock`: **Required**. Allows the agent to collect Docker container stats.
- `/app/data`: **Recommended**. Persists the generated API Key and local metric history.
