.PHONY: api nukeDB proto todoApp

api:
	# restart the api
	@docker-compose up --detach --build && docker-compose logs -f api

nukeDB:
	@# wipe the database and remove the mounted db volume to ensure the DDL initialization scripts are re-run
	@rm -rf ./db/pgdata && docker-compose up --detach --build && docker-compose logs -f postgres

proto:
	@# compile api protobuf files
	@go generate ./api/proto/proto.go

todoApp:
	@docker-compose up -d --force-recreate && docker-compose logs -f
