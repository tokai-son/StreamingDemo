# パス設定
PROTO_DIR := api/proto                  # .protoファイルの格納ディレクトリ
OUT_DIR := api/generated               # 生成されたコードの出力先ディレクトリ
PROTOC := protoc                       # Protocol Buffers Compiler

# .proto ファイルのリスト
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# タスク
.PHONY: all generate clean

# デフォルトタスク
all: generate

# コード生成
generate:
	@echo "Generating Go code from Proto files..."
	@mkdir -p $(OUT_DIR)  # フォルダがない場合は作成
	$(PROTOC) --proto_path=$(PROTO_DIR) \
	          --go_out=$(OUT_DIR) \
	          --go-grpc_out=$(OUT_DIR) \
	          $(PROTO_FILES)/*.proto
	@echo "Code generation completed!"

# クリーンアップ（生成されたコードを削除）
clean:
	@echo "Cleaning generated files..."
	rm -rf $(OUT_DIR)/*.pb.go
	@echo "Cleanup completed!"

run-server:
	@echo "Starting hot reload for cmd/server/main.go..."
	air -c .air.server.toml

send-sample-request:
	@echo "Sending sample request to server..."
	grpcurl -d '{"videoID": "example_video_id", "startTime": 0, "quality": "HIGH"}' -plaintext localhost:5432 streaming.VideoStreamService/StreamVideo
