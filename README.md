# GO BASIC REST API REVISIT
## WIP

Basic Rest API Golang

Package | Description 
--- | --- 
[Gorilla](https://github.com/gorilla/mux) | Request router and dispatcher for matching incoming requests to their respective handler
[Gorm](https://github.com/go-gorm/gorm) | ORM
[GoDotEnv](https://github.com/joho/godotenv) | .env lib 
[Faker](https://github.com/jaswdr/faker) | faker (for seeding)
[Uuid](https://github.com/gofrs/uuid) | generate uuid
[JWT GO](https://github.com/dgrijalva/jwt-go) | json web token
[Validator](https://github.com/go-playground/validator) | json validator



## How to run on your local
1. copy `example.env` into `.env`, Fill with yours value
2. create your database name from `.env` `DB` value
3. add type user_role on your db
    ```sql
      CREATE TYPE user_role AS ENUM ( 'ADMIN', 'USER');
    ```
4. install [air](https://github.com/cosmtrek/air) to achieve live reload
5. or if you won't use air. you can run with `go run .`



## How to db:?
1. run seeder
    ```bash
    go run cmd/seeder/main.go 
    ```
2. run migrate
    ```bash
    go run cmd/migrate/main.go 
    ```


## Folder Structure
    .
    └── go-basic-rest-api-revisit/
        ├── cmd/
        │   └── seeder        #db seeder
        │   └── migrate       #db migrate
        ├── config            #.env var value, json validator
        ├── fake              #data faker
        ├── database          #get database instance
        ├── handler           #rest api handler / controller
        ├── service           #rest api service [for bussiness logic]
        ├── middleware        #rest api middleware, before visiting controller
        ├── router            #rest api router
        ├── model             #rest api model
        └── util              #buch off helper function