# Gator

Gator is a command-line blog aggregator based on RSS (Really Simple Syndicate) feeds, built as a part of a learning guided project from boot.dev

## Prerequisites (Check these before installing Gator)

1. Postgres and psql should be installed. Gator uses postressql to store the users, feeds and articles information. To install postgres, follow the instructions in this link: https://www.postgresql.org/download/
2. Go should be installed. We will be installing gator using Go. To install go, follow instructions in this link: https://go.dev/doc/install 

## Installation

After checking the prerequisites, install gator by running the following command in your terminal: `go install https://github.com/vigneshsekar314/gator`
Gator uses a configuration file to store the information of the logged in user and postgres connection string (which stores the users, feeds and posts information).

Gator expects the configuration file to be present. Create the configuration file in your home directory with the name: `.gatorconfig.json`
The configuration file (`.gatorconfig.json`) should have the following json structure:

```{"db_url":"postgres://[username]:[password]@localhost:5432/gator?sslmode=disable"}```

In the above line, replace `[username]` with your postgres username and `[password]` with your postgres password for that user.

## Overview

Gator uses multiple commands to register a user, login as a different user, fetch, store, browse and follow articles from various feed links.
After installing Gator, register a new user using this command:

`gator register [username]`. Replace `[username]` with your name

The following commands are supported:

1. `gator register [username]`
2. `gator login [username]`
3. `gator users`
4. `gator addfeed [name] [url]`
5. `gator follow [url]`
6. `gator unfollow [url]`
7. `gator following`
8. `gator feeds`
9. `gator browse [number of articles]`
9. `gator agg [time duration]`


Run `gator agg [time duration]` as a background job to fetch latest feed articles and store their results. For example, `gator agg 1m` will fetch articles from following feeds every 1 minute.
Few other examples: Run `gator agg 1h` for 1 hour. Run `gator agg 1h30m` for 1 hour and 30 minutes.

