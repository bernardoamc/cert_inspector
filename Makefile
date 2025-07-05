APP_NAME=cert_inspector
CMD_DIR=cmd/$(APP_NAME)
BUILD_DIR=build
BINARY=$(BUILD_DIR)/$(APP_NAME)
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_VOLUME=cert_logs

.PHONY: all
all: build

.PHONY: test
test:
	go test ./...

## 🔨 Build Go binary locally
.PHONY: build
build:
	@echo ">> Building $(APP_NAME)..."
	mkdir -p $(BUILD_DIR)
	go build -v -o $(BINARY) ./$(CMD_DIR)

## 🧪 Run locally (needs logs dir and targets.txt)
.PHONY: run
run: build logs
	@echo ">> Running locally..."
	./$(BINARY) -binary ./gungnir -targets targets.txt -log-dir logs

## 📂 Ensure local logs dir exists & owned by container user
.PHONY: logs
logs:
	mkdir -p logs
	sudo chown -R 65532:65532 logs

## 🐳 Build Docker image
.PHONY: docker-build
docker-build:
	docker build --no-cache -t $(DOCKER_IMAGE) .

## ▶️  Run Docker with local logs folder (bind mount)
.PHONY: docker-run-local
docker-run-local: logs
	@echo ">> Running Docker container with local bind mount logs..."
	docker run --rm \
		-v $(PWD)/logs:/app/logs \
		$(DOCKER_IMAGE)

## 🐳 Create Docker volume for logs
.PHONY: volume-create
volume-create:
	docker volume create $(DOCKER_VOLUME)

## 🐳 Remove Docker volume for logs
.PHONY: volume-rm
volume-rm:
	docker volume rm $(DOCKER_VOLUME)

## 🛠 Fix permissions inside Docker volume so nonroot user can write
.PHONY: volume-chown
volume-chown:
	docker run --rm -v $(DOCKER_VOLUME):/app/logs alpine chown -R 65532:65532 /app/logs

## ▶️  Run Docker with volume (container-native)
.PHONY: docker-run-volume
docker-run-volume: volume-rm volume-create volume-chown
	@echo ">> Running Docker container with Docker volume..."
	docker run --rm \
		-v $(DOCKER_VOLUME):/app/logs \
		$(DOCKER_IMAGE)

## 📦 Copy logs from Docker volume to host ./logs_backup
.PHONY: copy-logs
copy-logs:
	mkdir -p logs_backup
	docker run --rm \
		-v $(DOCKER_VOLUME):/data \
		-v $(PWD)/logs_backup:/backup \
		alpine \
		cp -r /data/. /backup/

## 🧹 Clean local build
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

## 📦 Tidy go.mod
.PHONY: tidy
tidy:
	go mod tidy
