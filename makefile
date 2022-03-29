models:
	#java -jar third_party/openapi/swagger-codegen-cli-3.0.20.jar generate -i swagger.yaml -l go -D models --model-package models -o pkg/swagger
	#java -jar third_party/openapi/swagger-codegen-cli-3.0.20.jar generate -i swagger.yaml -l typescript-angular --additional-properties modelPropertyNaming=snake_case -D models -o client/src/swagger
	java -jar third_party/openapi/swagger-codegen-cli-3.0.20.jar generate -i swagger.yaml -l kotlin-server -D models --model-package come.example.rewinder -o rewinder/app/src/main/java/com/example/rewinder/swagger
	#cp rewinder/app/src/main/java/com/example/rewinder/swagger/src/main/kotlin/io/swagger/server/models/* rewinder/app/src/main/java/com/example/rewinder/swagger/
	#rm -r rewinder/app/src/main/java/com/example/rewinder/swagger/src
