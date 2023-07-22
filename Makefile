
run:
	go run ./cmd/server/

test:
	go test -race ./...

tidy:
	go mod tidy

# will format the proto file to look like go code, using the .clang-format file.
clang-format:
	clang-format -i api/v1/*.proto

proto:
	protoc api/v1/*.proto \
	--go_out=. \
	--go_opt=paths=source_relative \
	--proto_path=.