echo Building Client...
go build -o built/client.exe src/client/main.go

echo Building Server...
go build -o built/server.exe src/server/main.go
