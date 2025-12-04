# Usage: iwr http://server/install.ps1 -UseBasicParsing | iex
# Arguments are passed via environment variables if running manually, 
# or we can parse them if the user downloads and runs.
# For "One Command" web install, we usually recommend:
# $env:SERVER_URL="http://..."; $env:API_KEY="..."; iwr ... | iex

param (
    [string]$ServerUrl = $env:SERVER_URL,
    [string]$ApiKey = $env:API_KEY
)

if ([string]::IsNullOrEmpty($ServerUrl) -or [string]::IsNullOrEmpty($ApiKey)) {
    Write-Host "Usage: Set SERVER_URL and API_KEY env vars or pass arguments." -ForegroundColor Red
    Write-Host "Example: `$env:SERVER_URL='http://your-server'; `$env:API_KEY='key'; iwr ... | iex"
    exit 1
}

Write-Host "Installing Server Monitor Agent..." -ForegroundColor Cyan

# 1. Create Directory
$InstallDir = "C:\Program Files\ServerMonitor"
if (!(Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

# 2. Download Binary
$AgentUrl = "$ServerUrl/downloads/agent-windows-amd64.exe"
$AgentPath = "$InstallDir\agent.exe"

Write-Host "Downloading Agent from $AgentUrl..."
try {
    Invoke-WebRequest -Uri $AgentUrl -OutFile $AgentPath -UseBasicParsing
} catch {
    Write-Host "Failed to download agent: $_" -ForegroundColor Red
    exit 1
}

# 3. Install as Service using NSSM or sc.exe
# Go binaries don't automatically register as Windows services unless using a library like kardianos/service.
# Since we didn't add that library yet, we will use a simple workaround: Scheduled Task or just run it.
# Ideally, we should add 'kardianos/service' to the agent code.
# For now, let's use 'sc.exe' but the binary needs to handle service control signals.
# If the binary is just a loop, 'sc' might timeout waiting for it to "start".
# A better approach for a raw binary is using a wrapper or just running it in background for now.

# WAIT: To make this "Robust", I should really add service support to the Go code.
# But for now, let's use a Scheduled Task which is robust enough for "Start on Boot".

$TaskName = "ServerMonitorAgent"
Write-Host "Creating Scheduled Task '$TaskName'..."

$Action = New-ScheduledTaskAction -Execute $AgentPath
$Trigger = New-ScheduledTaskTrigger -AtStartup
$Settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -Hidden

# We need to pass env vars. Scheduled tasks are tricky with env vars.
# We will write a wrapper script.
$WrapperPath = "$InstallDir\start_agent.ps1"
$WrapperContent = @"
`$env:SERVER_URL = '$ServerUrl'
`$env:API_KEY = '$ApiKey'
Start-Process -FilePath '$AgentPath' -WindowStyle Hidden -Wait
"@
Set-Content -Path $WrapperPath -Value $WrapperContent

$ActionWrapper = New-ScheduledTaskAction -Execute "powershell.exe" -Argument "-ExecutionPolicy Bypass -File `"$WrapperPath`""

Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false -ErrorAction SilentlyContinue
Register-ScheduledTask -Action $ActionWrapper -Trigger $Trigger -Settings $Settings -TaskName $TaskName -User "SYSTEM" -RunLevel Highest

Write-Host "Starting Agent..."
Start-ScheduledTask -TaskName $TaskName

Write-Host "Agent Installed and Started!" -ForegroundColor Green
