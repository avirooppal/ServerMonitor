git clone https://github.com/avirooppal/ServerMonitor.git
cd ServerMonitor

# Start the Dashboard
docker-compose up -d --build
```

The dashboard will be available at `http://YOUR_VPS_IP:8082`.

### 2. Get Your Master Key
On the first run, a secure Master API Key is generated. Retrieve it:

```bash
docker exec server-moni-dashboard-1 cat data/api_key.txt
```

### 3. Login & Add Agents
1. Open the Dashboard.
2. Go to **Settings**.
3. Enter your **Master API Key** to unlock admin features.
4. Copy the **One-Click Installer** command shown in the dashboard.

### 4. Install Agents
Run the copied command on **any** Linux server you want to monitor:

```bash
curl -sL https://raw.githubusercontent.com/avirooppal/ServerMonitor/main/web/public/get-key.sh | bash
```

This script will:
1. Check for Docker (and install if missing).
2. Pull and run the Agent container.
3. Output the **API Key** and **URL** for that agent.

### 5. Add to Dashboard
Copy the **Name**, **URL**, and **API Key** from the script output back into your Dashboard's "Add System" form.

---

## 🛠️ Local Development

To run the entire stack (Dashboard + Agent) locally on your machine:

1. **Start the Stack**:
   ```bash
   docker-compose up -d --build
   ```

2. **Access Dashboard**:
   - Open `http://localhost:8082`

3. **Connect Local Agent**:
   - The local agent is running on port `8081`.
   - In the Dashboard, add a system with:
     - **URL**: `http://agent:8080` (Internal Docker Network)
     - **API Key**: Retrieve from `docker exec server-moni-agent-1 cat data/api_key.txt`

---

## 🔧 Configuration

### Environment Variables

| Service | Variable | Description | Default |
|---------|----------|-------------|---------|
| **Dashboard** | `PORT` | Port to serve the UI/API | `8080` |
| **Agent** | `API_PORT` | Port for Agent API | `8080` |
| **Agent** | `COLLECTION_INTERVAL_SECONDS` | Metrics update frequency | `5` |

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

MIT License. Free for personal and commercial use.
