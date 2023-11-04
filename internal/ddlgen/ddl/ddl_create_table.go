package ddl

import (
	"regexp"

	"github.com/kunitsucom/ddlgen/internal/ddlgen/lang/util"
)

var _ Stmt = (*CreateTableStmt)(nil)

type CreateTableStmt struct {
	SourceFile  string
	SourceLine  int
	Comments    []string                 // -- <Comment>
	CreateTable string                   // CREATE TABLE [IF NOT EXISTS] <Table>
	Columns     []*CreateTableColumn     // ( <Column>, ...
	Constraints []*CreateTableConstraint // <Constraint> )
	Options     []*CreateTableOption     // <Options>;
}

func (stmt *CreateTableStmt) GetSourceFile() string {
	return stmt.SourceFile
}

func (stmt *CreateTableStmt) GetSourceLine() int {
	return stmt.SourceLine
}

func (*CreateTableStmt) private() {}

//nolint:gochecknoglobals
var stmtRegexCreateTable = &util.StmtRegex{
	Regex: regexp.MustCompile(`\s*CREATE\s+TABLE\s+(IF\s+NOT\s+EXISTS\s+)?(\S+)`),
	Index: 2,
}

func (stmt *CreateTableStmt) SetCreateTable(createTable string) {
	if len(stmtRegexCreateTable.Regex.FindStringSubmatch(createTable)) > stmtRegexCreateTable.Index {
		stmt.CreateTable = createTable
		return
	}

	stmt.CreateTable = "CREATE TABLE " + createTable
}

type CreateTableColumn struct {
	Comments       []string
	Column         string
	TypeConstraint string
}

type CreateTableConstraint struct {
	Comments   []string
	Constraint string
}

type CreateTableOption struct {
	Comments []string
	Option   string
}