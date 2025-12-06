# PowerShell Script to Install Agent
param (
    [string]$ServerUrl,
    [string]$Token
)

[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

if ([string]::IsNullOrEmpty($ServerUrl) -or [string]::IsNullOrEmpty($Token)) {
    Write-Host "Usage: .\install_agent.ps1 -ServerUrl http://your-server.com -Token YOUR_API_KEY"
    exit 1
}

$AgentUrl = "https://github.com/avirooppal/ServerMonitor/raw/main/downloads/agent-windows-amd64.exe"
$InstallDir = "C:\Program Files\ServerMonitor"
$AgentPath = "$InstallDir\agent.exe"

# Stop existing service if running
$ServiceName = "ServerMoniAgent"
if (Get-Service $ServiceName -ErrorAction SilentlyContinue) {
    Write-Host "Stopping existing service..."
    Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
    Start-Sleep -Seconds 2
    
    # Uninstall old service to ensure clean state
    if (Test-Path $AgentPath) {
        & $AgentPath --service uninstall
    }
}

# Kill process if still running (to unlock file)
Get-Process -Name "agent" -ErrorAction SilentlyContinue | Stop-Process -Force

# Create Directory
New-Item -ItemType Directory -Force -Path $InstallDir

# Download Agent
Write-Host "Downloading Agent from $AgentUrl..."
Invoke-WebRequest -Uri $AgentUrl -OutFile $AgentPath

# Install Service
Write-Host "Installing Service..."
& $AgentPath -server $ServerUrl -token $Token --service install

# Start Service
Write-Host "Starting Service..."
& $AgentPath --service start

Write-Host "Agent installed and started!"
