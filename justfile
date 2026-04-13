export OS := `uname | tr '[:upper:]' '[:lower:]'`
export ARCH := `uname -m |  tr -d '_' | sed s/aarch64/arm64/`

default: clean build

clean:
	rm -rf ./dist

build:
	mkdir -p dist/bin
	go build -C ./cmd -o ../dist/bin/pelmgr

publish-setup:
    aws s3 cp ./scripts/setup.sh s3://hub.platform.engineering/get/setup.sh
    aws s3 cp ./scripts/formae.sh s3://hub.platform.engineering/get/formae.sh

publish-bin: build
    aws s3 cp ./dist/bin/pelmgr s3://hub.platform.engineering/get/binaries/{{OS}}-{{ARCH}}/pelmgr
