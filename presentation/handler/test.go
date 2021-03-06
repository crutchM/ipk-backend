package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ipk/domain/model"
	"ipk/domain/model/stat"
	"net/http"
	"time"
)

type Input struct {
	Name   string        `json:"name"`
	Blocks []model.Block `json:"blocks"`
}

func (h *Handler) CreateTest(c *gin.Context) {
	var input Input
	if err := c.BindJSON(&input); err != nil { //получаем данные из тела
		logrus.Error(err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	test := model.Test{Name: input.Name, Blocks: input.Blocks}
	id, err := h.services.CreateTest(test) //записываем данные в базу
	if err != nil {
		logrus.Error(err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.addHeaders(c)
	c.Header("Access-Control-Allow-Origin", "http://192.168.11.40:3000")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	SendJSONResponse(c, "testId", id)
}

func (h *Handler) GetTest(c *gin.Context) {
	//id, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	logrus.Error(err.Error())
	//	newErrorResponse(c, http.StatusBadRequest, err.Error())
	//	return
	//}
	test, err := h.services.GetTest(1)
	if err != nil {
		logrus.Error(err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Access-Control-Allow-Origin", "http://192.168.11.40:3000")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	SendJSONResponse(c, "test", test)
}

type ResInput struct {
	Blocks []model.Block `json:"blocks"`
	Test   int           `json:"test"`
}

func (h *Handler) SendResult(c *gin.Context) {
	var input ResInput
	if err := c.BindJSON(&input); err != nil {
		logrus.Error(err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	h.addHeaders(c)
	err := h.services.AddResult(input.Blocks, input.Test)
	if err != nil {
		logrus.Error(err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Access-Control-Allow-Origin", "http://192.168.11.40:3000")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	SendJSONResponse(c, "status", "success")
}

type myInput struct {
	stat.Stat `json:"stat"`
	AnketDate int64 `json:"anketDate"`
	Date      int64 `json:"date"`
}

func (h *Handler) SendStat(c *gin.Context) {
	var input stat.Stat
	var inp myInput
	if err := c.BindJSON(&inp); err != nil {
		logrus.Error(err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input = inp.Stat
	input.LessonDate = time.Unix(inp.Date, 0).Format("2006-01-02")
	input.AnketDate = time.Unix(inp.AnketDate, 0).Format("2006-01-02")
	h.addHeaders(c)
	id, err := h.services.AddRow(input)
	if err != nil {
		logrus.Error(err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Access-Control-Allow-Origin", "http://192.168.11.40:3000")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	SendJSONResponse(c, "rowId", id)
}

func (h *Handler) removeUser(c *gin.Context) {
	var input inpt
	err := c.BindJSON(&input)
	if err != nil {
		logrus.Error(err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.RemoveUser(input.Id)
	if err != nil {
		logrus.Error(err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Access-Control-Allow-Origin", "http://192.168.11.40:3000")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	SendJSONResponse(c, "status", "successful")
}
