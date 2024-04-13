clean:
	go mod tidy
	rm -rf dist

build:
	make clean
	mkdir -p dist
	go build -o dist/ResponsePlan main.go

dev:
	make clean
	mkdir -p dist
	clear
	gow -c -e=go,html,svg run . -d

sudev:
	make build
	clear
	sudo ./dist/ResponsePlan -d

run:
	make build
	clear
	./dist/ResponsePlan

install:
	make build
	sudo cp dist/ResponsePlan /usr/local/bin/ResponsePlan
	make clean