models:
	java --add-opens=java.base/java.util=ALL-UNNAMED -jar third_party/openapi/swagger-codegen-cli-3.0.20.jar generate -i api/swagger.yaml -l go -D models --model-package models -o pkg/swagger
	java --add-opens=java.base/java.util=ALL-UNNAMED -jar third_party/openapi/swagger-codegen-cli-3.0.20.jar generate -i api/swagger.yaml -l typescript-angular --additional-properties modelPropertyNaming=snake_case -D models -o client/src/swagger
	java --add-opens=java.base/java.util=ALL-UNNAMED -jar third_party/openapi/swagger-codegen-cli-3.0.20.jar generate -i api/swagger.yaml -l kotlin-server --additional-properties sourceFolder="" -D models -o rewinder/app/src/main/java

.dockerignore: client/.gitignore rewinder/.gitignore .gitignore
	scripts/gen_docker_ignore.sh

.PHONY: server_image
server_image:
	docker build . -f build/server.Dockerfile -t vweb

bin/server:
	mkdir -p bin
	go build -o bin/server cmd/http_server/server.go

client/dist:
	cd client && yarn build