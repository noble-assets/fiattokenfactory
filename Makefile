.PHONY: proto-format proto-lint proto-gen format build local-image test-unit test-e2e
all: proto-all format test-unit build local-image test-e2e

###############################################################################
###                                  Build                                  ###
###############################################################################

build:
	@echo "🤖 Building simd..."
	@cd simapp && make build
	@echo "✅ Completed build!"

###############################################################################
###                          Formatting & Linting                           ###
###############################################################################

gofumpt_cmd=mvdan.cc/gofumpt

format:
	@echo "🤖 Running formatter..."
	@go run $(gofumpt_cmd) -l -w .
	@echo "✅ Completed formatting!"

###############################################################################
###                                Protobuf                                 ###
###############################################################################

BUF_VERSION=1.30.1
BUILDER_VERSION=0.14.0

proto-all: proto-format proto-lint proto-gen

proto-format:
	@echo "🤖 Running protobuf formatter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) format --diff --write
	@echo "✅ Completed protobuf formatting!"

proto-gen:
	@echo "🤖 Generating code from protobuf..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		ghcr.io/cosmos/proto-builder:$(BUILDER_VERSION) sh ./proto/generate.sh
	@echo "✅ Completed code generation!"

proto-lint:
	@echo "🤖 Running protobuf linter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) lint
	@echo "✅ Completed protobuf linting!"

###############################################################################
###                                 Testing                                 ###
###############################################################################

local-image:
	@echo "🤖 Building image..."
	@heighliner build --chain noble-fiattokenfactory-simd --local 1> /dev/null
	@echo "✅ Completed build!"

test-e2e:
	@echo "🤖 Running e2e tests..."
	@cd e2e && GOWORK=off go test -timeout 0 -race -v ./...
	@echo "✅ Completed e2e tests!"

test-unit:
	@echo "🤖 Running unit tests..."
	@go test -cover -coverprofile=coverage.out -race -v ./x/fiattokenfactory/keeper/...
	@echo "✅ Completed unit tests!"
