# Project Variables
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourdbpassword
DB_NAME=employeside
GEN_PATH=./generated
SERVER_PATH=./cmd/server/
CONFIG_FILE=config.yaml

# Remove generated Jet code
clean:
	rm -rf ./generated

# Generate Jet code from MySQL database
codegen:
	jet -source=MySQL -host=$(DB_HOST) -port=$(DB_PORT) -user=$(DB_USER) -password=$(DB_PASSWORD) -dbname=$(DB_NAME) -path=$(GEN_PATH)

# Build the project
build:
	go build $(SERVER_PATH)

# Run the projectne
run:
	go run $(SERVER_PATH) -cfg=$(CONFIG_FILE)

# Help command to display available options
help:
	@echo "Available commands:"
	@echo "  make clean       - Remove generated Jet code"
	@echo "  make codegen    - Generate Jet code from MySQL"
	@echo "  make build       - Build the project"
	@echo "  make run         - Run the project"
