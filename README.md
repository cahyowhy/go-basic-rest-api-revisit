# GO BASIC REST API REVISIT
## WIP

Basic Rest API Golang

Package | Description 
--- | --- 
[Gorilla](https://github.com/gorilla/mux) | Request router and dispatcher for matching incoming requests to their respective handler
[Gorm](https://github.com/go-gorm/gorm) | ORM
[GoDotEnv](https://github.com/joho/godotenv) | .env lib 
[Faker](github.com/jaswdr/faker) | faker (for seeding)
[Uuid](github.com/gofrs/uuid) | generate uuid
[JWT GO](github.com/dgrijalva/jwt-go) | json web token
[Validator](github.com/go-playground/validator) | json validator



## How to run on your local
1. copy `example.env` into `.env`, Fill with yours value
2. create your database name from `.env` `DB` value
3. add type user_role on your db
    ```sql
      CREATE TYPE user_role AS ENUM ( 'ADMIN', 'USER');
    ```
4. install [air](https://github.com/cosmtrek/air) to achieve live reload
5. or if you won't use air. you can run with `go run .`



## How to db:seed
1. run
    ```bash
    go run cmd/seeder/main.go 
    ```