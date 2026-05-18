package utils

import (
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

func GetCallerFromContext(ctx interface{ Value(any) any }) (types.CallerContext, error) {
	caller, ok := ctx.Value(types.CallerContextKey).(types.CallerContext)
	if !ok || caller.UserID == 0 {
		return types.CallerContext{}, errors.ErrUnauthorized{Msg: "Unauthorized."}
	}
	return caller, nil
}
