# Gator
Gator is a blog (RSS feed) aggregator using GO.

## Requirements
1. Install go 1.23+
2. Install postgresql v15 or later

## Tools used
1. SQLC `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` SQL queries to go code

## Config
1. `touch "$HOME/.gatorconfig.json"`
2. Change username to your own `echo "{ \"db_url\": \"postgres://username:@localhost:5432/gator" }" >> "$HOME/.gatorconfig.json"`
3. goose `go install github.com/pressly/goose/v3/cmd/goose@latest` for the creation of db migrations
4. run migrations `goose up`

## Install
Once the requirements are installed and the config is set up, run:
1. `go install`

## Usage
1. `gator register user`
2. `gator addfeed TechCrunch https://techcrunch.com/feed/`

## Feeds to add
- TechCrunch: https://techcrunch.com/feed/
- Hacker News: https://news.ycombinator.com/rss
- Boot.dev Blog: https://blog.boot.dev/index.xml
