package controllers

import (
	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/forms"
	"github.com/dbsSensei/filesystem-api/utils"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type ServerController struct {
	config *config.Config
	db     *gorm.DB
}

func NewServerController(config *config.Config, db *gorm.DB) *ServerController {
	return &ServerController{
		config: config,
		db:     db,
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Server
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response{data=forms.HealthCheckResponse}
// @Failure 500 {object} utils.Response{data=object}
// @Router /health [get]
func (s *ServerController) HealthCheck(c *gin.Context) {
	// Server Check
	dbConn, _ := s.db.DB()

	parsedDsn, _ := url.Parse(s.config.DBSource)
	dbHost := parsedDsn.Host
	dbName := parsedDsn.Path

	if dbHost == "" {
		// Parse DSN server format
		pairs := strings.Split(dbName, " ")
		serverData := make(map[string]string)
		for _, pair := range pairs {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				serverData[parts[0]] = parts[1]
			}
		}
		dbHost = serverData["host"] + ":" + serverData["port"]
		dbName = serverData["dbname"]
	}

	var databaseStatus string
	if err := dbConn.Ping(); err != nil {
		databaseStatus = "error"
	} else {
		databaseStatus = "ok"
	}

	c.JSON(http.StatusOK, utils.ResponseData("success", "Server running well", forms.HealthCheckResponse{
		ServerStatus:   "ok",
		DatabaseStatus: databaseStatus,
		DatabaseName:   dbName,
		DatabaseHost:   dbHost,
	},
	))
}
