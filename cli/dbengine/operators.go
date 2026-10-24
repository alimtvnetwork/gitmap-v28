package dbengine

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine/enums"
)

// SqlOperatorType is the canonical SQL operator enum, defined in the dbengine enums package.
type SqlOperatorType = enums.SqlOperatorType

// SqlOperator is a backward-compatible alias for SqlOperatorType.
type SqlOperator = enums.SqlOperatorType

const (
	SqlOpEqual              = enums.SqlOpEqual
	SqlOpNotEqual           = enums.SqlOpNotEqual
	SqlOpNotEqualAlt        = enums.SqlOpNotEqualAlt
	SqlOpLessThan           = enums.SqlOpLessThan
	SqlOpLessThanOrEqual    = enums.SqlOpLessThanOrEqual
	SqlOpGreaterThan        = enums.SqlOpGreaterThan
	SqlOpGreaterThanOrEqual = enums.SqlOpGreaterThanOrEqual
	SqlOpLike               = enums.SqlOpLike
	SqlOpNotLike            = enums.SqlOpNotLike
	SqlOpIn                 = enums.SqlOpIn
	SqlOpNotIn              = enums.SqlOpNotIn
	SqlOpIsNull             = enums.SqlOpIsNull
	SqlOpIsNotNull          = enums.SqlOpIsNotNull
)

// SqlOperators provides scoped access to SQL operators: SqlOperators.Equal, SqlOperators.Like.
var SqlOperators = enums.SqlOperators
