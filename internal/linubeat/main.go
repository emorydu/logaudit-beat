package main

import (
	"fmt"
	"github.com/emorydu/dbaudit/internal/common/client"
)

func main() {
	c, auditBeatClosed, err := client.NewAuditBeatClient("")
	if err != nil {
		panic(err)
	}
	defer auditBeatClosed()

	fmt.Println(c)
}
