GOOS=linux   GOARCH=amd64 go build -o issuetracker-linux-amd64;
GOOS=darwin  GOARCH=arm64 go build -o issuetracker-macos-arm64;
GOOS=windows GOARCH=amd64 go build -o issuetracker-windows.exe;
GOOS=linux   GOARCH=arm64 go build -o issuetracker-linux-arm64
