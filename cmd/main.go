package main

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/ghibranalj/sqauto"
)

type Client struct {
	ID        int    `sq:"id"`
	Username  string `sq:"username"`
	PartnerID int    `sq:"partner_id"`

	Partner Partner `sq:"partner"`
}

type Partner struct {
	ID   int    `sq:"id"`
	Name string `sq:"name"`
}

func main() {

	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, _, _ := sqauto.SelectJoin(b, "client", Client{}, "partner")

	fmt.Println(query)
}
