# Employee-Side Service

This project is a backend service for Employee-Side, built using Golang and MySQL.

## Installation & Setup

### **Prerequisites**
- Go 1.22.4+
- MySQL database
- Jet (for SQL code generation)
- Git


### **Clone the Repository**
```sh
git clone https://github.com/Employes-Side/Employeside-service.git
cd Employeside-service
```

### **Configuration Variables**
Create a `config.yaml` file and configure the required environment variables:
```env
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=

http port
```

### **Using Makefile Commands**

| Command | Description |
|---------|-------------|
| `make clean` | Remove generated Jet code |
| `make generate` | Generate Jet code from MySQL |
| `make build` | Build the project binary |
| `make run` | Run the project with `config.yaml` |

### **Manual Commands (Without Makefile)**
```sh
# Remove generated Jet code
rm -rf ./generated

# Generate Jet code
jet -source=MySQL -host=localhost -port=3306 -user=root -password=yourpasswd -dbname=users -path=./internal

install jet before generating above command -: below the installing jet
go get -u github.com/go-jet/jet/v2/cmd/jet
export PATH=$PATH:$HOME/go/bin
source ~/.zshrc


# Build the project
go build -o employeside ./cmd/server/

# Run the project
./employeside -cfg=config.yaml
```


