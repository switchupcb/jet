package postgres

import (
	"github.com/switchupcb/jet/v2/notinternal/jet"
	"github.com/switchupcb/jet/v2/notinternal/utils/is"
)

type OnConflict interface {
	ON_CONSTRAINT(name string) ConflictTarget
	WHERE(indexPredicate BoolExpression) ConflictTarget
	ConflictTarget
}

type ConflictTarget interface {
	DO_NOTHING() InsertStatement
	DO_UPDATE(action ConflictAction) InsertStatement
}

type OnConflictClause struct {
	insertStatement  InsertStatement
	constraint       string
	indexExpressions []jet.ColumnExpression
	whereClause      jet.ClauseWhere
	do               jet.Serializer
}

func (o *OnConflictClause) ON_CONSTRAINT(name string) ConflictTarget {
	o.constraint = name
	return o
}

func (o *OnConflictClause) WHERE(indexPredicate BoolExpression) ConflictTarget {
	o.whereClause.Condition = indexPredicate
	return o
}

func (o *OnConflictClause) DO_NOTHING() InsertStatement {
	o.do = jet.Keyword("DO NOTHING")
	return o.insertStatement
}

func (o *OnConflictClause) DO_UPDATE(action ConflictAction) InsertStatement {
	o.do = action
	return o.insertStatement
}

func (o *OnConflictClause) Serialize(statementType jet.StatementType, out *jet.SQLBuilder, options ...jet.SerializeOption) {
	if is.Nil(o.do) {
		return
	}

	out.NewLine()
	out.WriteString("ON CONFLICT")
	if len(o.indexExpressions) > 0 {
		out.WriteString("(")
		jet.SerializeColumnExpressions(o.indexExpressions, statementType, out, jet.ShortName)
		out.WriteString(")")
	}

	if o.constraint != "" {
		out.WriteString("ON CONSTRAINT")
		out.WriteString(o.constraint)
	}

	o.whereClause.Serialize(statementType, out, jet.SkipNewLine, jet.ShortName)

	out.IncreaseIdent(7)
	jet.Serialize(o.do, statementType, out)
	out.DecreaseIdent(7)
}
