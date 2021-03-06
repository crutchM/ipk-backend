package handler

import (
	"github.com/gin-gonic/gin"
	"ipk/data/service"
)

//головной объект хендлера запросов, вся его суть просто навешивать методы на каждый путь и передавать им контекст gin
type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.OPTIONS("/", h.opt)
	//просто options запросы чтобы cors на фронте не ругался
	router.OPTIONS("/auth/sign-up", h.opt)
	router.OPTIONS("/api/user/teachers", h.opt)
	router.OPTIONS("/api/user/", h.opt)
	router.OPTIONS("/api/user/experts", h.opt)
	router.OPTIONS("/api/user/employments", h.opt)
	router.OPTIONS("/api/test/", h.opt)
	router.OPTIONS("/api/test/sendStat", h.opt)
	router.OPTIONS("/api/test/sendResults", h.opt)
	router.OPTIONS("/api/stat/getStat", h.opt)
	router.OPTIONS("/api/stat/getIndividual", h.opt)
	router.OPTIONS("/api/stat/remove", h.opt)

	//очевидно, группа запросов на авторизацию и регистрацию
	auth := router.Group("/auth")
	{
		auth.OPTIONS("/sign-in", h.opt)
		auth.GET("/", h.check)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.GET("/getall", h.getAll)

	}
	//основной набор методов апи userIdentity-просто метод, который проверяет валидность jwt токена, полученного после авторизации
	api := router.Group("/api", h.userIdentity)
	{
		user := api.Group("/user")
		{
			user.GET("/", h.getUser)
			user.GET("/teachers", h.getTeachers)
			user.GET("/experts", h.getExperts)
			user.GET("/employments", h.getEmployments)
		}
		chair := api.Group("/chair")
		{
			//вспомогательные методы для кафедр, возможно не будут использоваться
			chair.GET("/getall", h.getAllChairs)
			chair.POST("/create", h.createChair)
		}

		test := api.Group("/test")
		{
			//методв create-задел на будущее, пока необходим только один вариант опроса
			test.GET("/", h.GetTest)
			test.POST("/create", h.CreateTest)
			test.POST("/sendResults", h.SendResult)

			test.POST("/sendStat", h.SendStat)
		}
		stat := api.Group("/stat")
		{
			stat.POST("/getStat", h.getStat)
			stat.POST("/getIndividual", h.getStatByTeacher)
			stat.POST("/remove", h.removeUser)
		}
	}

	return router
}
