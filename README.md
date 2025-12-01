# Server Monitor

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/backend-Go-00ADD8.svg)
![React](https://img.shields.io/badge/frontend-React-61DAFB.svg)
![Docker](https://img.shields.io/badge/deployment-Docker-2496ED.svg)

**Server Monitor** is a lightweight, self-hosted infrastructure monitoring solution. It provides a centralized dashboard to track the health and performance of your Linux servers in real-time.

![Dashboard Preview](https://github.com/user-attachments/assets/2779b9a6-8ab2-407e-9a9a-f44ffc6afcbe)

---

## ✨ Features

- **Real-Time Metrics**: Monitor CPU, RAM, Disk, and Network usage with live updates.
- **Multi-Server Support**: Manage unlimited servers from a single dashboard.
- **Lightweight Agent**: Low-overhead Go-based agent (~10MB) that runs on any Linux distribution.
- **Secure Architecture**:
    - **Master API Key**: Protects the dashboard from unauthorized access.
    - **Agent Tokens**: Unique, auto-generated keys ensure secure communication between agents and the dashboard.
    - **Proxy Mode**: The dashboard acts as a proxy, keeping your agents secure behind firewalls without exposing extra ports.
- **One-Click Installation**: Deploy agents instantly using a simple `curl | bash` script.
- **Docker Native**: Built for containerized environments with easy `docker-compose` deployment.

---

## 🚀 Quick Start

Get your monitoring dashboard up and running in minutes.

### Prerequisites
- A Linux server (VPS) with **Docker** and **Docker Compose** installed.

### 1. Deploy the Dashboard
Clone the repository and start the services:

```bash
git clone https://github.com/avirooppal/ServerMonitor.git
cd ServerMonitor
docker-compose up -d --build
```

The dashboard will be available at `http://YOUR_SERVER_IP:8082`.

### 2. Initial Setup
On the first run, a secure **Master API Key** is generated. You need this key to log in and manage the system.

Retrieve the key:
```bash
docker exec server-moni-dashboard-1 cat data/api_key.txt
```

1. Open the Dashboard in your browser.
2. Go to the **Settings** tab.
3. Enter your **Master API Key**.

### 3. Add Agents (Monitor Servers)
To monitor a server (including the one hosting the dashboard), you need to install the Agent.

1. In the Dashboard **Settings**, copy the **One-Click Installer** command.
2. SSH into the target server.
3. Run the command:
   ```bash
   curl -sL https://raw.githubusercontent.com/avirooppal/ServerMonitor/main/web/public/get-key.sh | bash
   ```
4. The script will output the **Agent Name**, **URL**, and **API Key**.
5. Go back to the Dashboard, click **Add System**, and enter these details.

---

## 🏗️ Architecture

The system consists of a central Dashboard (Server) and multiple Agents (Clients).

```mermaid
graph TD
    User[User / Browser] -->|HTTPS| Dashboard[Central Dashboard]
    Dashboard -->|Proxy Request| Agent1[Agent: DB Server]
    Dashboard -->|Proxy Request| Agent2[Agent: Web Server]
    Agent1 -->|Collects| Docker1[Docker Engine]
    Agent2 -->|Collects| Docker2[Docker Engine]
```

- **Dashboard**: Stores configuration (SQLite), serves the UI, and proxies requests to agents.
- **Agent**: Collects system metrics and exposes them via a secured API.

---

## 🛠️ Local Development

To run the project locally for development:

1. **Start the Stack**:
   ```bash
   docker-compose up -d --build
   ```

2. **Access Dashboard**:
   Open `http://localhost:8082`.

3. **Connect Local Agent**:
   The local setup includes an agent running on port `8081`.
   - **URL**: `http://agent:8080` (Internal Docker Network URL)
   - **API Key**: `docker exec server-moni-agent-1 cat data/api_key.txt`

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

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.
