# gator-go-cli

## Requirements
- PostgreSQL
- Go

## Installation
```bash
go install github.com/amengdv/gator
```

## Usage

Create a config repo in your home directory called
`.gatorconfig.json`

The inital content of the file should be
```
{"db_url":"postgres://example"}
```
replace the url accordingly

Some command you can run

```bash
gator register <name>
gator login <name>
gator addfeed <feed name> <url>
gator follow <feed url>
gator unfollow <feed url>
gator following
gator browse <limit> (optional) default: 2
```

To run the aggregator
```bash
gator agg <interval time> e.g "1s", "1h", "1m", "4s"
```

