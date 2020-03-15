# Kumparan Test

## Tools
1. MySQL
2. Redis
3. Elasticsearch local

## Installation
1. Your local directory should be `bitbucket.org/michaelchandrag/kumparan-test`
2. Clone this repository `git clone https://github.com/michaelchandrag/kumparan-test.git`
3. Setup environment variable at `.env`, references on `.env.example`
4. Run `source .env`
5. Run `cd build`
6. Run `go build`
6. Run `go run main.go` or `./main`

## Environment Variables
```
export PORT=8080 -> Port used for running Go
export DB_HOST=localhost -> Database Host
export DB_USER=root -> Database username
export DB_PASS= -> Database password
export DB_NAME=kumparan -> Database name
export ES_HOST=http://localhost:9200 -> Elastic Search host
export Q_PORT=localhost:6379 -> Redis host
```

## Database
1. Make sure to create database with the same name on Environment Variables `DB_NAME`
2. SQL file can be found at this directory `database/news.sql`