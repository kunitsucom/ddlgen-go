package ddlgengo

import (
	"context"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	errorz "github.com/kunitsucom/util.go/errors"

	"github.com/kunitsucom/ddlgen/internal/config"
	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
	"github.com/kunitsucom/ddlgen/internal/ddlgen/lang/util"
	"github.com/kunitsucom/ddlgen/internal/logs"
	apperr "github.com/kunitsucom/ddlgen/pkg/errors"
)

//nolint:cyclop
func Parse(ctx context.Context, src string) (*ddlast.DDL, error) {
	// MEMO: get absolute path for parser.ParseFile()
	sourceAbs, err := filepath.Abs(src)
	if err != nil {
		return nil, errorz.Errorf("filepath.Abs: %w", err)
	}

	info, err := os.Stat(sourceAbs)
	if err != nil {
		return nil, errorz.Errorf("os.Stat: %w", err)
	}

	ddl := ddlast.NewDDL(ctx)

	if info.IsDir() {
		if err := filepath.WalkDir(sourceAbs, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err //nolint:wrapcheck
			}

			if d.IsDir() {
				return nil
			}

			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			stmts, err := parseFile(ctx, path)
			if err != nil {
				if errors.Is(err, apperr.ErrDDLKeyGoNotFoundInSource) {
					logs.Debug.Printf("parseFile: %s: %v", path, err)
					return nil
				}
				return errorz.Errorf("parseFile: %w", err)
			}

			ddl.Stmts = append(ddl.Stmts, stmts...)

			return nil
		}); err != nil {
			return nil, errorz.Errorf("filepath.WalkDir: %w", err)
		}

		return ddl, nil
	}

	stmts, err := parseFile(ctx, sourceAbs)
	if err != nil {
		return nil, errorz.Errorf("Parse: %w", err)
	}
	ddl.Stmts = append(ddl.Stmts, stmts...)

	return ddl, nil
}

//nolint:cyclop,funlen,gocognit
func parseFile(ctx context.Context, filename string) ([]ddlast.Stmt, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, errorz.Errorf("parser.ParseFile: %w", err)
	}

	ddlSrc, err := extractDDLSource(ctx, fset, f)
	if err != nil {
		return nil, errorz.Errorf("extractDDLSource: %w", err)
	}

	dumpDDLSource(fset, ddlSrc)

	stmts := make([]ddlast.Stmt, 0)
	for _, r := range ddlSrc {
		stmt := &ddlast.CreateTableStmt{}

		// source
		source := fset.Position(r.StructType.Pos())
		stmt.SourceFile = source.Filename
		stmt.SourceLine = source.Line

		// comments
		comments := strings.Split(strings.Trim(r.CommentGroup.Text(), "\n"), "\n")
		for _, comment := range comments {
			logs.Debug.Printf("[COMMENT DETECTED]: %s:%d: %s", stmt.SourceFile, stmt.SourceLine, comment)
		}
		// stmt.Comments = append(stmt.Comments, util.TrimTailEmptyCommentElement(util.TrimDDLGenCommentElement(comments))...)
		stmt.Comments = append(stmt.Comments, comments...)

		// CREATE TABLE / CONSTRAINT / OPTIONS
		for _, comment := range comments {
			if matches := util.StmtRegexCreateTable.Regex.FindStringSubmatch(comment); len(matches) > util.StmtRegexCreateTable.Index {
				stmt.SetCreateTable(matches[util.StmtRegexCreateTable.Index])
			} else if matches := util.StmtRegexCreateTableConstraint.Regex.FindStringSubmatch(comment); len(matches) > util.StmtRegexCreateTableConstraint.Index {
				stmt.Constraints = append(stmt.Constraints, &ddlast.CreateTableConstraint{
					Constraint: matches[util.StmtRegexCreateTableConstraint.Index],
				})
			} else if matches := util.StmtRegexCreateTableOptions.Regex.FindStringSubmatch(comment); len(matches) > util.StmtRegexCreateTableOptions.Index {
				stmt.Options = append(stmt.Options, &ddlast.CreateTableOption{
					Option: matches[util.StmtRegexCreateTableOptions.Index],
				})
			}
		}
		if stmt.CreateTable == "" {
			stmt.SetCreateTable(r.TypeSpec.Name.String())
		}

		// columns
		for _, field := range r.StructType.Fields.List {
			column := &ddlast.CreateTableColumn{}

			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			// column name
			switch columnName := tag.Get(config.ColumnKeyGo()); columnName {
			case "-":
				logs.Info.Printf("[%s]: Ignore columns with column name \"-\": %s", stmt.CreateTable, field.Names[0].Name)
				continue
			case "":
				column.Column = field.Names[0].Name
			default:
				column.Column = columnName
			}

			// column type and constraint
			switch columnTypeConstraint := tag.Get(config.DDLKeyGo()); columnTypeConstraint {
			case "":
				logs.Warn.Printf("[%s]: Ignore columns with no type and constraints set: %s", stmt.CreateTable, field.Names[0].Name)
				continue
			default:
				column.TypeConstraint = columnTypeConstraint
			}

			// comments
			comments := strings.Split(strings.Trim(field.Doc.Text(), "\n"), "\n")
			column.Comments = append(column.Comments, util.TrimTailEmptyCommentElement(util.TrimDDLGenCommentElement(comments))...)

			stmt.Columns = append(stmt.Columns, column)
		}

		stmts = append(stmts, stmt)
	}

	sort.Slice(stmts, func(i, j int) bool {
		return fmt.Sprintf("%s:%09d", stmts[i].GetSourceFile(), stmts[i].GetSourceLine()) < fmt.Sprintf("%s:%09d", stmts[j].GetSourceFile(), stmts[j].GetSourceLine())
	})

	for i := range stmts {
		logs.Trace.Print(fmt.Sprintf("%s:%09d", stmts[i].GetSourceFile(), stmts[i].GetSourceLine()))
	}

	return stmts, nil
}
