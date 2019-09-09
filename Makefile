all: build
build:
	go build -o build/xmlpurifier .
	GOOS=windows GOARCH=386 go build -o build/xmlpurifier.exe .