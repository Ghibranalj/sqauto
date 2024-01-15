package main

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ghibranalj/sqauto"
	"github.com/guregu/null"
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
	AppointmentDate time.Time
	MaidenName null.String

	Car Car
	Clients []Client
}

type Car struct {
}

func main() {
	sb := sqauto.Join(sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		sqauto.Table{Name: "client", Object: Client{}},
		sqauto.JoinTable{Name: "partner", Object: Partner{}})

	query, _, _ := sb.ToSql()
	fmt.Println(query)
}
