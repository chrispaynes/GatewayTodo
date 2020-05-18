.PHONY: api apiTidy apiVendor nukeDB proto todoApp uiNG uiProduction unbindPort4200

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

uiNG: unbindPort4200
	# start the UI for development
	# note: the ng serve command requires having the Angular CLI installed https://cli.angular.io/
	@cd ui && ng serve

uiProduction:
	# note: the ng serve command requires having the Angular CLI installed https://cli.angular.io/
	@cd ui && ng build --prod=true
	@cd ui && ng build --prod=true --outputPath=../api/dist/todo-app

unbindPort4200:
	kill -9 $$(lsof -i :4200 | grep node | awk '{ print $$2}' | xargs) || true
