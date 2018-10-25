SERVER_PATH=./cmd/db


run:
	@rm -f $(SERVER_PATH)/main && go build -o $(SERVER_PATH)/db $(SERVER_PATH)/main.go && $(SERVER_PATH)/db