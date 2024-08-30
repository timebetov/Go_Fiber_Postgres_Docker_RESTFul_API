build:
	go build -o server cmd/main.go

run: build
	./server

watch:
	reflex -s -r '\.go$$' make run