package postgres

import "github.com/switchupcb/jet/v2/notinternal/jet"

// UNION effectively appends the result of sub-queries(select statements) into single query.
// It eliminates duplicate rows from its result.
func UNION(lhs, rhs jet.SerializerStatement, selects ...jet.SerializerStatement) SetStatement {
	return newSetStatementImpl(union, false, toSelectList(lhs, rhs, selects...))
}

// UNION_ALL effectively appends the result of sub-queries(select statements) into single query.
// It does not eliminates duplicate rows from its result.
func UNION_ALL(lhs, rhs jet.SerializerStatement, selects ...jet.SerializerStatement) SetStatement {
	return newSetStatementImpl(union, true, toSelectList(lhs, rhs, selects...))
}

// INTERSECT returns all rows that are in query results.
// It eliminates duplicate rows from its result.
func INTERSECT(lhs, rhs jet.SerializerStatement, selects ...jet.SerializerStatement) SetStatement {
	return newSetStatementImpl(intersect, false, toSelectList(lhs, rhs, selects...))
}

// INTERSECT_ALL returns all rows that are in query results.
// It does not eliminates duplicate rows from its result.
func INTERSECT_ALL(lhs, rhs jet.SerializerStatement, selects ...jet.SerializerStatement) SetStatement {
	return newSetStatementImpl(intersect, true, toSelectList(lhs, rhs, selects...))
}

// EXCEPT returns all rows that are in the result of query lhs but not in the result of query rhs.
// It eliminates duplicate rows from its result.
func EXCEPT(lhs, rhs jet.SerializerStatement) SetStatement {
	return newSetStatementImpl(except, false, toSelectList(lhs, rhs))
}

// EXCEPT_ALL returns all rows that are in the result of query lhs but not in the result of query rhs.
// It does not eliminates duplicate rows from its result.
func EXCEPT_ALL(lhs, rhs jet.SerializerStatement) SetStatement {
	return newSetStatementImpl(except, true, toSelectList(lhs, rhs))
}

type SetStatement interface {
	setOperators

	ORDER_BY(orderByClauses ...OrderByClause) SetStatement

	LIMIT(limit int64) SetStatement
	OFFSET(offset int64) SetStatement
	// OFFSET_e can be used when an integer expression is needed as offset, otherwise OFFSET can be used
	OFFSET_e(offset IntegerExpression) SetStatement

	AsTable(alias string) SelectTable
}

type setOperators interface {
	Statement
	jet.HasProjections
	Expression

	UNION(rhs SelectStatement) SetStatement
	UNION_ALL(rhs SelectStatement) SetStatement
	INTERSECT(rhs SelectStatement) SetStatement
	INTERSECT_ALL(rhs SelectStatement) SetStatement
	EXCEPT(rhs SelectStatement) SetStatement
	EXCEPT_ALL(rhs SelectStatement) SetStatement
}

type setOperatorsImpl struct {
	parent setOperators
}

func (s *setOperatorsImpl) UNION(rhs SelectStatement) SetStatement {
	return UNION(s.parent, rhs)
}

func (s *setOperatorsImpl) UNION_ALL(rhs SelectStatement) SetStatement {
	return UNION_ALL(s.parent, rhs)
}

func (s *setOperatorsImpl) INTERSECT(rhs SelectStatement) SetStatement {
	return INTERSECT(s.parent, rhs)
}

func (s *setOperatorsImpl) INTERSECT_ALL(rhs SelectStatement) SetStatement {
	return INTERSECT_ALL(s.parent, rhs)
}

func (s *setOperatorsImpl) EXCEPT(rhs SelectStatement) SetStatement {
	return EXCEPT(s.parent, rhs)
}

func (s *setOperatorsImpl) EXCEPT_ALL(rhs SelectStatement) SetStatement {
	return EXCEPT_ALL(s.parent, rhs)
}

type SetStatementImpl struct {
	jet.ExpressionStatement

	setOperatorsImpl

	setOperator jet.ClauseSetStmtOperator
}

func newSetStatementImpl(operator string, all bool, selects []jet.SerializerStatement) SetStatement {
	newSetStatement := &SetStatementImpl{}
	newSetStatement.ExpressionStatement = jet.NewExpressionStatementImpl(Dialect, jet.SetStatementType, newSetStatement,
		&newSetStatement.setOperator)

	newSetStatement.setOperator.Operator = operator
	newSetStatement.setOperator.All = all
	newSetStatement.setOperator.Selects = selects
	newSetStatement.setOperator.Limit.Count = -1

	newSetStatement.setOperatorsImpl.parent = newSetStatement

	return newSetStatement
}

func (s *SetStatementImpl) ORDER_BY(orderByClauses ...OrderByClause) SetStatement {
	s.setOperator.OrderBy.List = orderByClauses
	return s
}

func (s *SetStatementImpl) LIMIT(limit int64) SetStatement {
	s.setOperator.Limit.Count = limit
	return s
}

func (s *SetStatementImpl) OFFSET(offset int64) SetStatement {
	s.setOperator.Offset.Count = Int(offset)
	return s
}

func (s *SetStatementImpl) OFFSET_e(offset IntegerExpression) SetStatement {
	s.setOperator.Offset.Count = offset
	return s
}

func (s *SetStatementImpl) AsTable(alias string) SelectTable {
	return newSelectTable(s, alias, nil)
}

const (
	union     = "UNION"
	intersect = "INTERSECT"
	except    = "EXCEPT"
)

func toSelectList(lhs, rhs jet.SerializerStatement, selects ...jet.SerializerStatement) []jet.SerializerStatement {
	return append([]jet.SerializerStatement{lhs, rhs}, selects...)
}
