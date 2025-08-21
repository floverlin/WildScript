$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o bin/linux/sigil .

$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o bin/mac/arm/sigil .

$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o bin/mac/amd/sigil .

$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o bin/win/sigil.exe .
