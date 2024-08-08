PROJECT_NAME=github.com/nakamaFramework/cgp-chinese-poker-module
APP_NAME=chinese-poker.so
APP_PATH=$(PWD)

NAKAMA_VER=3.19.0

update-submodule-dev:
	git checkout develop && git pull origin develop
	git submodule update --init
	git submodule update --remote
	cd ./cgp-common && git checkout develop && git pull origin develop && cd ..
	go get github.com/nakamaFramework/cgp-common@develop
update-submodule-stg:
	git checkout staging && git pull origin staging
	git submodule update --init
	git submodule update --remote
	cd ./cgp-common && git checkout staging && git pull origin staging && cd ..
	go get github.com/nakamaFramework/cgp-common@staging

build:
	# ./sync_pkg_3.11.sh
	go mod tidy
	go mod vendor
	docker run --rm -w "/app" -v "${APP_PATH}:/app" "heroiclabs/nakama-pluginbuilder:${NAKAMA_VER}" build -buildvcs=false --trimpath --buildmode=plugin -o ./bin/${APP_NAME}

syncdev:
	rsync -aurv --delete ./bin/${APP_NAME} root@cgpdev:/root/cgp-server-dev/dist/data/modules/bin/
	ssh root@cgpdev 'cd /root/cgp-server-dev && docker restart nakama_dev'

syncstg:
	rsync -aurv --delete ./bin/${APP_NAME} root@cgpdev:/root/cgp-server/dist/data/modules/bin
	ssh root@cgpdev 'cd /root/cgp-server && docker restart nakama'

dev: update-submodule-dev build

stg: update-submodule-stg build