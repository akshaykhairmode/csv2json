build: build-windows-x64 build-windows-x32 build-darwin build-darwin-m1 build-linux-32 build-linux-64

build-windows-x64:
	GOOS=windows GOARCH=amd64 go build -o dist/win64/csv2json.exe main.go

build-windows-x32:
	GOOS=windows GOARCH=386 go build -o dist/win32/csv2json.exe main.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/csv2json main.go

build-darwin-m1:
	GOOS=darwin GOARCH=arm64 go build -o dist/darwin-m1/csv2json main.go

build-linux-32:
	GOOS=linux GOARCH=386 go build -o dist/linux32/csv2json main.go

build-linux-64:
	GOOS=linux GOARCH=amd64 go build -o dist/linux64/csv2json main.go
	