PROTO_DIR := delivery/proto/domain
GEN_DIR := delivery/proto/gen

.PHONY: gen-proto
gen-proto: clean
	@echo "Generating protobuf files..."
	@find $(PROTO_DIR) -name "*.proto" -exec protoc \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		-I=$(PROTO_DIR) \
		{} \;
	@echo "Protobuf files generated successfully!"

.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	@rm -rf $(GEN_DIR)/*
	@echo "Clean complete!"
