# Visit Service

A visit-tracking backend for [shivankkapoor.com](https://shivankkapoor.com). Records page visits, enriches them with IP geolocation via [Meridian](https://github.com/ShivankKapoor/meridian-go), sends Discord notifications, and generates daily summaries.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Service info |
| `POST` | `/track` | Record a page visit |
| `GET` | `/health` | Database health check |
| `GET` | `/admin/run-summary` | Manually trigger daily summary |

### POST /track

```json
{
  "ipAddress": "1.2.3.4",
  "pageVisited": "/projects",
  "deviceInfo": "Mozilla/5.0 ..."
}
```

## Setup

```bash
cp .env.example .env  # fill in values
go run main.go
```

Service runs on **port 8088**.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `DISCORD_WEBHOOK_URL` | Discord webhook for visit notifications |
| `MERIDIAN_URL` | Base URL for internal IP geolocation cache service |
| `CUSTOM_EMOJI_ID` | Optional custom emoji for Discord messages |
| `PROD` | `true` enables CORS restrictions and referer validation |
| `ALLOWED_ORIGINS` | Comma-separated origins allowed in production |

## Architecture

```
POST /track
  └── AllowedReferer middleware     # referer check (prod only)
  └── RateLimit middleware          # 20 req/min per IP
  └── InsertPageVisit               # sync DB write
  └── goroutine
        ├── GetLocation (Meridian)  # IP geolocation with Redis cache
        └── SendVisitMessage        # Discord webhook
```

**Daily cron** (midnight CST): aggregates yesterday's visits into `daily_visit_stats` and sends a Discord summary.

## Docker

```bash
docker build -t visit-service .
docker run -p 8088:8088 visit-service
```

## Database

PostgreSQL with two tables:

- `page_visits` — raw visit records
- `daily_visit_stats` — daily aggregated counts
