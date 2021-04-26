package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/database"
	"gorm.io/gorm"
)

func main() {
	cf := config.GetConfig(os.Args[1:]...)
	db := database.GetDatabase(cf)

	fmt.Println(cf)

	dbDeleteAllRow(db)
}

func dbDeleteAllRow(db *gorm.DB) {
	fineChDest, ubookChDest, usessChDest := make(map[string]interface{}), make(map[string]interface{}), make(map[string]interface{})
	errFineHCh, errUbookCh, errUsessCh := deleteAll(db, "user_fine_histories", &fineChDest), deleteAll(db, "user_books", &ubookChDest), deleteAll(db, "user_sessions", &usessChDest)
	errFineH, errUbook, errUsess := <-errFineHCh, <-errUbookCh, <-errUsessCh
	errs := []error{errFineH, errUbook, errUsess}

	uChDest, bChDest := make(map[string]interface{}), make(map[string]interface{})
	errUCh, errBCh := deleteAll(db, "users", &uChDest), deleteAll(db, "books", &bChDest)
	errU, errB := <-errUCh, <-errBCh
	errs = append(errs, errU, errB)

	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
			break
		}
	}
}

func deleteAll(db *gorm.DB, tableName string, dest interface{}) <-chan error {
	r := make(chan error)

	go func() {
		defer close(r)
		r <- db.Debug().Raw(fmt.Sprintf("DELETE FROM %s WHERE id IS NOT NULL", tableName)).Scan(dest).Error
	}()

	return r
}
