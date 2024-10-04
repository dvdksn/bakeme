package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"text/tabwriter"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users = map[int]*user{}
	seq   = 1
	lock  = sync.Mutex{}
	// the version of the program
	version string
)

//----------
// Handlers
//----------

func createUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := &user{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	users[id].Email = u.Email
	return c.JSON(http.StatusOK, users[id])
}

func deleteUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

func getAllUsers(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, users)
}

func newServer() *echo.Echo {
	e := echo.New()

	// Hide startup banner
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", getAllUsers)
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	return e
}

func main() {
	server := newServer()
	fmt.Print("Users API\n")

	if version != "" {
		fmt.Printf("version: %s\n", version)
	}

	configureLogging(server)

	fmt.Print("\nRoutes:\n\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 3, ' ', 0)
	fmt.Fprintln(w, "\tGET\t/users\tGet all users")
	fmt.Fprintln(w, "\tPOST\t/users\tCreate a user")
	fmt.Fprintln(w, "\tGET\t/users/{id}\tGet a user")
	fmt.Fprintln(w, "\tPUT\t/users/{id}\tEdit a user")
	fmt.Fprintln(w, "\tDELETE\t/users/{id}\tDelete a user")
	fmt.Fprintln(w)
	w.Flush()

	// Start the server
	server.Logger.Fatal(server.Start(":1323"))
}
