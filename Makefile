source = ./cmd/display/main.go
arm_v7_dest = ./bin/ssc_hmc_display

dev:
	go run $(source)

lint:
	golangci-lint run -c ./.golangci.yml ./...

arm-v7-compile:
	GOOS=linux GOARCH=arm GOARM=7 go build -o $(arm_v7_dest) $(source)

clean-arm-v7-compile:
	rm $(arm_v7_dest)