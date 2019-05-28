# v2

A discord bot written in go.

[![Go report](http://goreportcard.com/badge/cabrinha/v2)](http://goreportcard.com/report/cabrinha/v2)

## Libraries

The main library in use is [discordgo](https://github.com/bwmarrin/discordgo).

Other libraries include, but are not limited to:

* [dgrouter](https://github.com/Necroforger/dgrouter)
* [logrus](https://github.com/sirupsen/logrus)
* [viper](https://github.com/spf13/viper)

# Usage

The bot expects a config file `config.yaml` to be present in the current working directory.

See `config.yaml.example` for an example config file.

## Building the Docker image

In order to build the docker image, simply run: 

```
docker build . -t v2:latest
```

## Running the Bot

In order to run the bot locally, simply run:

```
go run main.go
```

# TODO

The following features have yet to be implemented.

### Search

The bot will search the following endpoints and return the first (or best matching) result.

* Amazon - `!a`
* Google - `!g`
* Wikipedia - `!wiki`
* YouTube - `!yt`

### Quotes

A quote grabbing system, which stores quotes from users to be called back at a later time.

* Upon receiving a `!grab` command, the bot will store the last thing said in chat.
* The `!grab` command can take arguments such as a user mention or a word.
  * On receiving a user mention, the last thing said by the mentioned user will be grabbed.
  * On receiving word, the bot will look through the last n messages searching for that word and grab the first match.
  
To play back quotes (hopefully out of context to improve humor levels), a command `!rq` will be used to grab a quote at random.

The `!rq` command can take in the same arguments as `!grab`, either a user mention or a word.

* Upon user mention, a random quote is selected from all quotes stored under that user.
* Upon word, a random quote is selected from all quotes containing that word.

### Stocks and Crypto

The bot will take in a ticker symbol for stocks or crypto currencies and return the current (or latest) known price.

APIs to use are still TBD.
