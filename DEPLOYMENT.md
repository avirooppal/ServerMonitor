# Deployment Guide

This guide explains how to deploy the Server Monitor in a "SaaS" configuration:
1.  **Backend**: Hosted on your VPS (Oracle Cloud, DigitalOcean, etc.).
2.  **Frontend**: Hosted on Vercel (Free).
3.  **Agents**: Installed on client servers using a single command.

---

## Part 1: Deploy Backend to VPS

### 1. Prepare the VPS
Login to your VPS (SSH) and run the following commands to install the backend.

**Option A: Build on VPS (Recommended)**
Requires `git`, `go`, and `npm` installed on the VPS.
```bash
# 1. Clone the repo
git clone https://github.com/aviroop/server-moni
cd server-moni

# 2. Run the install script
sudo ./scripts/install_server_linux.sh
```

**Option B: Upload Binary**
If you don't want to install Go on the VPS:
1.  **On your PC**: Run `set GOOS=linux` then `go build -o server-moni ./cmd/server`.
2.  **Upload**: Copy `server-moni` to your VPS.
3.  **Run**: `./server-moni` (or create a systemd service manually).

### 2. Verify Backend
Your backend should now be running on `http://YOUR_VPS_IP:8080`.
*   Test it: `curl http://YOUR_VPS_IP:8080/api/v1/ping`

### 3. Enable HTTPS (Required for Vercel)
Since Vercel uses HTTPS, your backend MUST also use HTTPS.
The easiest way (without buying a domain) is using **Cloudflare Tunnel**.

**Run this on your VPS:**
```bash
# 1. Download cloudflared
curl -L --output cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared.deb

# 2. Start a Quick Tunnel
cloudflared tunnel --url http://localhost:8080
```
*   Copy the URL that looks like `https://random-name.trycloudflare.com`.
*   **Keep this running** (use `tmux` or a service to keep it alive).

---

## Part 2: Deploy Frontend to Vercel

1.  **Push to GitHub**: Ensure your latest code is on GitHub.
2.  **Import to Vercel**:
    *   Go to Vercel Dashboard -> Add New -> Project.
    *   Select your `server-moni` repository.
3.  **Configure Build**:
    *   **Root Directory**: `web` (Important! The frontend is in the `web` folder).
    *   **Build Command**: `npm run build`
    *   **Output Directory**: `dist`
4.  **Environment Variables**:
    *   Add a new variable: `VITE_API_URL`
    *   Value: `http://YOUR_VPS_IP:8080/api/v1` (Note: If Vercel forces HTTPS, you might need to set up SSL on your VPS using Nginx/Certbot, or use a Cloudflare tunnel).
5.  **Deploy**: Click Deploy.

---

## Part 3: The "One Command" Experience (For Users)

Now that your backend is running, you can give this command to your users (or run it on your other servers).

**Command:**
```bash
curl -L http://YOUR_VPS_IP:8080/downloads/agent-linux-amd64 -o agent && \
chmod +x agent && \
sudo ./agent
```

# Deployment Guide

This guide explains how to deploy the Server Monitor in a "SaaS" configuration:
1.  **Backend**: Hosted on your VPS (Oracle Cloud, DigitalOcean, etc.).
2.  **Frontend**: Hosted on Vercel (Free).
3.  **Agents**: Installed on client servers using a single command.

---

## Part 1: Deploy Backend to VPS

### 1. Prepare the VPS
Login to your VPS (SSH) and run the following commands to install the backend.

**Option A: Build on VPS (Recommended)**
Requires `git`, `go`, and `npm` installed on the VPS.
```bash
# 1. Clone the repo
git clone https://github.com/aviroop/server-moni
cd server-moni

# 2. Run the install script
sudo ./scripts/install_server_linux.sh
```

**Option B: Upload Binary**
If you don't want to install Go on the VPS:
1.  **On your PC**: Run `make build-server-linux`.
2.  **Upload**: Copy `bin/server-linux-amd64` to your VPS.
3.  **Run**: `./server-linux-amd64` (or create a systemd service manually).

### 2. Verify Backend
Your backend should now be running on `http://YOUR_VPS_IP:8080`.
*   Test it: `curl http://YOUR_VPS_IP:8080/api/v1/ping`

---

## Part 2: Deploy Frontend to Vercel

1.  **Push to GitHub**: Ensure your latest code is on GitHub.
2.  **Import to Vercel**:
    *   Go to Vercel Dashboard -> Add New -> Project.
    *   Select your `server-moni` repository.
3.  **Configure Build**:
    *   **Root Directory**: `web` (Important! The frontend is in the `web` folder).
    *   **Build Command**: `npm run build`
    *   **Output Directory**: `dist`
4.  **Environment Variables**:
    *   Add a new variable: `VITE_API_URL`
    *   Value: `http://YOUR_VPS_IP:8080/api/v1` (Note: If Vercel forces HTTPS, you might need to set up SSL on your VPS using Nginx/Certbot, or use a Cloudflare tunnel).
5.  **Deploy**: Click Deploy.

---

## Part 3: The "One Command" Experience (For Users)

Now that your backend is running, you can give this command to your users (or run it on your other servers).

**Command:**
```bash
curl -L http://YOUR_VPS_IP:8080/downloads/agent-linux-amd64 -o agent && \
chmod +x agent && \
sudo ./agent
```

*Wait!* To make it truly "Zero Touch" (auto-connect), you should use the install script.

**Better Command:**
1.  Host the `scripts/install_agent_linux.sh` on your server (or GitHub Gist).
2.  User runs:
    ```bash
    curl -sL http://YOUR_VPS_IP:8080/downloads/install.sh | sudo bash -s -- --server=http://YOUR_VPS_IP:8080 --token=USER_API_KEY
    ```

### Windows Installation
For Windows servers, run this in PowerShell (Administrator):

```powershell
$env:SERVER_URL="http://YOUR_VPS_IP:8080"; $env:API_KEY="USER_API_KEY"; 
iwr http://YOUR_VPS_IP:8080/downloads/install.ps1 -UseBasicParsing | iex
```
*(Note: You need to copy `scripts/install_agent_windows.ps1` to `downloads/install.ps1` on your VPS first)*

### How to Enable the "Better Command"
1.  Copy `scripts/install_agent_linux.sh` to `downloads/install.sh` on your VPS.
    ```bash
    cp scripts/install_agent_linux.sh /var/lib/server-moni/downloads/install.sh
    ```
2.  Copy `scripts/install_agent_windows.ps1` to `downloads/install.ps1` on your VPS.
    ```bash
    cp scripts/install_agent_windows.ps1 /var/lib/server-moni/downloads/install.ps1
    ```
3.  Now the commands above work perfectly!
