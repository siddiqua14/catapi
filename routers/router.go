package routers

import (
	"catapi/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Define routes for the application
	web.Router("/", &controllers.CatController{}, "get:ServeSingleCatImage")
	web.Router("/api/catimages", &controllers.CatController{}, "get:ServeMultipleCatImages")

	web.Router("/vote", &controllers.CatController{}, "post:CreateVote")
	web.Router("/votes", &controllers.CatController{}, "get:GetVotes")
	web.Router("/api/breeds", &controllers.CatController{}, "get:GetBreeds")
	web.Router("/api/breed-images", &controllers.CatController{}, "get:GetBreedImages")
	web.Router("/createFavorite", &controllers.CatController{}, "post:CreateFavorite")
	web.Router("/getFavorites", &controllers.CatController{}, "get:GetFavorites")
	web.Router("/deleteFavorite/:id", &controllers.CatController{}, "delete:DeleteFavorite")
}
