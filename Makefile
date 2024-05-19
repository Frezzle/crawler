.PHONY: all
all: flowchart

.PHONY: flowchart
flowchart:
	docker run --rm -u `id -u`:`id -g` -v ./:/data minlag/mermaid-cli -i flowchart.mermaid -o flowchart.png --scale 4