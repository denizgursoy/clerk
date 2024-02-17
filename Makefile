
generate:
	protoc -I=internal/server --go_out=. internal/server/server.proto