package sports

import (
	"mike/pkg/application"
	"mike/pkg/routes/sports/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CheckSportExists(c echo.Context) error {
	app := c.Get("app").(*application.App)
	sportId := c.Param("id")
	sportIdInt, err := strconv.Atoi(sportId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sport ID"})
	}
	exists, err := models.CheckSportExists(app.DB, sportIdInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
}

func GetSportDetails(c echo.Context) error {
	app := c.Get("app").(*application.App)
	sportId := c.Param("id")
	sportIdInt, err := strconv.Atoi(sportId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sport ID"})
	}
	sport, err := models.GetSportDetails(app.DB, sportIdInt)
	if err != nil {
		if err.Error() == "sport not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Sport not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get sport details"})
	}
	return c.JSON(http.StatusOK, sport)
}

func GetAllSports(c echo.Context) error {
	app := c.Get("app").(*application.App)
	sports, err := models.GetAllSports(app.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get sports"})
	}
	return c.JSON(http.StatusOK, sports)
}
