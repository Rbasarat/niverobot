# niverobot
Telegram chatbot build in Go

## Installation
Install and run Niverobot with Docker-compose:
1. Install [Docker](https://docs.docker.com/get-docker/)
2. Create a file named `docker-compose.yml`
3. Copy the following contents to the `docker-compose.yml` file:
```yaml
version: '3.1'

services:
  niverobot:
    image: ghcr.io/rbasarat/niverobot:1.0.1
    restart: unless-stopped
    networks:
      - default
    depends_on:
      - database
    environment:
      BOT_API_TOKEN: "InsertBotFatherTokenHere"
      DB_HOST: database
      DB_USER: niverobot
      DB_PASSWORD: "InsertSecureDatabasePasswordHere"
      DB_SCHEMA: niverobot
  database:
    image: library/postgres:13.4
    networks:
      - default
    ports:
      - 5432:5432
    volumes:
      - C:/Users/rasja/tmp:/var/lib/postgresql/data
    environment:
      POSTGRES_DB: niverobot
      POSTGRES_USER: niverobot
      POSTGRES_PASSWORD: "InsertSecureDatabasePasswordHere"

networks:
  default:
```
5. Change the `POSTGRES_PASSWORD` and `DB_PASSWORD` keys to a secure password.
6. Get your API token from [Botfather](https://core.telegram.org/bots#6-botfather)
7. Change the value of ``BOT_API_TOKEN`` to the API token retrieved from Bothfather.
8. Run docker `compose up -d` to start the Niverobot.