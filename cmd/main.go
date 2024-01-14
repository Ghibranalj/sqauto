package main

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/ghibranalj/sqauto"
)

type Client struct {
	ID        int
	Username  string
	PartnerID int
	Partner Partner
}

type Partner struct {
	ID   int
	Name string
	CarID int

	Car Car
	Clients []Client
}

type Car struct {
	ID   int
	Name string
}

func main() {
	sb := sqauto.Join(sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		sqauto.Table{Name: "client", Object: Client{}},
		sqauto.JoinTable{Name: "partner", Object: Partner{}, Alias: "partner"})

	query, _, _ := sb.ToSql()
	fmt.Println(query)

	query, _, _ = sqauto.Insert(sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		sqauto.Table{Name: "client", Object: Client{ID: 1, Username: "ghibran", PartnerID: 1}},
	).ToSql()
	fmt.Println(query)

	query, _,_ = sqauto.CoalesceUpdate(sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		sqauto.Table{Name: "client", Object: Client{Username: "ghibran", PartnerID: 1}},
	).ToSql()
	fmt.Println(query)

}
