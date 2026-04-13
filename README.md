# Introduction

Gator is a blog aggregator written in Go. It uses postgres to track users, feeds and posts.
<br>In order to run Gator you will need Go and Postgres.
<br>
Go instalation files can be found at https://go.dev/doc/install
### Install Postgres (Mac OS)
```
brew install postgresql@15
```

### Install Postgres (Linux)
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```
### Installing Gator
After cloning the repository, you can install gator by running ```go install``` inside the root project directory.
<br>
### Config
In order for Gator to work, you will need to set up a config file in your home directory ```~/.gatorconfig.json``` with the following content:
<br>
```
{
  "db_url": "postgres://user:password@host:port/database?sslmode=disable"
}
```
Be sure to replace the example values with your actual credentials.

### Commands
The following commands can be used with gator as follows:
<br><br>
```gator login <username>```
Logs you into the specified account.
<br><br>
```gator register <username>```
Adds the specified account to the users database.
<br><br>
```gator reset```
Resets the database's contents.
<br><br>
```gator addfeed <name> <url>```
Adds a new feed to the database.
<br><br>
```gator feeds```
Lists the currently available feeds.
<br><br>
```gator following```
Lists the currently followed feeds.
<br><br>
```gator follow <url>```
Adds a feed to the user's following list.
<br><br>
```gator unfollow <url>```
Removes a feed from the user's following list.
<br><br>
```gator agg <time_between_reqs>```
Scrapes the currently followed feeds and adds found posts to the posts database.
<br><br>
```gator browse <limit>```
Lists the specified number of posts, the default being 2.