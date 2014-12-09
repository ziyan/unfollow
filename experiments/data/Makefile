.PHONY: all
all: twitter.txt nodes.json

twitter.txt: twitter.txt.gz
	gzip -cd twitter.txt.gz > twitter.txt

nodes.json: nodes.json.gz
	gzip -cd nodes.json.gz > nodes.json

.PHONY: clean
clean:
	rm -f twitter.txt nodes.json
