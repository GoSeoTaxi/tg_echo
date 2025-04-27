# tg\_echo – Minimal HTTP‑to‑Telegram notifier 🌐

This project exposes a single HTTP endpoint 💬

```plaintext
GET /tg?msg=Your+message
```

Every request is forwarded to a Telegram **group or channel** via a bot and annotated with UTC timestamp and caller IP 📡.

---

## Features 📜

| Feature                 | Description                                                        |
| ----------------------- | ------------------------------------------------------------------ |
| 💬 Message relay        | Any text sent via `msg` parameter is posted to Telegram            |
| 🕒 UTC timestamp        | Each post includes the moment the server received the request      |
| 🌐 Client IP            | Real client IP is appended (uses `X‑Real‑IP` / `X‑Forwarded-For`)  |
| 🚀 Zero‑downtime deploy | Provided GitHub Actions workflow builds, ships & reloads container |
| 🛑 Graceful shutdown    | SIGINT/SIGTERM → 5 s timeout → goodbye message                     |

---

## Quick start (local) ⚡

```bash
git clone https://github.com/GoSeoTaxi/tg_echo.git
cd tg_echo

# export your bot token & chat id
export BOT_TOKEN="123456789:ABC…"
export CHAT_ID="-1002355475020"

# run from source
go run ./cmd/tg_echo
# → listening on :8080

curl "http://localhost:8080/tg?msg=Hello%20world"   # Telegram receives it
```

---

## Docker 📦

```bash
docker build -t tg_echo .

docker run -d --name tg_echo \
  -e BOT_TOKEN="…" \
  -e CHAT_ID="…" \
  -e PORT="8080" \
  -p 8080:8080 \
  tg_echo:latest
```

---

## Production deploy with GitHub Actions ⚙️

> **TL;DR** Push to `or` → image is built, copied to your VPS over SSH and container restarted 🔄.

### 1. Create a user on the server 👤

```bash
# as root
useradd -m -s /bin/bash deploy
usermod -aG docker deploy     # let the user run docker
```

### 2. Generate an SSH key (on the server) 🔐

```bash
sudo -u deploy ssh-keygen -t ed25519 -f ~/.ssh/id_github_actions -N ""
cat ~/.ssh/id_github_actions.pub >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

Copy (private key) to your workstation – it will be added to GitHub secrets 🗝️.

### 3. Repository settings → *Secrets & Variables → Actions* 🔐

| Type         | Name        | Value                                 |
| ------------ | ----------- | ------------------------------------- |
| **Secret**   | `SSH_HOST`  | `your.server.ip`                      |
| **Secret**   | `SSH_USER`  | `deploy`                              |
| **Secret**   | `SSH_KEY`   | \*contents of \**`id_github_actions`* |
| **Secret**   | `BOT_TOKEN` | Telegram bot token                    |
| **Secret**   | `CHAT_ID`   | Chat or group ID                      |
| **Variable** | `PORT`      | `8080`                                |
| **Variable** | `LOG_LEVEL` | `info`                                |

### 4. Workflow overview 🔁

```plaintext
.github/workflows/deploy.yml
├── checkout code
├── docker build (with Buildx)
├── docker save → image.tar
├── scp image.tar ▶ ~/tg_echo on VPS
└── ssh: docker load & restart container
```

The container is started with:

```bash
docker run -d --name tg_echo --restart unless-stopped \
  -e BOT_TOKEN -e CHAT_ID -e PORT -e LOG_LEVEL \
  -p PORT:PORT tg_echo:latest
```

---

## Environment variables 🛠️

| Var         | Default | Description                                   |
| ----------- | ------- | --------------------------------------------- |
| `BOT_TOKEN` | –       | **Required**. Token from @BotFather           |
| `CHAT_ID`   | –       | **Required**. Chat / group numeric ID         |
| `PORT`      | `8080`  | HTTP listen port                              |
| `LOG_LEVEL` | `info`  | zap logger level (`debug`, `info`, `warn`, …) |

---

## API 🛎️

```plaintext
GET /tg?msg=Hello%20TG
```

### Response

- `200 OK` and body `ok`             – success
- `400 Bad Request` (`msg` missing) – client error
- `500 Internal Server Error`       – Telegram API failed

---

## License 📜

Released to the **public domain** under [The Unlicense](https://unlicense.org/).

