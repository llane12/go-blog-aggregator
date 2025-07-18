# gator: Blog Aggregator in Go

Scrapes registered blogs (RSS feeds), parses the XML and saves the posts to a database.

Has a simple CLI interface, see commands below

## Setup
Requires Go, Goose and Postgres

### Install Go
See [Go docs](https://go.dev/doc/install)

For WSL:
```
wget https://dl.google.com/go/go1.23.0.linux-amd64.tar.gz
sudo tar -xvf go1.23.0.linux-amd64.tar.gz
sudo mv go /usr/local
echo "export GOROOT=/usr/local/go" >> ~/.bashrc
echo "export GOPATH=\$HOME/go" >> ~/.bashrc
echo "export PATH=\$GOPATH/bin:\$GOROOT/bin:\$PATH" >> ~/.bashrc
source ~/.bashrc
```

### Install Postgres (Linux/WSL)

1. Install Postgres v15 or later:
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

2. Ensure the installation worked:
```
psql --version
```

3. Update postgres password:
```
sudo passwd postgres
```
Enter a simple password e.g. "postgres".

4. Start the Postgres server in the background
```
sudo service postgresql start
```

5. Connect to the server using the psql shell
```
sudo -u postgres psql
```

6. Create a the database:
```
CREATE DATABASE gator;
```

7. Connect to the new database:
```
\c gator
```

8. Set the user password:
```
ALTER USER postgres PASSWORD 'postgres';
```

### Config file
A JSON config file .gatorconfig.json is expected in your home directory.
On Windows this is probably the root of your user folder e.g. C:\Users\<user_name>\.gatorconfig.json
On Linux, you already know where your home directory is ;)

You will need a database connection string for the config file and to create the database using goose

Format: `postgres://postgres:<user>@<server>:<port>/gator`

Example: `postgres://postgres:postgres@localhost:5432/gator?sslmode=disable`

Add `?sslmode=disable` if your OS username doesn't match the postgres user

```
{
    "db_url":"postgres://postgres:<user>@<server>:<port>/gator?sslmode=disable"
}
```

Then create the database:
```
cd sql/queries
goose postgres <db_connection_string> up
```

## Commands
- `register <user_name>`
- `login <user_name>`
- `users`: list all registered users

- `addfeed <user_name> <feed_url>`
- `feeds`: list all registered feeds

- `follow <url>`: follow the given feed, for the logged-in user
- `following`: list the followed feed, for the logged-in user
- `unfollow <url>`: un-follow the given feed, for the logged-in user

- `agg <limit e.g 10s, 1m>`: long-running scrape of all registered feeds, pausing for the passed duration between feeds
- `browse [num posts, default 2]`: (after running agg to scrape feeds) list the n last posts from feeds followed by the logged-in user

- `reset`: clear database
