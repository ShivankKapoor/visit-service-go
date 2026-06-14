# Visit Service

A visit-tracking backend for [shivankkapoor.com](https://shivankkapoor.com). Records page visits, enriches them with IP geolocation via [Meridian](https://github.com/ShivankKapoor/meridian-go), sends Discord notifications, and generates daily summaries.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Service info |
| `POST` | `/track` | Record a page visit |
| `GET` | `/health` | Database health check |
| `GET` | `/run-daily-summary` | Manually trigger daily summary (non-prod only) |

### POST /track

```json
{
  "pageVisited": "/projects",
  "deviceInfo": "Mozilla/5.0 ..."
}
```

## Setup

```bash
cp example.env .env  # fill in values
go run main.go
```

Service runs on **port 8088**.

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `DISCORD_WEBHOOK_URL` | Yes | Discord webhook for visit and daily summary notifications |
| `MERIDIAN_URL` | Yes | Base URL for internal IP geolocation service |
| `PROD` | No | Set to `true` to enable CORS restrictions and referer validation |
| `ALLOWED_ORIGINS` | In prod | Comma-separated list of allowed origins e.g. `https://shivankkapoor.com,https://www.shivankkapoor.com` |
| `CUSTOM_EMOJI_ID` | No | Custom Discord emoji to append to notifications |

## Architecture

```
POST /track
  └── CORS middleware
  └── RateLimit middleware          # 20 req/min per IP (evicts after 5 min idle)
  └── AllowedReferer middleware     # host-based referer check (prod only)
  └── InsertPageVisit               # DB write (context.Background, not tied to request)
  └── goroutine
        ├── GetLocation (Meridian)  # IP → city, region, country
        └── SendVisitMessage        # Discord webhook (async)
```

**Daily cron** runs at midnight CST (DST-aware). Counts visits from the previous CST day, upserts into `daily_visit_stats`, and sends a Discord summary.

## Database

PostgreSQL. Schema in `schema.sql`.

- `page_visits` — raw visit records with IP, page, device, user agent, and timestamp
- `daily_visit_stats` — daily aggregated visit counts keyed by date

## Docker

```bash
docker build -t visit-service .
docker run --env-file .env -p 8088:8088 visit-service
```
