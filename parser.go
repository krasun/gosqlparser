package gosqlparser

import (
	"fmt"
	"strings"
)

type parser struct {
	lexer *lexer

	statement Statement
	err       error
}

// Parse parses the input and returns the statement object or an error.
func Parse(input string) (Statement, error) {
	p := &parser{newLexer(input), nil, nil}

	go p.lexer.run()
	p.run()

	return p.statement, p.err
}

// ColumnType for predefined column types.
type ColumnType int

const (
	TypeInteger ColumnType = iota
	TypeString
)

type Operator int

const (
	OperatorEquals Operator = iota
	OperatorLogicalAnd
)

// Statement represents parsed SQL statement. Can be one of
// Select, Insert, Update, Delete, CreateTable or DropTable.
type Statement interface {
	i()
}

func (*Select) i()      {}
func (*Insert) i()      {}
func (*Update) i()      {}
func (*Delete) i()      {}
func (*CreateTable) i() {}
func (*DropTable) i()   {}

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
	Where   *Where
}

// Delete represents DELETE query.
//
//	DELETE FROM Table
//	WHERE
//		...
type Delete struct {
	Table string
	Where *Where
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
	Where   *Where
	Limit   string
}

// Where represent conditional expressions.
type Where struct {
	Expr Expr
}

// Expr represents expression that can be used in WHERE statement.
type Expr interface {
	i()
}

func (ExprIdentifier) i()   {}
func (ExprValueInteger) i() {}
func (ExprValueString) i()  {}
func (ExprOperation) i()    {}

type ExprIdentifier struct {
	Name string
}

type ExprValueInteger struct {
	Value string
}

type ExprValueString struct {
	Value string
}

type ExprOperation struct {
	Left     Expr
	Operator Operator
	Right    Expr
}

type parseFunc func(*parser) parseFunc

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
	// TODO: add token position

	return p.error(fmt.Errorf(format, args...))
}

func (p *parser) error(err error) parseFunc {
	p.err = err

	return nil
}

func (p *parser) scanFor(tokenTypes ...tokenType) (token, error) {
	t := p.next(true)
	if t.tokenType == tokenError {
		return token{}, fmt.Errorf(t.value)
	}

	expectedTokens := []string{}
	for _, tokenType := range tokenTypes {
		if tokenType == t.tokenType {
			return t, nil
		}

		expectedTokens = append(expectedTokens, tokenType.String())
	}

	return token{}, fmt.Errorf("expected %s, but got %s: \"%s\"", strings.Join(expectedTokens, ", "), t.tokenType, t.value)
}

func (p *parser) statementReady(statement Statement) parseFunc {
	p.statement = statement

	return nil
}

func parseStatement(p *parser) parseFunc {
	t := p.next(true)

	switch t.tokenType {
	case tokenError:
		return p.errorf(t.value)
	case tokenSelect:
		return parseSelect
	case tokenInsert:
		return parseInsert
	case tokenUpdate:
		return parseUpdate
	case tokenDelete:
		return parseDelete
	case tokenCreate:
		return parseCreateTable
	case tokenDrop:
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
	s := &Select{}
	for {
		t, err := p.scanFor(tokenIdentifier)
		if err != nil {
			return p.error(err)
		}

		s.Columns = append(s.Columns, t.value)

		t, err = p.scanFor(tokenFrom, tokenDelimeter)
		if err != nil {
			return p.error(err)
		}

		if t.tokenType == tokenFrom {
			break
		}
	}

	t, err := p.scanFor(tokenIdentifier)
	if err != nil {
		return p.error(err)
	}

	s.Table = t.value

	t, err = p.scanFor(tokenWhere, tokenLimit, tokenEnd)
	if err != nil {
		return p.error(err)
	}

	switch t.tokenType {
	case tokenWhere:
		expr, err := parseExpression(p)
		if err != nil {
			return p.error(err)
		}

		s.Where = &Where{expr}
	case tokenLimit:
		t, err = p.scanFor(tokenInteger)
		if err != nil {
			return p.error(err)
		}

		s.Limit = t.value
	}

	return p.statementReady(s)
}

// parseDelete parses INSERT statement.
func parseInsert(p *parser) parseFunc {
	t, err := p.scanFor(tokenIdentifier, tokenInto)
	if err != nil {
		return p.error(err)
	}

	if t.tokenType == tokenInto {
		// INTO is an optional keyword
		t, err = p.scanFor(tokenIdentifier)
		if err != nil {
			return p.error(err)
		}
	}

	insert := &Insert{t.value, []string{}, []string{}}

	_, err = p.scanFor(tokenLeftParenthesis)
	if err != nil {
		return p.error(err)
	}

	for {
		t, err := p.scanFor(tokenIdentifier)
		if err != nil {
			return p.error(err)
		}

		insert.Columns = append(insert.Columns, t.value)

		t, err = p.scanFor(tokenDelimeter, tokenRightParenthesis)
		if err != nil {
			return p.error(err)
		}

		if t.tokenType == tokenRightParenthesis {
			break
		}
	}

	t, err = p.scanFor(tokenValues)
	if err != nil {
		return p.error(err)
	}

	t, err = p.scanFor(tokenLeftParenthesis)
	if err != nil {
		return p.error(err)
	}

	for {
		t, err = p.scanFor(tokenInteger, tokenString)
		if err != nil {
			return p.error(err)
		}

		insert.Values = append(insert.Values, t.value)

		t, err = p.scanFor(tokenDelimeter, tokenRightParenthesis)
		if err != nil {
			return p.error(err)
		}

		if t.tokenType == tokenRightParenthesis {
			break
		}
	}

	_, err = p.scanFor(tokenEnd)
	if err != nil {
		return p.error(err)
	}

	return p.statementReady(insert)
}

// parseDelete parses UPDATE statement.
func parseUpdate(p *parser) parseFunc {
	t, err := p.scanFor(tokenIdentifier)
	if err != nil {
		return p.error(err)
	}

	update := &Update{t.value, []string{}, []string{}, nil}

	_, err = p.scanFor(tokenSet)
	if err != nil {
		return p.error(err)
	}

	for {
		t, err := p.scanFor(tokenIdentifier)
		if err != nil {
			return p.error(err)
		}

		update.Columns = append(update.Columns, t.value)

		t, err = p.scanFor(tokenAssign)
		if err != nil {
			return p.error(err)
		}

		t, err = p.scanFor(tokenString, tokenInteger)
		if err != nil {
			return p.error(err)
		}

		update.Values = append(update.Values, t.value)

		t, err = p.scanFor(tokenDelimeter, tokenEnd)
		if err != nil {
			return p.error(err)
		}

		if t.tokenType == tokenEnd {
			break
		}
	}

	return p.statementReady(update)
}

// parseDelete parses DELETE statement.
func parseDelete(p *parser) parseFunc {
	t, err := p.scanFor(tokenFrom)
	if err != nil {
		return p.error(err)
	}

	t, err = p.scanFor(tokenIdentifier)
	if err != nil {
		return p.error(err)
	}

	delete := &Delete{t.value, nil}

	_, err = p.scanFor(tokenEnd)
	if err != nil {
		return p.error(err)
	}

	return p.statementReady(delete)
}

// parseCreateTable parses CREATE TABLE statement.
func parseCreateTable(p *parser) parseFunc {
	t, err := p.scanFor(tokenTable)
	if err != nil {
		return p.error(err)
	}

	t, err = p.scanFor(tokenIdentifier)
	if err != nil {
		return p.error(err)
	}

	createTable := &CreateTable{t.value, []ColumnDefinition{}}

	_, err = p.scanFor(tokenLeftParenthesis)
	if err != nil {
		return p.error(err)
	}

	for {
		t, err := p.scanFor(tokenIdentifier)
		if err != nil {
			return p.error(err)
		}

		columnName := t.value

		t, err = p.scanFor(tokenTypeInteger, tokenTypeString)
		if err != nil {
			return p.error(err)
		}

		var columnType ColumnType
		switch t.tokenType {
		case tokenTypeInteger:
			columnType = TypeInteger
		case tokenTypeString:
			columnType = TypeString
		}

		createTable.Columns = append(createTable.Columns, ColumnDefinition{columnName, columnType})

		t, err = p.scanFor(tokenDelimeter, tokenRightParenthesis)
		if err != nil {
			return p.error(err)
		}

		if t.tokenType == tokenRightParenthesis {
			break
		}
	}

	_, err = p.scanFor(tokenEnd)
	if err != nil {
		return p.error(err)
	}

	return p.statementReady(createTable)
}

// parseDropTable parses DROP TABLE statement.
func parseDropTable(p *parser) parseFunc {
	t, err := p.scanFor(tokenTable)
	if err != nil {
		return p.error(err)
	}

	t, err = p.scanFor(tokenIdentifier)
	if err != nil {
		return p.error(err)
	}

	dropTable := &DropTable{t.value}

	_, err = p.scanFor(tokenEnd)
	if err != nil {
		return p.error(err)
	}

	return p.statementReady(dropTable)
}

func parseExpression(p *parser) (Expr, error) {
	t, err := p.scanFor(tokenLeftParenthesis, tokenIdentifier, tokenTypeInteger, tokenTypeString)
	if err != nil {
		return nil, err
	}

	var left Expr
	switch t.tokenType {
	case tokenIdentifier:
		left = ExprIdentifier{t.value}
	case tokenTypeInteger:
		left = ExprValueInteger{t.value}
	case tokenTypeString:
		left = ExprValueString{t.value}
	}

	t, err = p.scanFor(tokenEquals)
	if err != nil {
		return nil, err
	}

	var operator Operator
	switch t.tokenType {
	case tokenEquals:
		operator = OperatorEquals
	}

	t, err = p.scanFor(tokenIdentifier, tokenTypeInteger, tokenTypeString)
	if err != nil {
		return nil, err
	}

	var right Expr
	switch t.tokenType {
	case tokenIdentifier:
		right = ExprIdentifier{t.value}
	case tokenTypeInteger:
		right = ExprValueInteger{t.value}
	case tokenTypeString:
		right = ExprValueString{t.value}
	}

	var expr = ExprOperation{left, operator, right}

	return expr, nil
}
