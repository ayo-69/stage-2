Below is a **clean, structured `README.md`** with **clear sections** — perfect for your **Go + PostgreSQL Country API** submission.

---

```markdown
# Country Currency & Exchange API

**A RESTful backend that fetches country data, computes estimated GDP, caches results, and generates a summary image.**

Built with:
- **Go** (Gin, GORM)
- **PostgreSQL**
- **In-Memory Cache**
- **Image Generation** (`gg`)
- **External APIs**: `restcountries.com`, `open.er-api.com`

---

## Table of Contents

- [Features](#features)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Local Setup](#local-setup)
- [Testing](#testing)
- [Deployment](#deployment)
- [Submission Instructions](#submission-instructions)
- [Troubleshooting](#troubleshooting)
- [Author](#author)

---

## Features

| Feature | Description |
|--------|-------------|
| `POST /countries/refresh` | Fetch + cache data + generate PNG |
| `GET /countries` | Filter by `region`, `currency`, `sort=gdp_desc` |
| `GET /countries/:name` | Get single country |
| `DELETE /countries/:name` | Remove country |
| `GET /status` | Total count + last refresh |
| `GET /countries/image` | Download `summary.png` |
| **In-Memory Cache** | TTL via `CACHE_TTL_SECONDS` |
| **GDP Calculation** | `population × random(1000–2000) ÷ rate` |
| **Error Handling** | 400, 404, 500, 503 JSON responses |

---

## API Endpoints

| Method | Endpoint | Query Params | Description |
|-------|----------|-------------|-----------|
| `POST` | `/countries/refresh` | — | Refresh data + image |
| `GET` | `/countries` | `region`, `currency`, `sort` | List with filters |
| `GET` | `/countries/:name` | — | Get one country |
| `DELETE` | `/countries/:name` | — | Delete country |
| `GET` | `/status` | — | Stats |
| `GET` | `/countries/image` | — | Summary PNG |

---

## Project Structure

```
country-api-go/
├── main.go
├── go.mod
├── go.sum
├── .env
├── .env.example
├── cache/
│   └── summary.png
├── models/
│   └── country.go
├── services/
│   ├── country_service.go
│   ├── exchange_service.go
│   └── image_service.go
├── controllers/
│   └── country_controller.go
├── utils/
│   └── cache.go
├── test_api.sh
└── README.md
```

---

## Prerequisites

- Go `1.21+`
- PostgreSQL `15+`
- Docker (optional)
- `curl`, `jq`

---

## Local Setup

### 1. Clone Repository

```bash
git clone https://github.com/yourname/country-api-go.git
cd country-api-go
```

### 2. Start PostgreSQL (Docker)

```bash
docker run -d \
  --name pg-test \
  -e POSTGRES_PASSWORD=pass \
  -e POSTGRES_DB=countries_db \
  -p 5432:5432 \
  postgres:15
```

### 3. Configure Environment

```bash
cp .env.example .env
```

Edit `.env`:

```env
PORT=3000
DATABASE_URL=host=localhost user=postgres password=pass dbname=countries_db port=5432 sslmode=disable
CACHE_TTL_SECONDS=3600
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run Server

```bash
go run main.go
```

API: `http://localhost:3000`

---

## Testing

### Run Full Test Suite

```bash
chmod +x test_api.sh
./test_api.sh
```

**Expected Output**:
```
All tests passed! API is ready for submission.
```

Tests cover:
- Refresh
- Filters & sorting
- Single country
- Image generation
- Status
- Error cases

---

## Deployment

### Deploy to Railway

```bash
npm install -g @railway/cli
railway login
railway init
railway up
```

Set in **Railway Dashboard**:
- `DATABASE_URL` → PostgreSQL add-on
- `CACHE_TTL_SECONDS` → `3600`

Live URL: `https://your-app.up.railway.app`

---

## Submission Instructions

1. **Deploy** your app
2. **Run**:
   ```bash
   ./test_api.sh https://your-app.up.railway.app
   ```
   → Must pass all tests
3. Go to **#stage-2-backend** in Slack
4. Run:
   ```
   /stage-two-backend
   ```
5. Submit:
   - **API Base URL**
   - **GitHub Repo URL**
   - **Full Name**
   - **Email**

---

## Troubleshooting

| Issue | Solution |
|------|----------|
| `connection refused` | Start PostgreSQL |
| `port 3000 in use` | `lsof -ti:3000 \| xargs kill -9` |
| `image not found` | Wait 5s after `/refresh` |
| `build error` | `go mod tidy` |
| `password failed` | Check `.env` password |

---

## Author

**Your Name**  
Backend Engineer | Go | PostgreSQL | APIs

---

**Ready?**  
Run `./test_api.sh` → **All tests passed!** → Submit via Slack.

**Good luck!**
```

---

### Save This File

```bash
cat > README.md << 'EOF'
[PASTE THE ABOVE CONTENT]
EOF
```

Then:
```bash
git add README.md
git commit -m "Add complete README with sections"
git push
```

---

**Your README is now submission-ready**  
Clear, professional, and structured.

Let me know when you deploy — I’ll help you **submit in 30 seconds**.
