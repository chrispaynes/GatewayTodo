.PHONY: api apiTidy apiVendor nukeDB proto todoApp ui

api:
	# restart the api
	@docker-compose up --detach --build && docker-compose logs -f api

apiVendor:
	@cd ./api && go mod vendor

apiTidy:
	@cd ./api && go mod tidy

nukeDB:
	@# wipe the database and remove the mounted db volume to ensure the DDL initialization scripts are re-run
	@rm -rf ./db/pgdata && docker-compose up --detach --build && docker-compose logs -f postgres

proto:
	@# compile api protobuf files
	@go generate ./api/proto/proto.go

todoApp:
	@docker-compose up -d --force-recreate && docker-compose logs -f

ui:
	# start the UI for development
	# note: the ng serve command requires having the Angular CLI installed https://cli.angular.io/
	@cd ui && ng serve
