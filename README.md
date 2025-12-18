# openmon

A lightweight monitor for OpenRouter that tracks new model releases and sends rich notifications to a Discord webhook.

## Features

- **Real-time Monitoring**: Checks for new models added to OpenRouter every minute.
- **Rich Discord Notifications**: Sends detailed embeds including:
    - Model Name and Slug (with direct links).
    - Modality (Input 🡒 Output).
    - Context window size.
    - Pricing (Input/Output per 1M tokens).
    - Model Description.
- **Simple Configuration**: Easy to set up via a single YAML file.
- **Systemd Support**: Includes ready-to-use service files for Linux deployments.

## Installation

### Binary
Download the latest pre-built binary for your architecture from the [Releases](https://github.com/coalaura/openmon/releases/latest) page.

### Source
```bash
go build -o openmon .
```

## Configuration

Create a `config.yml` in the same directory as the binary:

```yaml
api-key: "sk-or-v1-your-api-key"
webhook: "https://discord.com/api/webhooks/your-webhook-url"
```

## Deployment (systemd)

The repository provides automated setup for systemd in the `conf/` directory:

1. Navigate to the `conf` directory.
2. Review the `openmon.service` and `openmon.conf` files.
3. Run the provided setup script:
   ```bash
   chmod +x setup.sh
   ./setup.sh
   ```

## How it works

1. On startup, `openmon` fetches the current list of available models from OpenRouter to establish a baseline.
2. Every 60 seconds, it refreshes the list.
3. If new models (detected by `created_at` timestamp) are found, it generates a Discord embed and sends it to the configured webhook.

## [License](LICENSE)
