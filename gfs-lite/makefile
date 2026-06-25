.PHONY: proto test master chunk status clean-data

proto:
	protoc \
		--go_out=. \
		--go_opt=module=github.com/dsbudziwojski/gfs-lite \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/dsbudziwojski/gfs-lite \
		api/gfs/v1/master.proto \
		api/gfs/v1/chunkserver.proto

test:
	go test ./...

master:
	go run ./cmd/gfsmaster --addr :8000

chunk:
	go run ./cmd/gfschunk --master localhost:8000 --addr localhost:8100

status:
	go run ./cmd/gfsctl --master localhost:8000 --chunk localhost:8100

clean-data:
	rm -rf data