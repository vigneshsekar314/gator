# Gator

Gator is a command-line blog aggregator based on RSS (Really Simple Syndication) feeds, built as a part of a learning guided project from boot.dev

## Prerequisites (Check these before installing Gator)

1. Postgres and psql should be installed. Gator uses postressql to store the users, feeds and articles information. To install postgres, follow the instructions in this link: https://www.postgresql.org/download/
2. Go should be installed. We will be installing gator using Go. To install go, follow instructions in this link: https://go.dev/doc/install 

## Installation

After checking the prerequisites, install gator by running the following command in your terminal: `go install github.com/vigneshsekar314/gator@latest`
Gator uses a configuration file to store the information of the logged in user and postgres connection string (postgres SQL DB stores the users, feeds and posts information).

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
10. `gator agg [time duration]`
11. `gator reset`


Run `gator agg [time duration]` as a background job to fetch latest feed articles and store their results. For example, `gator agg 1m` will fetch articles from following feeds every 1 minute.
Few other examples: Run `gator agg 1h` for 1 hour. Run `gator agg 1h30m` for 1 hour and 30 minutes.

## Commands

### `gator register [username]`
Register a new user. Gator tracks feeds and posts by user. Each user will have different feeds they follow. When you register a new user, Gator automatically logs in as the new user.

### `gator login [username]`
Login as a different user. The login username should be first registered before loggin in.

### `gator users`
Lists all registered users in Gator database. The logged in user will be shown as (current) next to their name in the list.

### `gator addfeed [name] [url]`
Add a feed to retrieve and store articles. Name refers to the name of the website and Url should point to the RSS link of the website.

### `gator follow [url]`
Add a feed to retrieve and store articles. This command is used to follow a feed which was previously added by a different user. This command takes the same website name provided when this feed was initially added using `addfeed` command.

### `gator unfollow [url]`
Remove a feed from your lists.

### `gator following`
Shows a list of feeds being followed.

### `gator feeds`
Lists down all the feeds that are stored in the database. It shows all feeds stored by all users in Gator.

### `gator browse [number of articles]`
Displays the title of the article and its description. An optional argument [number of articles] limits the number of articles displayed. If [number of articles] is not provided, `browse` command by default, limits the number of articles as 2. 

### `gator agg [time duration]`
Aggregates the articles every n duration from different feeds. [time duration specifies] how long Gator waits before fetching feeds. **Do not specify duration for 1 second or less**, as it might cause a denial of service and may get your ip address blocked by the feed provider.

### `gator reset`
Deletes all users and feeds from Gator. **This action is permanent and not reversible**. Take caution before executing this command.
