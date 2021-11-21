package gosqlparser

import "fmt"

type parser struct {
	lexer *lexer

	statement Statement
	err       error
}

// ColumnType for predefined column types.
type ColumnType int

const (
	TypeInteger ColumnType = iota
	TypeString
)

// Statement represents parsed SQL statement. Can be one of
// Select, Insert, Update, Delete, CreateTable or DropTable.
type Statement interface {
	i()
}

func (s *Select) i()      {}
func (s *Insert) i()      {}
func (s *Update) i()      {}
func (s *Delete) i()      {}
func (s *CreateTable) i() {}
func (s *DropTable) i()   {}

// Insert represents INSERT query.
//
//	INSERT [INTO] Table
//		(Columns[0], Columns[1], ...Columns[n])
//		VALUES (Values[0], Values[1], ...Values[n])
type Insert struct {
	Table   string
	Columns []string
	Values  []string
}

// Update represents UPDATE query.
//
//	UPDATE Table
//	SET
//		Columns[0] = Values[0],
//		Columns[1] = Values[1],
//		Columns[n] = Values[n]
//	WHERE
//		...
type Update struct {
	Table   string
	Columns []string
	Values  []string
	Where   Where
}

// Delete represents DELETE query.
//
//	DELETE FROM Table
//	WHERE
//		...
type Delete struct {
	Table string
	Where Where
}

// CreateTable represents CREATE TABLE statement.
//
//
type CreateTable struct {
	Name    string
	Columns []ColumnDefinition
}

// ColumnDefinition
type ColumnDefinition struct {
	Name string
	Type ColumnType
}

// CreateTable represents DROP TABLE statement.
//
// 	DROP TABLE Table
type DropTable struct {
	Table string
}

// Select represents parsed SELECT SQL statement.
//
//
type Select struct {
	Table   string
	Columns []string
	Where   Where
	Limit   int
}

// Where represent conditional expressions.
type Where struct {
	Expr Expr
}

// Expr represents expression that can be used in WHERE statement.
type Expr struct {
}

type parseFunc func(*parser) parseFunc

// Parse parses the input and returns the statement object or an error.
func Parse(input string) (Statement, error) {
	p := &parser{newLexer(input), nil, nil}

	go p.lexer.run()
	p.run()

	return p.statement, p.err
}

func (p *parser) run() {
	for state := parseStatement; state != nil; {
		state = state(p)
	}

	p.lexer.drain()
}

func (p *parser) next(skipSpace bool) token {
	for {
		t := p.lexer.nextToken()
		if !(skipSpace && t.tokenType == tokenSpace) {
			return t
		}
	}
}

func (p *parser) errorf(format string, args ...interface{}) parseFunc {
	p.statement = nil
	p.err = fmt.Errorf(format, args...)

	// TODO: add token position

	return nil
}

func (p *parser) asSelect() *Select {
	return (p.statement).(*Select)
}

func (p *parser) asInsert() *Insert {
	return (p.statement).(*Insert)
}

func (p *parser) asUpdate() *Update {
	return (p.statement).(*Update)
}

func (p *parser) asDelete() *Delete {
	return (p.statement).(*Delete)
}

func (p *parser) asCreateTable() *CreateTable {
	return (p.statement).(*CreateTable)
}

func (p *parser) asDropTable() *DropTable {
	return (p.statement).(*DropTable)
}

func parseStatement(p *parser) parseFunc {
	t := p.next(true)

	switch t.tokenType {
	case tokenError:
		return p.errorf(t.value)
	case tokenSelect:
		p.statement = &Select{}

		return parseSelect
	case tokenInsert:
		p.statement = &Insert{}

		return parseInsert
	case tokenUpdate:
		p.statement = &Update{}

		return parseUpdate
	case tokenDelete:
		p.statement = &Delete{}

		return parseDelete
	case tokenCreate:
		p.statement = &CreateTable{}

		return parseCreateTable
	case tokenDrop:
		p.statement = &DropTable{}

		return parseDropTable
	default:
		return p.errorf(
			"expected %s, %s, %s, %s, %s or %s, but got %s: %s",
			tokenSelect,
			tokenInsert,
			tokenUpdate,
			tokenDelete,
			tokenCreate,
			tokenDrop,
			t.tokenType,
			t.value,
		)
	}
}

// parseSelect initiates SELECT statement parsing
func parseSelect(p *parser) parseFunc {
	t := p.next(true)
	if t.tokenType == tokenError {
		return p.errorf(t.value)
	}

	if t.tokenType == tokenIdentifier {
		s := p.asSelect()
		s.Columns = append(s.Columns, t.value)

		// FIND DELIMITER
	}

	// TODO continue

	return nil
}

// parseDelete parses INSERT statement.
func parseInsert(p *parser) parseFunc {
	return nil
}

// parseDelete parses UPDATE statement.
func parseUpdate(p *parser) parseFunc {
	return nil
}

// parseDelete parses DELETE statement.
func parseDelete(p *parser) parseFunc {
	return nil
}

// parseCreateTable parses CREATE TABLE statement.
func parseCreateTable(p *parser) parseFunc {
	return nil
}

// parseDropTable parses DROP TABLE statement.
func parseDropTable(p *parser) parseFunc {
	t := p.next(true)
	if t.tokenType != tokenTable {
		return p.errorf("expected %s, but got %s: %s", tokenTable, t.tokenType, t.value)
	}

	t = p.next(true)
	if t.tokenType != tokenIdentifier {
		return p.errorf("expected %s, but got %s: %s", tokenIdentifier, t.tokenType, t.value)
	}

	p.asDropTable().Table = t.value

	return nil
}
