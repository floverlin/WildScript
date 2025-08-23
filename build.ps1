$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o bin/linux/wild .

$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o bin/mac/arm/wild .

$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o bin/mac/amd/wild .

$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o bin/win/wild.exe .
