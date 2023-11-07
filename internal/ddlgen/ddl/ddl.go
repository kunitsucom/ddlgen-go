package ddl

import (
	"context"
	"time"

	"github.com/kunitsucom/ddlgen/internal/contexts"
)

type DDL struct {
	Indent string
	Header []string
	Stmts  []Stmt
}

func NewDDL(ctx context.Context) *DDL {
	return &DDL{
		Indent: "    ",
		Header: []string{
			"Code generated by ddlgen. DO NOT EDIT.",
			"",
			"Date: " + contexts.Now(ctx).Format(time.RFC3339),
			"",
		},
	}
}

type Stmt interface {
	GetSourceFile() string
	GetSourceLine() int

	private()
}
