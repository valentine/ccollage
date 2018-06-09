build:
	go run scripts/generate.go
	go build ./cmd/ccollage

generate:
	go run scripts/generate.go

dev:
	go run -tags dev cmd/ccollage/main.go

.PHONY: build generate dev
