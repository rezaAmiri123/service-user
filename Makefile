.PHONY: proto
proto:
	protoc -I ./proto \
		--go_out ./gen/pb --go_opt paths=source_relative \
		--go-grpc_out ./gen/pb --go-grpc_opt paths=source_relative,require_unimplemented_servers=false \
		--grpc-gateway_out ./gen/pb --grpc-gateway_opt paths=source_relative --grpc-gateway_opt logtostderr=true \
		--openapiv2_out ./gen/swagger --openapiv2_opt logtostderr=true \
		proto/*.proto
