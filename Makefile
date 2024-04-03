clean:
	go mod tidy
	rm -rf dist

build:
	make clean
	mkdir -p dist
	cd tty2web && make tty2web
	cp tty2web/tty2web dist/tty2web
	go build -o dist/ResponsePlan main.go

dev:
	make clean
	mkdir -p dist
	cd tty2web && make tty2web
	cp tty2web/tty2web dist/tty2web
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