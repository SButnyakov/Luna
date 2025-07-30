package httphandlers

import (
	"net/http"

	"github.com/SButnyakov/luna/audio-processing/internal/dto"
	"github.com/gin-gonic/gin"
)

type processAudioRequest struct {
	AudioLink string `json:"audio_link"`
}

type processAudioResponse struct {
	Status  string         `json:"status"`
	Results map[int]string `json:"results"` // битрейт : ссылка на плейлист
}

func ProcessAudio(process func(dto.ProcessAudioDTO) (map[int]string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req processAudioRequest

		if err := c.BindJSON(&req); err != nil {
			c.IndentedJSON(http.StatusBadRequest, nil) // TODO: httputils.ErrorStatus
			return
		}

		processDTO := dto.ProcessAudioDTO{
			FileLink: req.AudioLink,
		}

		results, err := process(processDTO)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, nil) // TODO: errors generics
			return
		}

		c.IndentedJSON(http.StatusOK, processAudioResponse{
			Status:  "OK",
			Results: results,
		})
	}
}
