package user_handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var limitPtr, offsetPtr *int
	var orderByPtr, orderPtr *string

	// limit
	if l := q.Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limitPtr = &val
		} else {
			utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid limit value."})
			return
		}
	}

	// offset
	if o := q.Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offsetPtr = &val
		} else {
			utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid offset value."})
			return
		}
	}

	if ob := q.Get("orderBy"); ob != "" {
		orderByPtr = &ob
	}

	if o := q.Get("order"); o != "" {
		orderPtr = &o
	}

	users, err := config.UserService.GetUsers(limitPtr, offsetPtr, orderByPtr, orderPtr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: fmt.Sprintf("Could not fetch users. %v", err)})
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Users retrieved successfully.",
		Data:    users,
	})
}
