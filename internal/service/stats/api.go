package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ssstatistics/internal/bots/telegram"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
)

func ApiMmrGET(c *gin.Context) {
	var players []*domain.Player
	r.Database.
		Select("Ckey", "MMR").
		Find(&players)

	type mmr struct {
		Ckey string
		MMR  int32
	}

	var mmrs []*mmr
	for _, player := range players {
		mmrs = append(mmrs, &mmr{
			Ckey: player.Ckey,
			MMR:  player.MMR,
		})
	}

	c.JSON(200, mmrs)
}

func ApiMapsGET(c *gin.Context) {
	var MapStatistics []*domain.MapStats
	r.Database.
		Preload("MapAttributes").
		Find(&MapStatistics)

	type simpleMapStats struct {
		MapName    string
		ServerID   string
		Attributes map[string]string
	}

	var maps []*simpleMapStats
	for _, stats := range MapStatistics {
		simpleMapStat := &simpleMapStats{
			MapName:    stats.MapName,
			ServerID:   stats.ServerID,
			Attributes: make(map[string]string),
		}
		for _, attribute := range stats.MapAttributes {
			simpleMapStat.Attributes[attribute.Name] = attribute.Value
		}
		maps = append(maps, simpleMapStat)
	}

	c.JSON(200, maps)
}

func ApiSendFeedback(c *gin.Context) {
	type FeedbackForm struct {
		Username string `json:"username"`
		Text     string `json:"text"`
	}

	var form FeedbackForm
	if err := c.BindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, "Некорректный запрос.")
		return
	}

	if len(form.Username) < 3 || len(form.Username) > 50 {
		c.JSON(http.StatusBadRequest, "В имени должно быть от 3х до 50 символов.")
		return
	}

	if len(form.Text) < 10 || len(form.Text) > 255 {
		c.JSON(http.StatusBadRequest, "В сообщение должно быть от 10 до 255 символов..")
		return
	}

	msg := fmt.Sprintf(`
	Name: *%s*
	
Text: %s`, form.Username, form.Text)

	err := telegram.Bot.Send(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Сообщение отправлено.")
}

func ApiHeatmapsGET(c *gin.Context) {
	type Request struct {
		// explosions, deaths
		Type          string `form:"type"`
		MapResolution string `form:"mapresolution"`
		MapName       string `form:"mapname"`
	}

	type Response struct {
		MapBase64 string
		Error     string
	}

	var query Request
	err := c.BindQuery(&query)
	if err != nil {
		c.JSON(200, Response{Error: "The request is incorrect"})
		return
	}

	c.JSON(200, Response{
		MapBase64: "SoSi He-x bibu",
		Error:     "228 Server shlet tebya нахуй."})
}
