package enums

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SqlOperatorType represents comparison and logical operators in SQL queries.
type SqlOperatorType string

// SqlOperator is a backward-compatible alias for SqlOperatorType.
type SqlOperator = SqlOperatorType

const (
	SqlOpEqual              SqlOperatorType = "="
	SqlOpNotEqual           SqlOperatorType = "!="
	SqlOpNotEqualAlt        SqlOperatorType = "<>"
	SqlOpLessThan           SqlOperatorType = "<"
	SqlOpLessThanOrEqual    SqlOperatorType = "<="
	SqlOpGreaterThan        SqlOperatorType = ">"
	SqlOpGreaterThanOrEqual SqlOperatorType = ">="
	SqlOpLike               SqlOperatorType = "LIKE"
	SqlOpNotLike            SqlOperatorType = "NOT LIKE"
	SqlOpIn                 SqlOperatorType = "IN"
	SqlOpNotIn              SqlOperatorType = "NOT IN"
	SqlOpIsNull             SqlOperatorType = "IS NULL"
	SqlOpIsNotNull          SqlOperatorType = "IS NOT NULL"
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

// IsValid reports whether the operator is a recognized enumeration value.
func (o SqlOperator) IsValid() bool {
	return o.IsEnum()
}

// IsEqual checks if the operator is an equality comparison.
func (o SqlOperator) IsEqual() bool {
	return o == SqlOpEqual
}

// IsNotEqual checks if the operator is an inequality comparison.
func (o SqlOperator) IsNotEqual() bool {
	return slices.Contains([]SqlOperator{SqlOpNotEqual, SqlOpNotEqualAlt}, o)
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

// ParseSqlOperatorType parses a raw string into a validated SqlOperatorType.
func ParseSqlOperatorType(s string) (SqlOperatorType, *apperror.AppError) {
	val := SqlOperatorType(strings.ToUpper(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSqlOperatorType", "E1000", fmt.Sprintf("invalid sql operator: %q", s)).WithContext("value", s)
}

// AllSqlOperatorTypes returns all valid SQL operators.
func AllSqlOperatorTypes() []SqlOperatorType {
	return []SqlOperatorType{
		SqlOpEqual,
		SqlOpNotEqual,
		SqlOpNotEqualAlt,
		SqlOpLessThan,
		SqlOpLessThanOrEqual,
		SqlOpGreaterThan,
		SqlOpGreaterThanOrEqual,
		SqlOpLike,
		SqlOpNotLike,
		SqlOpIn,
		SqlOpNotIn,
		SqlOpIsNull,
		SqlOpIsNotNull,
	}
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
		return apperror.NewValidationError(fmt.Sprintf("invalid sql operator: %s", s))
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
		return apperror.NewValidationError(fmt.Sprintf("invalid sql operator: %s", str))
	}

	*o = target

	return nil
}
