rm -rf cover.out && rm -rf cover.html && go run cmd/migrate/main.go .test.env && go test -coverprofile cover.out -v ./... && go tool cover -html=cover.out