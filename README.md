# Intro

# Requirements
1. Install go 1.23+
2. Install postgresql v15 or later
3. Install goose `go install github.com/pressly/goose/v3/cmd/goose@latest`
4. Install SQLC `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

# Config
1. `touch ~/.gatorconfig.json`
2. Change username to your own `echo "{ \"db_url\": \"postgres://username:@localhost:5432/gator" }"` >> ~/.gatorconfig.json
3. run migrations `goose postgres "postgres://username:@localhost:5432/gator" up`