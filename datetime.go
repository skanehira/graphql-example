package main

import (
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func coerceDateTime(value interface{}) interface{} {
	switch v := value.(type) {
	case *string:
		if v == nil {
			return nil
		}
		coerceDateTime(*v)
	case string:
		if v == "" {
			return nil
		}
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil
		}
		return coerceDateTime(i)
	case int64:
		return time.Unix(0, v).In(time.UTC)
	case *int64:
		if v == nil {
			return nil
		}
		return coerceDateTime(*v)
	case time.Time:
		return v.UnixNano()
	case *time.Time:
		return coerceDateTime(*v)
	}

	return nil
}

var DateTime = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "DateTime",
	Description: "golang time.Time",
	Serialize:   coerceDateTime,
	ParseValue:  coerceDateTime,
	ParseLiteral: func(v ast.Value) interface{} {
		return nil
	},
})
