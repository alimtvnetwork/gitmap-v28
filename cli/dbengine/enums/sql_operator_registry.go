package enums

import (
	"encoding/json"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SqlOperatorRegistry provides scoped access to SQL operators.
type SqlOperatorRegistry struct {
	Equal              SqlOperator
	NotEqual           SqlOperator
	NotEqualAlt        SqlOperator
	LessThan           SqlOperator
	LessThanOrEqual    SqlOperator
	GreaterThan        SqlOperator
	GreaterThanOrEqual SqlOperator
	Like               SqlOperator
	NotLike            SqlOperator
	In                 SqlOperator
	NotIn              SqlOperator
	IsNull             SqlOperator
	IsNotNull          SqlOperator
}

// All returns a slice of all valid SQL operators.
func (r SqlOperatorRegistry) All() []SqlOperator {
	return []SqlOperator{
		r.Equal,
		r.NotEqual,
		r.NotEqualAlt,
		r.LessThan,
		r.LessThanOrEqual,
		r.GreaterThan,
		r.GreaterThanOrEqual,
		r.Like,
		r.NotLike,
		r.In,
		r.NotIn,
		r.IsNull,
		r.IsNotNull,
	}
}

// Names returns the string representations of all operators.
func (r SqlOperatorRegistry) Names() []string {
	all := r.All()
	names := make([]string, len(all))
	for i, op := range all {
		names[i] = op.String()
	}

	return names
}

// IsEnum checks if an operator is valid.
func (r SqlOperatorRegistry) IsEnum(target SqlOperator) bool {
	return sqlOperatorValidMap[target]
}

// IsEqual checks if target matches Equal.
func (r SqlOperatorRegistry) IsEqual(target SqlOperator) bool {
	return target == r.Equal
}

// ToJSON serializes the registry to JSON.
func (r SqlOperatorRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r.Names())
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize sql operator registry to json")
	}

	return string(b), nil
}

// SqlOperators provides scoped access to SQL operators: SqlOperators.Equal, SqlOperators.Like.
var SqlOperators = SqlOperatorRegistry{
	Equal:              SqlOpEqual,
	NotEqual:           SqlOpNotEqual,
	NotEqualAlt:        SqlOpNotEqualAlt,
	LessThan:           SqlOpLessThan,
	LessThanOrEqual:    SqlOpLessThanOrEqual,
	GreaterThan:        SqlOpGreaterThan,
	GreaterThanOrEqual: SqlOpGreaterThanOrEqual,
	Like:               SqlOpLike,
	NotLike:            SqlOpNotLike,
	In:                 SqlOpIn,
	NotIn:              SqlOpNotIn,
	IsNull:             SqlOpIsNull,
	IsNotNull:          SqlOpIsNotNull,
}

var sqlOperatorValidMap = map[SqlOperator]bool{
	SqlOperators.Equal:              true,
	SqlOperators.NotEqual:           true,
	SqlOperators.NotEqualAlt:        true,
	SqlOperators.LessThan:           true,
	SqlOperators.LessThanOrEqual:    true,
	SqlOperators.GreaterThan:        true,
	SqlOperators.GreaterThanOrEqual: true,
	SqlOperators.Like:               true,
	SqlOperators.NotLike:            true,
	SqlOperators.In:                 true,
	SqlOperators.NotIn:              true,
	SqlOperators.IsNull:             true,
	SqlOperators.IsNotNull:          true,
}
