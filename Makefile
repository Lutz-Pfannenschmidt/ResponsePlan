.PHONY: build run clean

build:
	go build -o ./dist/ResponsePlan main.go

dev: build
	cd web && bun run build && cd ..
	sudo ./dist/ResponsePlan -d

clean:
	rm -rf dist