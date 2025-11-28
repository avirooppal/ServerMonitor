# Docker Metrics Monitor

## 🎯 Business Model

### Open Source Agent (FREE)
The **metrics collection agent** is completely open source and free to use:
- ✅ Collects Docker container metrics
- ✅ Collects system metrics (CPU, RAM, Disk, Network)
- ✅ Stores data locally in SQLite
- ✅ Provides REST API for data access
- ✅ No limitations, no restrictions

### Cloud Visualization (PREMIUM)
The **dashboard and visualization** is a cloud-hosted SaaS service:
- 🌟 **Free Tier**: Monitor 1 system
- 💎 **Pro Tier**: Monitor up to 10 systems ($19/month)
- 🚀 **Enterprise Tier**: Unlimited systems ($99/month)

---

## 🆓 Free Tier Features

**Perfect for individuals and small projects:**
- ✅ Monitor 1 server/system
- ✅ Real-time metrics dashboard
- ✅ 24 hours data retention
- ✅ Container logs viewer
- ✅ Basic charts and graphs
- ✅ Email support

**Sign up at:** [https://metrics.yourcloud.com/signup](https://metrics.yourcloud.com/signup)

---

## 💎 Pro Tier ($19/month)

**For growing teams and multiple servers:**
- ✅ Monitor up to 10 systems
- ✅ 30 days data retention
- ✅ Advanced analytics
- ✅ Custom alerts (Email, Slack, Discord)
- ✅ API access for integrations
- ✅ Priority support
- ✅ Team collaboration (up to 5 users)

---

## 🚀 Enterprise Tier ($99/month)

**For large organizations:**
- ✅ Unlimited systems
- ✅ 90 days data retention
- ✅ Advanced AI-powered insights
- ✅ Custom integrations
- ✅ SSO/SAML authentication
- ✅ Dedicated support
- ✅ SLA guarantee (99.9% uptime)
- ✅ Unlimited team members
- ✅ White-label option available

---

## 🚀 Quick Start

### Step 1: Install the Agent (Open Source)

```bash
# Clone the repository
git clone https://github.com/yourusername/docker-metrics-agent
cd docker-metrics-agent

# Start the agent with Docker Compose
docker-compose -f docker-compose.yml up -d
```

The agent will start collecting metrics and expose them on `http://localhost:8080/api`

### Step 2: Sign Up for Cloud Dashboard

1. Visit [https://metrics.yourcloud.com/signup](https://metrics.yourcloud.com/signup)
2. Create a free account
3. Get your API key

### Step 3: Connect Your Agent

```bash
# Set your API key
export API_KEY="your-api-key-from-dashboard"

# Restart the agent with API key
docker-compose down
docker-compose -f docker-compose.yml up -d
```

### Step 4: Configure Cloud Access

1. Go to [https://metrics.yourcloud.com/dashboard](https://metrics.yourcloud.com/dashboard)
2. Click "Add System"
3. Enter your system details:
   - **Name**: My Production Server
   - **API Endpoint**: `https://your-server.com:8080/api` (or use Cloudflare Tunnel)
   - **API Key**: (auto-filled from your account)
4. Click "Connect"

That's it! Your metrics will now appear in the cloud dashboard.

---

## 🔐 Security & Privacy

### Your Data Stays Local
- ✅ All metrics are stored **locally** on your infrastructure
- ✅ The cloud dashboard **never stores** your raw metrics
- ✅ Only aggregated statistics are sent to the cloud (optional)
- ✅ You maintain full control of your data

### How It Works
1. Agent collects metrics locally
2. Cloud dashboard connects to your agent via API
3. Data is fetched **on-demand** when you view the dashboard
4. No persistent storage of your metrics in the cloud

### API Key Security
- 🔒 API keys are encrypted in transit (HTTPS)
- 🔒 Keys can be rotated anytime
- 🔒 IP whitelisting available (Pro/Enterprise)
- 🔒 Rate limiting prevents abuse

---

## 🌐 Exposing Your Agent to the Cloud

For the cloud dashboard to access your agent, you need to expose it to the internet. Here are secure options:

### Option 1: Cloudflare Tunnel (Recommended - FREE)

```bash
# Install cloudflared
curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o cloudflared
chmod +x cloudflared

# Create tunnel
./cloudflared tunnel --url http://localhost:8080
```

Copy the generated URL (e.g., `https://abc123.trycloudflare.com`) and use it in the cloud dashboard.

### Option 2: Reverse Proxy (Nginx + Let's Encrypt)

```nginx
server {
    listen 443 ssl;
    server_name metrics.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/metrics.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/metrics.yourdomain.com/privkey.pem;

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header X-API-Key $http_x_api_key;
    }
}
```

### Option 3: Tailscale (Private Network)

```bash
# Install Tailscale
curl -fsSL https://tailscale.com/install.sh | sh

# Connect to your Tailscale network
tailscale up

# Use your Tailscale IP in the cloud dashboard
# e.g., http://100.64.1.2:8080/api
```

---

## 📊 Feature Comparison

| Feature | Free | Pro | Enterprise |
|---------|------|-----|------------|
| **Systems** | 1 | 10 | Unlimited |
| **Data Retention** | 24h | 30d | 90d |
| **Users** | 1 | 5 | Unlimited |
| **Real-time Metrics** | ✅ | ✅ | ✅ |
| **Container Logs** | ✅ | ✅ | ✅ |
| **Historical Charts** | ✅ | ✅ | ✅ |
| **Email Alerts** | ❌ | ✅ | ✅ |
| **Slack/Discord Alerts** | ❌ | ✅ | ✅ |
| **API Access** | ❌ | ✅ | ✅ |
| **Custom Dashboards** | ❌ | ✅ | ✅ |
| **AI Insights** | ❌ | ❌ | ✅ |
| **SSO/SAML** | ❌ | ❌ | ✅ |
| **White Label** | ❌ | ❌ | ✅ |
| **SLA** | ❌ | ❌ | 99.9% |
| **Support** | Email | Priority | Dedicated |
| **Price** | FREE | $19/mo | $99/mo |

---

## 🔧 Agent Configuration

The agent supports these environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `API_PORT` | Port for API server | `8080` | No |
| `DB_NAME` | SQLite database path | `metrics.db` | No |
| `COLLECTION_INTERVAL_SECONDS` | Metrics collection interval | `5` | No |
| `RETENTION_HOURS` | Local data retention | `24` | No |
| `API_KEY` | Your cloud dashboard API key | - | **Yes** |
| `ALLOWED_ORIGINS` | CORS origins (auto-set) | - | No |

### Example Configuration

```yaml
# docker-compose.yml
services:
  metrics-agent:
    image: yourorg/docker-metrics-agent:latest
    environment:
      - API_KEY=${API_KEY}
      - COLLECTION_INTERVAL_SECONDS=5
      - RETENTION_HOURS=24
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - metrics_data:/data
    ports:
      - "8080:8080"
    restart: unless-stopped
```

---

## 📝 License

### Agent (Backend)
**MIT License** - Free and open source
- ✅ Use commercially
- ✅ Modify and distribute
- ✅ Private use
- ✅ No warranty

### Dashboard (Frontend)
**Proprietary** - Cloud-hosted SaaS
- ❌ Not open source
- ❌ Cannot self-host
- ✅ Free tier available
- ✅ Subscription required for advanced features

---

## 🤝 Contributing

We welcome contributions to the **agent (backend)**!

- 🐛 Bug reports: [GitHub Issues](https://github.com/yourorg/docker-metrics-agent/issues)
- 💡 Feature requests: [GitHub Discussions](https://github.com/yourorg/docker-metrics-agent/discussions)
- 🔧 Pull requests: [Contributing Guide](CONTRIBUTING.md)

**Note:** The dashboard/frontend is proprietary and not open for contributions.

---

## 📞 Support

### Free Tier
- 📧 Email: support@yourcloud.com
- 📚 Documentation: https://docs.yourcloud.com
- 💬 Community: Discord Server

### Pro Tier
- ⚡ Priority email support (24h response)
- 💬 Live chat support
- 📞 Phone support (business hours)

### Enterprise Tier
- 🎯 Dedicated account manager
- 📞 24/7 phone support
- 🔧 Custom integration assistance
- 📊 Quarterly business reviews

---

## 🎯 Roadmap

### Q1 2025
- [ ] Mobile app (iOS/Android)
- [ ] Advanced alerting rules
- [ ] Kubernetes support

### Q2 2025
- [ ] AI-powered anomaly detection
- [ ] Cost optimization recommendations
- [ ] Multi-cloud support (AWS, GCP, Azure)

### Q3 2025
- [ ] Custom plugins marketplace
- [ ] Terraform/Ansible integrations
- [ ] Advanced security scanning

---

## ❓ FAQ

**Q: Can I self-host the dashboard?**
A: No, the dashboard is cloud-only. However, the agent is open source and all your data stays on your infrastructure.

**Q: What happens if I exceed the free tier limit?**
A: You'll be prompted to upgrade to Pro. Your monitoring will continue, but you won't be able to add more systems.

**Q: Can I cancel anytime?**
A: Yes! Cancel anytime with no penalties. Your data remains accessible for 30 days after cancellation.

**Q: Is my data secure?**
A: Yes! Your metrics never leave your infrastructure. The cloud dashboard only fetches data on-demand when you view it.

**Q: Do you offer discounts?**
A: Yes! We offer:
- 20% off for annual billing
- 50% off for non-profits and educational institutions
- Custom pricing for enterprises (100+ systems)

---

**Ready to get started?**

👉 [Sign Up for Free](https://metrics.yourcloud.com/signup)

👉 [View Pricing](https://metrics.yourcloud.com/pricing)

👉 [Read Documentation](https://docs.yourcloud.com)

---

Made with ❤️ by the Docker Metrics Monitor team

---

## 🛠️ Development / Local Run

To run the entire stack (Dashboard + Agent) locally for development:

### 1. Start the Stack
```bash
docker-compose up --build -d
```
This starts:
- **Dashboard**: http://localhost:8082
- **Agent**: http://localhost:8081

### 2. Get API Keys
- **Dashboard Master Key**: Check `data/api_key.txt` (created after first run)
- **Agent API Key**: Check `agent_data/api_key.txt` (created after first run)

### 3. Connect Agent to Dashboard
1. Open Dashboard at [http://localhost:8082](http://localhost:8082)
2. Go to **Settings**
3. Enter the **Dashboard Master Key** to login.
4. Under "Monitored Systems", click **Add New System**:
   - **Name**: Local Agent
   - **URL**: `http://agent:8080` (Internal Docker Network - Recommended)
   - **API Key**: Enter the **Agent API Key**
5. Click **Add System**.

You should now see metrics flowing in!
