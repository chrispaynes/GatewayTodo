.PHONY: api nukeDB proto todoApp

api:
	docker-compose up --detach --build && docker-compose logs -f api

nukeDB:
	rm -rf ./db/pgdata && docker-compose up --detach --build && docker-compose logs -f postgres

proto:
	@cd ./api && go generate ./proto/proto.go

todoApp:
	docker-compose up -d --force-recreate && docker-compose logs -f
