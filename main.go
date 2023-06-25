package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Kucing struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World! pake echo")
}

func getKucingFunc(c echo.Context) error {
	// Get team and member from the query string
	kucing := c.QueryParam("kucing")
	status := c.QueryParam("status")

	dataType := c.Param("type")

	//return c.String(http.StatusOK, "Kucingnya jenis:"+kucing+", keadaan sekarang:"+status)
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("Ini sekarang jenis kucingnya adalah %s dan kucing ini keadaannya %s", kucing, status))
	}
	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"kucing": kucing,
			"status": status,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "tipe harus string atau json",
	})
}

func AddKucingFunc(c echo.Context) error {
	kucing := Kucing{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&kucing)
	if err != nil {
		log.Printf("Gagal melakukan decode %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "Gagal melakukan decode",
		})
	}

	// == save to database here ==
	log.Printf("Berhasil Menyimpan kucing dari request %v", kucing)
	return c.JSON(http.StatusOK, map[string]string{
		"status": "Success",
	})
}

func getDashboard(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "dashboard berhasil",
	})
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Server Pluto")
		c.Response().Header().Set("Developer", "t3d")
		return next(c)
	}
}

func main() {
	fmt.Println("server berjalan")

	e := echo.New()

	// Server header
	e.Use(ServerHeader)

	g := e.Group("/api/v1")

	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("andi")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("123456")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	g.GET("/dashboard", getDashboard)

	e.GET("/", home)
	e.GET("/getKucing/:type", getKucingFunc)
	e.POST("/addKucing", AddKucingFunc)

	e.Start(":8080")
}
