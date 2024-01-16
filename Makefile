build:
	GOOS=windows GOARCH=amd64 go build -o bin/NotCoinBot_windows_x64.exe .
	GOOS=darwin GOARCH=amd64 go build -o bin/NotCoinBot-amd64-macOS .
	GOOS=linux GOARCH=amd64 go build -o bin/NotCoinBot-amd64-linux .