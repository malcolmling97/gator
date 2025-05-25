# 🐊 gator - A Minimalist RSS Aggregator CLI

`gator` is a command-line RSS feed aggregator that lets you follow, fetch, and browse RSS feeds directly from your terminal. It supports background aggregation, user logins, and stores posts in a PostgreSQL database for easy viewing later.

---

## 🚀 Features

- User login & registration
- Add and follow RSS feeds (e.g. TechCrunch, Hacker News)
- Background aggregation loop (`agg`) that fetches new posts periodically
- Stores posts in PostgreSQL with deduplication by URL
- Browse followed feed posts in the terminal

---

## 📦 Requirements

To run `gator`, you’ll need:

- [Go](https://golang.org/dl/) 1.20 or newer
- [PostgreSQL](https://www.postgresql.org/download/) installed and running

---

## 🛠 Installation

Clone the repo and install the CLI:

```bash
git clone https://github.com/your-username/gator.git
cd gator
go install
```

## ⚙️ Configuration
Create a .gatorconfig file in your home directory with the following content:

toml
Copy
Edit

```
# ~/.gatorconfig
db_url = "postgres://user:password@localhost:5432/gator?sslmode=disable"
current_user_name = ""
```
Replace the db_url with your actual PostgreSQL connection string.

## 🧪 Example Usage
### Register or Login
```bash
gator register yourusername
gator login yourusername
```
### Add and Follow Feeds

```bash
gator addfeed techcrunch https://techcrunch.com/feed/
gator follow https://techcrunch.com/feed/
```
### Start Aggregation Loop (every 30 seconds)
```bash
gator agg 30s
```

### Browse Posts
```bash
gator browse       # default: show 2 recent posts
gator browse 5     # show 5 recent posts
```

## 🧹 Development
Run commands directly for testing:

```bash
go run . addfeed dev https://blog.boot.dev/index.xml
go run . browse 3
```

## 🐘 Migrations
This project uses goose for schema migrations.

Install goose (if you haven’t already):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Apply migrations:

```bash
goose postgres "postgres://user:pass@localhost:5432/gator?sslmode=disable" up
```

## 📂 Project Structure
```bash
.
├── internal/
│   ├── config/           # Config loader
│   ├── database/         # sqlc-generated DB layer
│   └── rss/              # RSS feed fetching
├── sql/
│   ├── queries/          # Config loader
│   ├── schema/           # DB migrations (goose)
├── main.go               # CLI entry point
└── README.md
```