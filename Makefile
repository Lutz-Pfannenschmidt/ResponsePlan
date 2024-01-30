clean:
	go mod tidy
	rm -rf dist

build:
	make clean
	export DEV_MODE=false && go build -o dist/ResponsePlan main.go

build_all:
	make clean
	export DEV_MODE=false && export GOOS=windows && export GOARCH=amd64 && go build -o dist/ResponsePlan.exe main.go
	export DEV_MODE=false && export GOOS=linux && export GOARCH=amd64 && go build -o dist/ResponsePlan main.go


dev:
	make clean
	clear
	export DEV_MODE=true && gow -c -e=go,html,svg run .

sudev:
	make clean
	make build
	sudo ./dist/ResponsePlan --dev

run:
	make build
	clear
	./dist/ResponsePlan