models:
	java -jar third_party/openapi/swagger-codegen-cli-3.0.20.jar generate -i swagger.yaml -l go -D models -o pkg/swagger
	java -jar third_party/openapi/swagger-codegen-cli-3.0.20.jar generate -i swagger.yaml -l typescript-angular --additional-properties modelPropertyNaming=snake_case -D models -o client/src/models
	java -jar third_party/openapi/swagger-codegen-cli-3.0.20.jar generate -i swagger.yaml -l kotlin-client -D models -o pkg/models
