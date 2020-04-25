# discord-bot

A discord bot that does a bunch of things.

## Perquisites

* A Discord Bot Token
* Docker installed
* [Go 1.14](https://golang.org/) installed (Optional)

## Getting a Discord Bot Token

1. Visit this [site](https://discordapp.com/developers/applications)
2. Create a new application
3. Add a bot to the application
4. Copy the Bot Token

## With Docker

### Building Docker Containers

1. Navigate to the Repo
2. Run the docker-compose build command

    ```bash
    docker-compose build
    ```

    *Note: On linux this should be done as sudo*

### Running Docker Containers

1. Navigate to the repo
2. Copy the `TEMPLATE.env` to `.env`
3. Fill in the `.env` file

    *Note: Ensure that the mongo data directory has been created*

    *Note: On a Windows host a drive must be shared first on docker hub*
    
4. Run the containers

    ```bash
    docker-compose up
    ```

    *Note: On linux this should be done as sudo*

    *Note: This command can be ran with a `-d` to run in detached mode*

### Stopping Docker Containers

1. Navigate to the repo
2. Stop the containers

    ```bash
    docker-compose down
    ```

    *Note: On linux this should be done as sudo*

## Without Docker

A note about this, a mongo database needs to be running.
(I have no idea how to do that without docker)

### Building

#### Creating the Bot Executable

1. Navigate to the repo directory
2. Run this command

    ```bash
    go build
    ```

This has now created an executable with the same name as the repo directory.

#### Adding the bot to the server

1. Visit this [site](https://discordapp.com/developers/applications)
2. Select the application you created
3. Select the `OAuth2` section of the application
4. Under Scopes select
    * bot
5. Under Bot Permission
    * Send Messages
    * Manage Messages
    * Read Message History

### Running

1. Navigate to the executable
2. Run the executable

    ```bash
    ./discord-bot.exe -t <token>
    ```

    For example

    ```bash
    ./discord-bot.exe -t ASDFGHJ.1234.ASDFGHJ
    ```
