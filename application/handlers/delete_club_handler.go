package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/delete_club"
	"sport_platform/application/models/shared"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

func DeleteClubHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	claimsRaw, exists := ctx.Get(middleware.ClaimsKey)
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Unauthorized",
			},
		)
		return
	}

	userClaims := claimsRaw.(claims.UserClaims)

	if userClaims.Role != shared.Teacher && userClaims.Role != shared.Admin {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{
				"message": "No permission",
			},
		)
		return
	}

	var request delete_club.DeleteClubRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid club ID format in URL",
		})
		return
	}

	if request.ID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid club ID",
		})
		return
	}

	if userClaims.Role == shared.Teacher {
		count, err := wrapper.Db.Queries.CheckClubOwnership(ctx, db_queries.CheckClubOwnershipParams{
			ID:        request.ID,
			TeacherID: userClaims.ID,
		})
		if err != nil || count == 0 {
			ctx.JSON(
				http.StatusForbidden,
				gin.H{
					"message": "You can only delete your own clubs",
				},
			)
			return
		}
	}

	err := wrapper.Db.Queries.SoftDeleteClub(ctx, request.ID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete club",
		})
		return
	}

	response := delete_club.DeleteClubResponse{
		ClubID: request.ID,
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
