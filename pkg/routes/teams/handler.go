package teams

import (
	"mike/pkg/application"
	"mike/pkg/routes/teams/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CheckTeamExists(c echo.Context) error {
	app := c.Get("app").(*application.App)
	teamId := c.Param("id")
	teamIdInt, err := strconv.Atoi(teamId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}
	exists, err := models.CheckTeamExists(app.DB, teamIdInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
}

func GetTeamDetails(c echo.Context) error {
	app := c.Get("app").(*application.App)
	teamId := c.Param("id")
	teamIdInt, err := strconv.Atoi(teamId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}
	team, err := models.GetTeamDetails(app.DB, teamIdInt)
	if err != nil {
		if err.Error() == "team not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Team not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get team details"})
	}
	return c.JSON(http.StatusOK, team)
}
