package insertadmin_handler

import (
	"net/http"

	"github.com/Robert076/doclane/backend/services"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func InsertAdminHandler(w http.ResponseWriter, r *http.Request) {
	adminParams := services.RegisterParams{
		Email:     "admin@admin.com",
		Password:  "adminadmin",
		FirstName: "Adminescu",
		LastName:  "Adminovici",
		Role:      "admin",
	}

	id, err := config.UserService.AddUser(r.Context(), adminParams)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Admin seeded successfully",
		Data:    id,
	})
}
