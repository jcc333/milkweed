# milkweed

An RSS to BlueSky Bot.
Because, uh, milkweed feeds monarch butterflies.
And RSS is a format that is in terms of feeds.
Listen, I'm not a smart man.

## Usage

Milkweed is entirely configered via environment variables, so the binary itself simply runs as `milkweed`.

## Configuration

Needs the following environment variables set:

```.environment
MILKWEED_USERNAME=<BlueSky Username>
MILKWEED_PASSWORD=<BlueSky Password>
MILKWEED_RSS_URL=<the URL where the RSS feed lives>
MILKWEED_SQLITE_PATH=<the path to the SQLITE DB>
MILKWEED_CRON_SCHEDULE=<a crontab schedule to refresh the RSS feed>
```
