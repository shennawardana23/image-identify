# image-identify

```mermaid
graph TD
    A[HTTP Request] --> B[Gin Router /api/image]
    B --> C[Database Connection]
    C --> D[Fetch Image URLs]
    D --> E[Goroutine Pool]
    E --> F1[Check URL 1]
    E --> F2[Check URL 2]
    E --> F3[Check URL n]
    F1 --> G[Atomic Counter]
    F2 --> G
    F3 --> G
    G --> H[Generate CSV Response]
```

## Technical Design

Components:
Gin HTTP Server
MySQL Database
Worker Pool for URL checking
Atomic Counter for statistics
CSV Response Generator

### Required Packages

> go get -u github.com/gin-gonic/gin
> go get -u gorm.io/gorm
> go get -u gorm.io/driver/mysql
> go get -u github.com/joho/godotenv

### Project Structure

```main
├── main.go
├── .env
├── config/
│   └── database.go
├── models/
│   └── website.go
├── handlers/
│   └── image_handler.go
├── services/
│   └── url_checker.go
└── utils/
    └── csv_generator.go
```

### Environments

```env
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=image_checker
DB_PORT=3306
SERVER_PORT=8080
WORKER_POOL_SIZE=10
```
