# Server Monitor

Server Monitor is a lightweight, self-hosted infrastructure monitoring solution. It provides a centralized dashboard to track the health and performance of your Linux and Windows servers in real-time.

## Features

- **Real-Time Metrics**: Monitor CPU, RAM, Disk, and Network usage with live updates.
- **Multi-Server Support**: Manage unlimited servers from a single dashboard.
- **Cross-Platform Agent**: Lightweight Go-based agent that runs on Linux and Windows.
- **Secure Architecture**:
    - **User Authentication**: Secure login and registration system.
    - **Agent Tokens**: Unique, auto-generated keys ensure secure communication between agents and the dashboard.
    - **Push-Based Architecture**: Agents push data to the central server, eliminating the need for complex firewall configurations.
- **One-Click Installation**: Deploy agents instantly using simple scripts.

## Architecture

The system consists of three main components:

1.  **Backend (Server)**: Written in Go, using Gin framework and SQLite database. It handles API requests, authentication, and metric ingestion.
2.  **Frontend (Dashboard)**: Built with React and Vite, hosted on Vercel. It provides the user interface for monitoring systems.
3.  **Agent**: A lightweight Go binary running on monitored servers. It collects system metrics (via `gopsutil`) and pushes them to the Backend.

## Prerequisites

- **Go 1.22+** (for Backend and Agent development)
- **Node.js 18+** (for Frontend development)
- **GCC** (for CGO/SQLite support)

## Installation and Setup

### 1. Backend Setup

The backend serves the API and the agent installation scripts.

1.  Clone the repository:
    ```bash
    git clone https://github.com/avirooppal/ServerMonitor.git
    cd ServerMonitor
    ```

2.  Install dependencies:
    ```bash
    go mod download
    ```

3.  Run the server:
    ```bash
    go run cmd/server/main.go
    ```
    The server will start on port `8080` (default).

### 2. Frontend Setup

The frontend is a React application located in the `web` directory.

1.  Navigate to the web directory:
    ```bash
    cd web
    ```

2.  Install dependencies:
    ```bash
    npm install
    ```

3.  Start the development server:
    ```bash
    npm run dev
    ```
    The dashboard will be available at `http://localhost:5173`.

### 3. Agent Installation

To monitor a server, you need to install the agent on it.

#### Linux

Run the following command on your Linux server:

```bash
curl -L https://raw.githubusercontent.com/avirooppal/ServerMonitor/main/scripts/install_agent_linux.sh | sudo bash -s -- --server=http://YOUR_BACKEND_IP:8080 --token=YOUR_API_KEY
```

#### Windows

Run the following command in PowerShell (Administrator):

```powershell
iwr https://raw.githubusercontent.com/avirooppal/ServerMonitor/main/scripts/install_agent_windows.ps1 -OutFile install.ps1; .\install.ps1 -ServerUrl http://YOUR_BACKEND_IP:8080 -Token YOUR_API_KEY
```

*Note: Replace `YOUR_BACKEND_IP` and `YOUR_API_KEY` with your actual backend URL and the system API key generated from the dashboard.*

## API Documentation

The backend exposes a RESTful API.

### Authentication

#### Register
Create a new user account.

- **URL**: `/api/v1/auth/register`
- **Method**: `POST`
- **Body**:
    ```json
    {
        "email": "user@example.com",
        "password": "securepassword"
    }
    ```

#### Login
Authenticate and receive a JWT token.

- **URL**: `/api/v1/auth/login`
- **Method**: `POST`
- **Body**:
    ```json
    {
        "email": "user@example.com",
        "password": "securepassword"
    }
    ```
- **Response**:
    ```json
    {
        "token": "eyJhbGciOiJIUzI1Ni..."
    }
    ```

### Systems

#### Get Systems
List all monitored systems.

- **URL**: `/api/v1/systems`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <TOKEN>`

#### Add System
Register a new system to monitor.

- **URL**: `/api/v1/systems`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <TOKEN>`
- **Body**:
    ```json
    {
        "name": "Production DB",
        "url": "push",
        "api_key": "generated-random-key"
    }
    ```

### Metrics

#### Get Metrics
Retrieve the latest metrics for a specific system.

- **URL**: `/api/v1/metrics?system_id=<ID>`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <TOKEN>`

#### Ingest Metrics (Agent)
Push metrics from the agent to the server.

- **URL**: `/api/v1/ingest`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <SYSTEM_API_KEY>`
- **Body**:
    ```json
    {
        "cpu_usage": 45.5,
        "memory_total": 16000000000,
        "memory_used": 8000000000,
        "disk_usage": 60.2
    }
    ```

### Health Check

Check if the server is running.

- **URL**: `/health`
- **Method**: `GET`
- **Response**:
    ```json
    {
        "status": "ok",
        "time": "2023-10-27T10:00:00Z"
    }
    ```

## Development

### Directory Structure

- `cmd/server`: Entry point for the backend server.
- `cmd/agent`: Entry point for the monitoring agent.
- `internal/api`: API handlers and router configuration.
- `internal/auth`: Authentication logic (JWT, bcrypt).
- `internal/db`: Database interaction (SQLite).
- `internal/metrics`: Metric collection and storage logic.
- `web`: React frontend application.
- `scripts`: Installation and utility scripts.

### Building

To build the binaries:

```bash
# Build Server
go build -o server-moni cmd/server/main.go

# Build Agent (Linux)
GOOS=linux GOARCH=amd64 go build -o server-moni-agent cmd/agent/main.go

# Build Agent (Windows)
GOOS=windows GOARCH=amd64 go build -o server-moni-agent.exe cmd/agent/main.go
```

## License

Distributed under the MIT License.
