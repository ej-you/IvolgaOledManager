# Russian southern scientific center project

## Program for managing OS of hydro-meteorological center with OLED display and buttons

## Hardware

1. BeagleBone Black platform (ARMv7)
2. OLED display 128x64 (I2C, SSD1306)
3. GPIO buttons with external pull-up to HIGH

## Production

### Compile

```shell
GOOS=linux GOARCH=arm GOARM=7 go build -o ./ssc_hmc_display ./cmd/ssc_hmc_display/main.go
```

### Config

Config file must be named `config.yml` and located in the same directory as the binary.

### Secret

DB connection settings used via ENV-variables.
By default app using tcp connection. Set DB_USE_SOCKET="false" and
specify DB_SOCKET to use unix-socket for connection.

For all DB env-variables use the next manual:

```env
DB_USER="db_user"
DB_PASSWORD="cool-password"
DB_NAME="example"

DB_HOST="127.0.0.1"
DB_PORT="3306"

DB_SOCKET="/var/run/mysqld/mysqld.sock"
DB_USE_SOCKET="false"
```

### Assets

Assets directory contains greetings image (path to it must be changed via `config.yml`)

### Service

Example of linux service for the program can be found in the `/build` directory

## DB

DB contains one table with logs from another app.

- ID - id of log record.
- Level - log level (see table `Log levels` below).
- Header - short log title.
- Content - log description.
- CreatedAt - datetime tag of log record.

### Log levels

| number | value |
|--------|-------|
| 0      | trace |
| 1      | debug |
| 2      | info  |
| 3      | warn  |
| 4      | error |
| 5      | fatal |
