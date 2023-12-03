PROJECT_NAME=github.com/ciaolink-game-platform/cgp-chinese-poker-module
APP_NAME=chinese-poker.so
APP_PATH=$(PWD)
NAKAMA_VER=3.19.0

update-submodule-dev:
	git checkout develop && git pull
	git submodule update --init
	git submodule update --remote
	cd ./cgp-common && git checkout develop && git pull && cd ..
	go get github.com/ciaolink-game-platform/cgp-common@develop
update-submodule-stg:
	git checkout staging && git pull
	git submodule update --init
	git submodule update --remote
	cd ./cgp-common && git checkout main && cd ..
	go get github.com/ciaolink-game-platform/cgp-common@main

build:
	go mod tidy
	go mod vendor
	docker run --rm -w "/app" -v "${APP_PATH}:/app" heroiclabs/nakama-pluginbuilder:${NAKAMA_VER} build -buildvcs=false --trimpath --buildmode=plugin -o ./bin/${APP_NAME}

syncdev:
	rsync -aurv --delete ./bin/${APP_NAME} root@cgpdev:/root/cgp-server-dev/dist/data/modules/bin/
	# ssh root@cgpdev 'cd /root/cgp-server-dev && docker restart nakama_dev'

bsync: build sync

dev: update-submodule-dev build
stg: update-submodule-stg build
proto:
	protoc -I ./ --go_out=$(pwd)/proto  ./proto/bandarqq_api.proto
