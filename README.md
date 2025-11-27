# Server Moni

A beautiful, real-time server monitoring dashboard inspired by the premium aesthetic of FiveNines.io. Built with Go and React.

![Dashboard Preview](https://via.placeholder.com/800x450.png?text=Server+Moni+Dashboard)

## Features

- **Real-Time Metrics**: Live updates for CPU, RAM, Swap, and Disk usage.
- **Network Monitoring**: Real-time upload and download rates with historical charts.
- **Process Manager**: View active processes with sorting, filtering, and resource usage (CPU/Mem).
- **Docker Integration**: Monitor running containers, their status, and resource consumption.
- **Modern UI**: Deep dark theme with neon accents, glassmorphism effects, and smooth animations.
- **Responsive**: Fully functional on desktop and mobile devices.

## Tech Stack

- **Backend**: Go (Golang), Gin Framework, gopsutil, Docker SDK
- **Frontend**: React, TypeScript, Vite, Tailwind CSS, Recharts, Lucide Icons
- **Database**: SQLite (Embedded)
- **Deployment**: Docker & Docker Compose

## Quick Start

The easiest way to run Server Moni is using Docker Compose.

1.  **Clone the repository**
    ```bash
    git clone https://github.com/yourusername/server-moni.git
    cd server-moni
    ```

2.  **Start the application**
    ```bash
    docker-compose up --build
    ```

3.  **Access the Dashboard**
    Open your browser and navigate to `http://localhost:8080`.

4.  **Authenticate**
    - On the first run, the server generates a secure **API Key**.
    - Check the container logs (`docker-compose logs`) or look at the file `data/api_key.txt` to find it.
    - Enter this key in the web interface to connect.

## Deployment

### Full Stack (Single Server)
Simply run the Docker container on your VPS. The frontend is served statically by the Go backend.
```bash
docker-compose up -d
```

### Separate Frontend (Optional)
You can host the frontend on a CDN (like Vercel, Netlify, or Cloudflare Pages) and the backend on your VPS.

1.  **Backend**: Run the container on your VPS. Ensure port `8080` (or your chosen port) is accessible.
2.  **Frontend**: Deploy the `web` folder to your host.
3.  **Configuration**: Set the `VITE_API_URL` environment variable on your frontend host to point to your backend (e.g., `http://your-vps-ip:8080/api/v1`).

## Development

### Backend
```bash
go run cmd/server/main.go
```

### Frontend
```bash
cd web
npm install
npm run dev
```

## License

MIT
