PROJECT_NAME=github.com/ciaolink-game-platform/cgp-bing-module
APP_NAME=bing.so
APP_PATH=$(PWD)

build:
	#docker run --rm -w "/go/src/${PROJECT_NAME}" -v "${APP_PATH}:/go/src/${PROJECT_NAME}" heroiclabs/nakama-pluginbuilder:3.3.0 build --buildmode=plugin -o ./bin/${APP_NAME}
	docker run --rm -w "/app" -v "${APP_PATH}:/app" heroiclabs/nakama-pluginbuilder:3.4.0 build --trimpath --buildmode=plugin -o ./bin/${APP_NAME}
	
sync:
	rsync -aurv --delete ./bin/${APP_NAME} cgpdev:/root/cgp-server/data/modules/

#buildplugin:
#	go build --trimpath --buildmode=plugin -o ${APP_NAME} .
#go mod vendor
#DOCKER_BUILDKIT=0 docker build --build-arg appName=${APP_NAME} --target bin --output bin/ .