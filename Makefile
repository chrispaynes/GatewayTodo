.PHONY: todoApp proto

todoApp:
	docker-compose up -d --force-recreate && docker-compose logs -f

proto:
	@cd ./api && go generate ./proto/proto.go
