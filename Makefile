.PHONY: all
all: build

.PHONY: build
build:
	go build -o build/app

.PHONY: run
run: build
	./build/app

.PHONY: flowchart
flowchart:
	docker run --rm -u `id -u`:`id -g` -v ./:/data minlag/mermaid-cli -i flowchart.mermaid -o flowchart.png  --configFile="mermaidConfig.json" --scale 4
