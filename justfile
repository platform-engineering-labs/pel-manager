default: clean build

clean:
	rm -rf ./dist

build:
	mkdir -p dist/bin
	mkdir -p tmp/
	go build -C ./cmd -o ../dist/bin/pelmgr