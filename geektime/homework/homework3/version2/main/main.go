package main

import (
	"exercise/geektime/homework3/version2"
	"fmt"
)

/*
zh 2022-12-01 19:21
@@
*/

type User struct {
	Name string
}

func main() {
	query, err := (&version2.Selector[User]{}).From("user").Where(version2.NewColumn("age").Eq(18).Or(version2.NewColumn("name").Eq("zh"))).Build()
	if err != nil {
		return
	}
	fmt.Println(query)
}
