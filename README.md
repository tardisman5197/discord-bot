# discord-bot

A discord bot that does a bunch of things.

## Perquisites

* [Go 1.14](https://golang.org/) installed
* A Discord Bot Token

## Getting a Discord Bot Token

1. Visit this [site](https://discordapp.com/developers/applications)
2. Create a new application
3. Add a bot to the application
4. Copy the Bot Token

## Building

### Creating the Bot Executable

1. Navigate to the repo directory
2. Run this command

    ```bash
    go build
    ```

This has now created an executable with the same name as the repo directory.

### Adding the bot to the server

1. Visit this [site](https://discordapp.com/developers/applications)
2. Select the application you created
3. Select the `OAuth2` section of the application
4. Under Scopes select
    * bot
5. Under Bot Permission
    * Send Messages
    * Manage Messages
    * Read Message History

## Running

1. Navigate to the executable
2. Run the executable

    ```bash
    ./discord-bot.exe -t <token>
    ```

    For example

    ```bash
    ./discord-bot.exe -t ASDFGHJ.1234.ASDFGHJ
    ```
