package dbengine

import (
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// SqlOperator represents comparison and logical operators in SQL queries.
type SqlOperator string

const (
	SqlOpEqual              SqlOperator = "="
	SqlOpNotEqual           SqlOperator = "!="
	SqlOpNotEqualAlt        SqlOperator = "<>"
	SqlOpLessThan           SqlOperator = "<"
	SqlOpLessThanOrEqual    SqlOperator = "<="
	SqlOpGreaterThan        SqlOperator = ">"
	SqlOpGreaterThanOrEqual SqlOperator = ">="
	SqlOpLike               SqlOperator = "LIKE"
	SqlOpNotLike            SqlOperator = "NOT LIKE"
	SqlOpIn                 SqlOperator = "IN"
	SqlOpNotIn              SqlOperator = "NOT IN"
	SqlOpIsNull             SqlOperator = "IS NULL"
	SqlOpIsNotNull          SqlOperator = "IS NOT NULL"
)

// Name returns the operator string name.
func (o SqlOperator) Name() string {
	return string(o)
}

// String returns the operator as a string.
func (o SqlOperator) String() string {
	return string(o)
}

// Value returns the raw operator value.
func (o SqlOperator) Value() string {
	return string(o)
}

// IsCompare compares the current operator to target.
func (o SqlOperator) IsCompare(target SqlOperator) bool {
	return o == target
}

// IsEnum checks if the operator is valid using an O(1) map.
func (o SqlOperator) IsEnum() bool {
	return sqlOperatorValidMap[o]
}

// IsEqual checks if the operator is an equality comparison.
func (o SqlOperator) IsEqual() bool {
	return o == SqlOpEqual
}

// IsNotEqual checks if the operator is an inequality comparison.
func (o SqlOperator) IsNotEqual() bool {
	if o == SqlOpNotEqual {
		return true
	}
	return o == SqlOpNotEqualAlt
}

// IsLessThan checks if the operator is less than.
func (o SqlOperator) IsLessThan() bool {
	return o == SqlOpLessThan
}

// IsLessThanOrEqual checks if the operator is less than or equal.
func (o SqlOperator) IsLessThanOrEqual() bool {
	return o == SqlOpLessThanOrEqual
}

// IsGreaterThan checks if the operator is greater than.
func (o SqlOperator) IsGreaterThan() bool {
	return o == SqlOpGreaterThan
}

// IsGreaterThanOrEqual checks if the operator is greater than or equal.
func (o SqlOperator) IsGreaterThanOrEqual() bool {
	return o == SqlOpGreaterThanOrEqual
}

// IsLike checks if the operator is a LIKE clause.
func (o SqlOperator) IsLike() bool {
	return o == SqlOpLike
}

// IsIn checks if the operator is an IN clause.
func (o SqlOperator) IsIn() bool {
	return o == SqlOpIn
}

// MarshalJSON marshals the operator to JSON.
func (o SqlOperator) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(o))
}

// UnmarshalJSON unmarshals the operator from JSON with validation.
func (o *SqlOperator) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	target := SqlOperator(s)
	if !sqlOperatorValidMap[target] {
		return fmt.Errorf("invalid sql operator: %s", s)
	}
	*o = target
	return nil
}

// ToJSON serializes the operator to a JSON string wrapped in an AppError.
func (o SqlOperator) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(string(o))
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize sql operator to json")
	}
	return string(b), nil
}

// FromJSON deserializes and validates the operator from a JSON string.
func (o *SqlOperator) FromJSON(s string) *apperror.AppError {
	var str string
	if err := json.Unmarshal([]byte(s), &str); err != nil {
		return apperror.WrapSimple(err, "deserialize sql operator from json")
	}
	target := SqlOperator(str)
	if !sqlOperatorValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid sql operator: %s", str), "validate sql operator from json")
	}
	*o = target
	return nil
}

type sqlOperatorRegistry struct {
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
func (r sqlOperatorRegistry) All() []SqlOperator {
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
func (r sqlOperatorRegistry) Names() []string {
	all := r.All()
	names := make([]string, len(all))
	for i, op := range all {
		names[i] = op.String()
	}
	return names
}

// IsEnum checks if an operator is valid.
func (r sqlOperatorRegistry) IsEnum(target SqlOperator) bool {
	return sqlOperatorValidMap[target]
}

// IsEqual checks if target matches Equal.
func (r sqlOperatorRegistry) IsEqual(target SqlOperator) bool {
	return target == r.Equal
}

// ToJSON serializes the registry to JSON.
func (r sqlOperatorRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r.Names())
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize sql operator registry to json")
	}
	return string(b), nil
}

// SqlOperators provides scoped access to SQL operators: SqlOperators.Equal, SqlOperators.Like.
var SqlOperators = sqlOperatorRegistry{
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
