# Gator

A CLI RSS feed aggreGATOR written in Go

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

## Requirements
This application requires a postgres database, which is used to store users, RSS feeds, and RSS feed posts.

### Config file
A json config file is used for storing database connection and the current user.  The config file must be called ".gatorconfig.json" and located in your root home directory, example for linux, ~/.gatorconfig.json.

Example config file

`{
  "db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name":"user1"
}`

#### Installation
Application can be installed by running the following command

`go install github.com/jkk290/gator@latest`

##### Commands/Usage

`gator register username`

Creates a new user and add them to the database

Example: `gator register user1`

`gator login username`

Sets the passed username as the current user in the config file

Example: `gator login user1`

`gator users`

List all users and shows who is set as the current user

`gator agg duration`

Continuously fetches new posts from feeds based off the passed duration.  Press ctrl+c to exit the application.

Durations examples: 1s, 1m, 1h

Example: `gator agg 5m` - Continuously fetches new posts every 5 minutes.

`gator addfeed "title" "url"`

Adds the given feed using the title and url to the database and sets the current user to follow that feed

Example: `gator addfeed "Boot.dev" "https://blog.boot.dev/index.xml"`

`gator feeds`

List all feeds and which user added them

`gator follow "url"`

Add current user as a follower to the passed feed

Example: `gator follow "https://blog.boot.dev/index.xml"`

`gator following`

Displays list of feeds current user is following

`gator unfollow "url"`

Removes current user as a follower to the passed feed

Example: `gator unfollow "https://blog.boot.dev/index.xml"`

`gator browse [limit]`

Displays posts from followed feeds ordered by most recent published date.  If no limit is passed, 2 posts will be displayed.

Example: `gator browse 5`  - Will display 5 most recently published posts from followed feeds
