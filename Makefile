bin: bin/ilogin_darwin bin/ilogin_linux bin/ilogin_windows

bin/ilogin_darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/ilogin_darwin cmd/ilogin/*.go

bin/ilogin_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/ilogin_linux cmd/ilogin/*.go

bin/ilogin_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/ilogin_windows cmd/ilogin/*.go

docker:
	docker build -t ilogin .