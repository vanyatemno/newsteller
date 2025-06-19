package routes

import "github.com/gofiber/fiber/v2"

type Routable interface {
	SetRoutes(app *fiber.App)
}

type Routes struct {
}

func New() *Routes {
	return &Routes{}
}

func (r *Routes) InitializeRoutes(app *fiber.App, routes ...Routable) {
	for i := range routes {
		routes[i].SetRoutes(app)
	}
}
