# A Restful API Boilerplate for Go

Clone this repo and start your API development.

## Start Application
### Docker
- Build Docker Images ```docker-compose build api```
- Run Docker Container ```docker-compose up api```

### Local
- Start dependencies ```docker-compose up -d db```
- Start API with *serve*: ```go run main.go serve```

## Database
- Integrated with PostgreSQL
- Integrated with migration package.
- Run migration: ```go run main.go migrate```
- Reset migration: ```go run main.go migrate --reset```

### Tables
- courses

| Columns       | Description               |
|---------------|---------------------------|
| id	        | Primary key of table      |
| name	        | Name of the course        |
| description   | Course description        | 
| author		| Author of description     |
| created_at 	| Created time of course    |

## Integrated Packages
- HTTP Router - https://github.com/go-chi/chi
- PostgreSQL - https://github.com/jackc/pgx
- Logging - https://github.com/sirupsen/logrus
- Cli - https://github.com/spf13/cobra
- Configuration - https://github.com/spf13/viper
