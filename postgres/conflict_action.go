package postgres

import "github.com/switchupcb/jet/v2/notinternal/jet"

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
