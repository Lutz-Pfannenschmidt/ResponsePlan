clean:
	go mod tidy
	rm -rf dist

build:
	make clean
	go build -o dist/ResponsePlan main.go

dev:
	make clean
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