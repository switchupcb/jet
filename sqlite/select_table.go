package sqlite

import "github.com/switchupcb/jet/v2/notinternal/jet"

// SelectTable is interface for MySQL sub-queries
type SelectTable interface {
	readableTable
	jet.SelectTable
}

type selectTableImpl struct {
	jet.SelectTable
	readableTableInterfaceImpl
}

func newSelectTable(selectStmt jet.SerializerHasProjections, alias string, columnAliases []jet.ColumnExpression) SelectTable {
	subQuery := &selectTableImpl{
		SelectTable: jet.NewSelectTable(selectStmt, alias, columnAliases),
	}

	subQuery.readableTableInterfaceImpl.parent = subQuery

	return subQuery
}
