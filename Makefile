tidy:
	go mod tidy

codegen:
	oapi-codegen \
	--config=openapi/config.yaml openapi/api.yaml