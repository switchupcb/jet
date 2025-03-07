package sqlite

import (
	"github.com/switchupcb/jet/v2/notinternal/jet"
	"github.com/switchupcb/jet/v2/notinternal/utils/is"
)

type OnConflict interface {
	WHERE(indexPredicate BoolExpression) ConflictTarget
	ConflictTarget
}

type ConflictTarget interface {
	DO_NOTHING() InsertStatement
	DO_UPDATE(action ConflictAction) InsertStatement
}

type OnConflictClause struct {
	insertStatement  InsertStatement
	indexExpressions []jet.ColumnExpression
	whereClause      jet.ClauseWhere
	do               jet.Serializer
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

	o.whereClause.Serialize(statementType, out, jet.SkipNewLine, jet.ShortName)

	out.IncreaseIdent(7)
	jet.Serialize(o.do, statementType, out)
	out.DecreaseIdent(7)
}

type ConflictAction interface {
	jet.Serializer
	WHERE(condition BoolExpression) ConflictAction
}

// SET creates conflict action for ON_CONFLICT clause
func SET(assigments ...ColumnAssigment) ConflictAction {
	ConflictAction := updateConflictActionImpl{}
	ConflictAction.doUpdate = jet.KeywordClause{Keyword: "DO UPDATE"}
	ConflictAction.Serializer = jet.NewSerializerClauseImpl(&ConflictAction.doUpdate, &ConflictAction.set, &ConflictAction.where)
	ConflictAction.set = assigments
	return &ConflictAction
}

type updateConflictActionImpl struct {
	jet.Serializer

	doUpdate jet.KeywordClause
	set      jet.SetClauseNew
	where    jet.ClauseWhere
}

func (u *updateConflictActionImpl) WHERE(condition BoolExpression) ConflictAction {
	u.where.Condition = condition
	return u
}
