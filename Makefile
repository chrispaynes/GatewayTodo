.PHONY: todoApp proto

todoApp:
	docker-compose up -d --force-recreate && docker-compose logs -f


nukeDB:
	rm -rf ./db/pgdata && docker-compose up --detach --build && docker-compose logs -f postgres

proto:
	@cd ./api && go generate ./proto/proto.go
