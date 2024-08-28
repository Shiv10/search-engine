package routes

import (
	"fmt"
	"os"
	"search-engine/db"
	"search-engine/utils"
	"search-engine/views"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type loginform struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type settingsform struct {
	Amount   int  `form:"amount"`
	SearchOn string `form:"searchOn"`
	AddNew   string `form:"addNew"`
}

type AdminClaims struct {
	User                 string `json:"user"`
	ID                   string `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

func LoginPostHandler(c *fiber.Ctx) error {
	input := loginform{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}

	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	if err != nil {
		c.Status(401)
		return c.SendString("<h2>Error: Authentication error</h2>")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)

	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong logging in</h2>")
	}

	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(200)
}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(200)
}

func DashboardHandler(c *fiber.Ctx) error {
	settings := db.SearchSetting{}
	err := settings.Get()
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Cannot get Settings</h2>")
	}

	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNew))

}

func DashboardPostHandler(c *fiber.Ctx) error {
	input := settingsform{}

	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Cannot parse Settings</h2>")
	}

	addNew := false;
	if input.AddNew == "on" {
		addNew = true
	}

	searchOn := false
	if input.SearchOn == "on" {
		searchOn = true
	}

	settings := &db.SearchSetting{}

	settings.Amount = uint(input.Amount)
	settings.AddNew = addNew
	settings.SearchOn = searchOn


	err := settings.Update()

	if err != nil {
		fmt.Println(err)
		c.Status(500)
		return c.SendString("<h2>Error: Cannot update Settings</h2>")
	}

	c.Append("HX-Refresh", "true")
	return c.SendStatus(200)
}

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")
	if cookie == "" {
		return c.Redirect("/login", 302)
	}

	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return c.Redirect("/login", 302)
	}

	_, ok := token.Claims.(*AdminClaims)

	if ok && token.Valid {
		return c.Next()
	}

	return c.Redirect("/login", 302)
}
