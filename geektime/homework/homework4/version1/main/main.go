package main

import (
	"database/sql"
	"fmt"
)

/*
@@zh 2022-12-17 21.54@@
*/
type PaulsModel struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  sql.NullString
}

func main() {
	fmt.Println("uryyb")
	//fmt.Println(reflect.TypeOf(&PaulsModel{}).Name())

}
