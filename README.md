# Russian southern scientific center project

## Program for managing OS of hydro-meteorological center with OLED display and buttons

## Hardware

1. BeagleBone Black platform (ARMv7)
2. OLED display 128x64 (I2C, SSD1306)

## Production

### Compile

```shell
GOOS=linux GOARCH=arm GOARM=7 go build -o ./ssc_hmc_display ./cmd/ssc_hmc_display/main.go
```

### Config

Config file must be named `config.yml` and located in the same directory as the binary.
