package postgres

import "testing"

func TestOnConflict(t *testing.T) {

	assertClauseSerialize(t, &OnConflictClause{}, "")

	OnConflict := &OnConflictClause{}
	OnConflict.DO_NOTHING()
	assertClauseSerialize(t, OnConflict, `
ON CONFLICT DO NOTHING`)

	OnConflict = &OnConflictClause{indexExpressions: ColumnList{table1ColBool}}
	OnConflict.DO_NOTHING()
	assertClauseSerialize(t, OnConflict, `
ON CONFLICT (col_bool) DO NOTHING`)

	OnConflict = &OnConflictClause{indexExpressions: ColumnList{table1ColBool}}
	OnConflict.ON_CONSTRAINT("table_pkey").DO_NOTHING()
	assertClauseSerialize(t, OnConflict, `
ON CONFLICT (col_bool) ON CONSTRAINT table_pkey DO NOTHING`)

	OnConflict = &OnConflictClause{indexExpressions: ColumnList{table1ColBool, table2ColFloat}}
	OnConflict.WHERE(table2ColFloat.ADD(table1ColInt).GT(table1ColFloat)).
		DO_UPDATE(
			SET(table1ColBool.SET(Bool(true)),
				table1ColInt.SET(Int(11))).
				WHERE(table2ColFloat.GT(Float(11.1))),
		)
	assertClauseSerialize(t, OnConflict, `
ON CONFLICT (col_bool, col_float) WHERE (col_float + col_int) > col_float DO UPDATE
       SET col_bool = $1::boolean,
           col_int = $2
       WHERE table2.col_float > $3`)
}
