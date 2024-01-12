package main

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/ghibranalj/sqauto"
)

type User struct {
	ID   int    `sq:"id"`
	Name string `sq:"name"`
	SomethingID int        `sq:"something_id"`
	Something   *Something `sq:"something"`
	CarID int        `sq:"car_id"`
	Car   *Car `sq:"car"`
}

type Something struct {
	ID   int    `sq:"id"`
	Name string `sq:"name"`
}

type Car struct {
	ID   int    `sq:"id"`
	Name string `sq:"name"`
}

func main() {

	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, _ := sqauto.SelectJoin(b, "user", User{}, "something", "car")

	fmt.Println(query)
	fmt.Println(args)
}
