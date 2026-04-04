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

func GetUsersByDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> GetUsersByDepartmentHandler called")
	claims, err := utils.GetClaimsFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	deptIDStr := r.URL.Query().Get("department_id")
	if deptIDStr == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Missing department_id query parameter."})
		return
	}

	deptID, err := strconv.Atoi(deptIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid department_id value."})
		return
	}

	users, err := config.UserService.GetUsersByDepartment(r.Context(), *claims, deptID)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Could not fetch users by department."})
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Users retrieved successfully.",
		Data:    users,
	})
}
