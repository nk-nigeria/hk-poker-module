PROJECT_NAME=github.com/ciaolink-game-platform/cgp-chinese-poker-module
APP_NAME=chinese-poker.so
APP_PATH=$(PWD)

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
	docker run --rm -w "/app" -v "${APP_PATH}:/app" heroiclabs/nakama-pluginbuilder:3.11.0 build -buildvcs=false --trimpath --buildmode=plugin -o ./bin/${APP_NAME}

sync:
	rsync -aurv --delete ./bin/${APP_NAME} root@cgpdev:/root/cgp-server/dev/data/modules/
	# ssh root@cgpdev 'cd /root/cgp-server && docker restart nakama'

bsync: build sync

dev: update-submodule-dev build
stg: update-submodule-stg build
proto:
	protoc -I ./ --go_out=$(pwd)/proto  ./proto/bandarqq_api.proto
