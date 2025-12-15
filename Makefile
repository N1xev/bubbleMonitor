PROJECT_NAME := bub
VERSION := 0.1.0
BUILD_DIR := build
COMPRESS ?= false

GOCMD := go
GOBUILD := $(GOCMD) build
GOMOD := $(GOCMD) mod
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test

TARGETS := windows/amd64 windows/arm64 linux/amd64 linux/arm linux/arm64 darwin/amd64 darwin/arm64 android/arm64

ifneq ($(TERM),)
	GREEN := \033[32m
	YELLOW := \033[33m
	RED := \033[31m
	BLUE := \033[34m
	RESET := \033[0m
endif

.PHONY: all build clean test help deps release archive

all: deps build

build:
	@echo "$(BLUE)================================================$(RESET)"
	@echo "$(BLUE)   Go Multi-Platform Build System$(RESET)"
	@echo "$(BLUE)   Building: $(PROJECT_NAME) v$(VERSION)$(RESET)"
	@echo "$(BLUE)================================================$(RESET)"
	@echo ""
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(PROJECT_NAME)-$(shell uname | tr '[:upper:]' '[:lower:]')-$(shell uname -m)-v$(VERSION) main.go
	@echo "$(GREEN)[SUCCESS]$(RESET) Built for current platform"

build-all: deps clean
	@echo "$(BLUE)================================================$(RESET)"
	@echo "$(BLUE)   Go Multi-Platform Build System$(RESET)"
	@echo "$(BLUE)   Building: $(PROJECT_NAME) v$(VERSION)$(RESET)"
	@echo "$(BLUE)================================================$(RESET)"
	@echo ""
	@mkdir -p $(BUILD_DIR)
	@$(eval SUCCESS_COUNT := 0)
	@$(eval FAIL_COUNT := 0)
	@$(foreach target,$(TARGETS),$(call build_target,$(target)))
	@echo "$(BLUE)------------------------------------------------$(RESET)"
	@echo ""
	@echo "$(BLUE)Build Summary:$(RESET)"
	@echo "$(BLUE)================================================$(RESET)"
	@echo "$(GREEN)[SUCCESS] Successful: $(SUCCESS_COUNT) / $(words $(TARGETS))$(RESET)"
	@if [ $(FAIL_COUNT) -gt 0 ]; then \
		echo "$(RED)[FAILED] Failed: $(FAIL_COUNT) / $(words $(TARGETS))$(RESET)"; \
	fi
	@echo ""
	@$(eval TOTAL_SIZE := $(shell du -sh $(BUILD_DIR) 2>/dev/null | cut -f1 || echo "0"))
	@echo "$(BLUE)[INFO] Total size: $(TOTAL_SIZE)$(RESET)"
	@echo "$(BLUE)[INFO] Output directory: $(BUILD_DIR)$(RESET)"
	@echo ""
	@echo "$(BLUE)[INFO] Built binaries:$(RESET)"
	@ls -lah $(BUILD_DIR)/$(PROJECT_NAME)-* 2>/dev/null || echo "   No binaries found"
	@echo ""
	@if [ $(SUCCESS_COUNT) -eq $(words $(TARGETS)) ]; then \
		echo "$(GREEN)[SUCCESS] ALL BUILDS COMPLETED SUCCESSFULLY!$(RESET)"; \
	else \
		if [ $(SUCCESS_COUNT) -gt 0 ]; then \
			echo "$(YELLOW)[WARNING] Some builds failed, but $(SUCCESS_COUNT) succeeded$(RESET)"; \
		else \
			echo "$(RED)[ERROR] ALL BUILDS FAILED!$(RESET)"; \
			exit 1; \
		fi; \
	fi

define build_target
	@$(eval GOOS := $(word 1,$(subst /, ,$(1))))
	@$(eval GOARCH := $(word 2,$(subst /, ,$(1))))
	@$(eval EXT := $(if $(filter windows,$(GOOS)),.exe,))
	@$(eval OUTPUT_NAME := $(PROJECT_NAME)-$(GOOS)-$(GOARCH)-v$(VERSION)$(EXT))
	@$(eval OUTPUT_PATH := $(BUILD_DIR)/$(OUTPUT_NAME))
	@echo "  $(BLUE)[BUILD]$(RESET) $(GOOS)/$(GOARCH) -> $(OUTPUT_NAME)"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -ldflags="-s -w -X main.Version=$(VERSION)" -o $(OUTPUT_PATH) main.go 2>/dev/null
	@if [ $$? -eq 0 ] && [ -f "$(OUTPUT_PATH)" ]; then \
		$(eval SUCCESS_COUNT := $(shell echo $$(($(SUCCESS_COUNT) + 1)))) \
		@$(eval FILE_SIZE := $(shell stat -f%z "$(OUTPUT_PATH)" 2>/dev/null || stat -c%s "$(OUTPUT_PATH)" 2>/dev/null || echo "0")) \
		@if [ $$FILE_SIZE -gt 1048576 ]; then \
			echo "  $(GREEN)[OK]$(RESET) Built successfully ($$(echo "scale=1; $$FILE_SIZE / 1048576" | bc) MB)"; \
		elif [ $$FILE_SIZE -gt 1024 ]; then \
			echo "  $(GREEN)[OK]$(RESET) Built successfully ($$(echo "scale=1; $$FILE_SIZE / 1024" | bc) KB)"; \
		else \
			echo "  $(GREEN)[OK]$(RESET) Built successfully ($$FILE_SIZE bytes)"; \
		fi \
		@if [ "$(COMPRESS)" = "true" ]; then \
			if command -v gzip >/dev/null 2>&1; then \
				gzip -k $(OUTPUT_PATH); \
				echo "        Compressed: $(OUTPUT_NAME).gz"; \
			elif command -v zip >/dev/null 2>&1; then \
				zip -j $(BUILD_DIR)/$(OUTPUT_NAME).zip $(OUTPUT_PATH) >/dev/null 2>&1; \
				echo "        Compressed: $(OUTPUT_NAME).zip"; \
			fi; \
		fi \
	else \
		$(eval FAIL_COUNT := $(shell echo $$(($(FAIL_COUNT) + 1)))) \
		echo "  $(RED)[FAIL]$(RESET) Build failed"; \
	fi
	@echo ""
endef

deps:
	@echo "$(BLUE)[INFO]$(RESET) Downloading dependencies..."
	@$(GOMOD) download
	@echo "$(GREEN)[SUCCESS]$(RESET) Dependencies downloaded"
	@echo ""

clean:
	@echo "$(YELLOW)[WARNING]$(RESET) Cleaning existing build directory..."
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)
	@echo "$(GREEN)[SUCCESS]$(RESET) Created build directory"
	@echo ""

test:
	@echo "$(BLUE)[INFO]$(RESET) Running tests..."
	@$(GOTEST) ./...
	@echo "$(GREEN)[SUCCESS]$(RESET) Tests completed"
release: build-all
	@echo ""
	@$(eval RELEASE_NAME := $(PROJECT_NAME)-v$(VERSION)-all-platforms.zip)
	@echo "$(BLUE)[INFO]$(RESET) Creating release archive: $(RELEASE_NAME)"
	@if command -v zip >/dev/null 2>&1; then \
		cd $(BUILD_DIR) && zip -r ../$(RELEASE_NAME) $(PROJECT_NAME)-*; \
		echo "$(GREEN)[SUCCESS]$(RESET) Release archive created: $(RELEASE_NAME)"; \
	elif command -v tar >/dev/null 2>&1; then \
		tar -czf $(RELEASE_NAME) -C $(BUILD_DIR) $(PROJECT_NAME)-*; \
		echo "$(GREEN)[SUCCESS]$(RESET) Release archive created: $(RELEASE_NAME)"; \
	else \
		echo "$(RED)[ERROR]$(RESET) No archive tool available (zip/tar)"; \
	fi

archive: build-all
	@echo ""
	@$(eval ARCHIVE_NAME := $(PROJECT_NAME)-v$(VERSION)-all-platforms.tar.gz)
	@echo "$(BLUE)[INFO]$(RESET) Creating compressed archive: $(ARCHIVE_NAME)"
	@if command -v tar >/dev/null 2>&1; then \
		tar -czf $(ARCHIVE_NAME) -C $(BUILD_DIR) $(PROJECT_NAME)-*; \
		$(eval ARCHIVE_SIZE := $(shell du -sh $(ARCHIVE_NAME) | cut -f1)) \
		echo "$(GREEN)[SUCCESS]$(RESET) Archive created: $(ARCHIVE_NAME) ($(ARCHIVE_SIZE))"; \
	else \
		echo "$(RED)[ERROR]$(RESET) tar not available"; \
	fi

install: deps build
	@echo "$(GREEN)[SUCCESS]$(RESET) Installation completed"
debug:
	@echo "$(BLUE)[INFO]$(RESET) Building with debug information..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(PROJECT_NAME)-debug main.go
	@echo "$(GREEN)[SUCCESS]$(RESET) Debug build completed"

help:
	@echo "$(BLUE)Available targets:$(RESET)"
	@echo ""
	@echo "$(GREEN)Build targets:$(RESET)"
	@echo "  $(YELLOW)make$(RESET) or $(YELLOW)make all$(RESET)     - Download dependencies and build for current platform"
	@echo "  $(YELLOW)make build$(RESET)            - Build for current platform only"
	@echo "  $(YELLOW)make build-all$(RESET)        - Build for all platforms"
	@echo "  $(YELLOW)make debug$(RESET)            - Build with debug information"
	@echo ""
	@echo "$(GREEN)Utility targets:$(RESET)"
	@echo "  $(YELLOW)make deps$(RESET)             - Download dependencies"
	@echo "  $(YELLOW)make clean$(RESET)            - Clean build directory"
	@echo "  $(YELLOW)make test$(RESET)             - Run tests"
	@echo "  $(YELLOW)make release$(RESET)          - Build all platforms and create ZIP archive"
	@echo "  $(YELLOW)make archive$(RESET)          - Build all platforms and create TAR.GZ archive"
	@echo "  $(YELLOW)make install$(RESET)          - Install dependencies and build"
	@echo ""
	@echo "$(GREEN)Information:$(RESET)"
	@echo "  $(YELLOW)make help$(RESET)             - Show this help message"
	@echo ""
	@echo "$(BLUE)Variables:$(RESET)"
	@echo "  $(YELLOW)PROJECT_NAME=$(PROJECT_NAME)$(RESET) - Set the project name"
	@echo "  $(YELLOW)VERSION=$(VERSION)$(RESET)          - Set the version"
	@echo "  $(YELLOW)COMPRESS=true$(RESET)        - Enable compression of individual binaries"
	@echo ""
	@echo "$(BLUE)Examples:$(RESET)"
	@echo "  $(YELLOW)make PROJECT_NAME=myapp VERSION=2.0.0 build-all$(RESET)"
	@echo "  $(YELLOW)make COMPRESS=true build-all$(RESET)"

info:
	@echo "$(BLUE)Current Configuration:$(RESET)"
	@echo "  Project Name: $(PROJECT_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Build Directory: $(BUILD_DIR)"
	@echo "  Go Version: $(shell $(GOCMD) version 2>/dev/null || echo 'Not installed')"
	@echo "  Targets: $(TARGETS)"
