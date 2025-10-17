# -----------------------------
# Go 应用简洁构建 Makefile
# -----------------------------

APP_NAME ?= avatar

# 默认使用当前平台
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# 默认输出文件名（当前目录下）
OUTPUT_FILE := $(APP_NAME)
ifeq ($(GOOS),windows)
OUTPUT_FILE := $(APP_NAME).exe
endif


# -----------------------------
# 目标定义
# -----------------------------

.PHONY: all build clean linux windows darwin darwin-amd darwin-arm help

## 默认目标
all: build

## 构建当前平台
build:
	@OUTPUT_FILE=$(OUTPUT_FILE); \
	echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH) → $$OUTPUT_FILE"; \
	GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -trimpath -o $$OUTPUT_FILE ./cmd/main.go; \
	echo "build Done: $$OUTPUT_FILE"

## 构建 Linux amd64
linux:
	@$(MAKE) GOOS=linux GOARCH=amd64 build

## 构建 Windows amd64
windows:
	@$(MAKE) GOOS=windows GOARCH=amd64 build

## 构建 macOS (darwin) 两个架构
darwin: darwin-amd darwin-arm

darwin-amd:
	@$(MAKE) GOOS=darwin GOARCH=amd64 build

darwin-arm:
	@$(MAKE) GOOS=darwin GOARCH=arm64 build

## 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(APP_NAME) $(APP_NAME).exe

## 查看帮助
help:
	@echo "  Makefile usage:"
	@echo "  make build              # 构建当前平台 (输出在当前目录)"
	@echo "  make linux              # 构建 Linux 版本"
	@echo "  make windows            # 构建 Windows 版本"
	@echo "  make darwin             # 构建 macOS 版本 (amd64 + arm64)"
	@echo "  make darwin-amd         # 构建 macOS amd64"
	@echo "  make darwin-arm         # 构建 macOS arm64"
	@echo "  make clean              # 清理产物"

