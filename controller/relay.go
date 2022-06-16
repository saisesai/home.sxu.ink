package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/saisesai/home.sxu.ink/model"
	"github.com/saisesai/home.sxu.ink/mqtt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

var ErrRelayNotExist = errors.New("relay not exist")

func HandleRelayAdd(L *logrus.Entry, ctx *gin.Context, pClientId string) {
	var exist = false
	_, err := model.FindRelayByClientId(pClientId)
	if err == nil {
		exist = true
	}
	relay, err := model.AddRelay(pClientId)
	if err != nil {
		L.WithError(err).Errorln("failed to add relay:", relay.ClientId)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// subscribe
	if !exist {
		token := mqtt.Client.Subscribe(pClientId+"/msg", 0, mqtt.RelayMsgHandler)
		token.Wait()
		if token.Error() != nil {
			L.WithError(token.Error()).Errorln("failed to subscribe:", relay.ClientId+"/msg")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": token.Error().Error(),
			})
			return
		}
		token = mqtt.Client.Subscribe(pClientId+"/oln", 0, mqtt.RelayOnlineHandler)
		token.Wait()
		if token.Error() != nil {
			L.WithError(token.Error()).Errorln("failed to subscribe:", relay.ClientId+"/oln")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": token.Error().Error(),
			})
			return
		}
	}
	L.Debugln(pClientId + " subscribed!")
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func HandleRelayDelete(L *logrus.Entry, ctx *gin.Context, pClientId string) {
	relay, err := model.FindRelayByClientId(pClientId)
	// error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // not found
			L.WithField("client_id", pClientId).Warningln("rely not exist!")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": ErrRelayNotExist.Error(),
			})
			return
		}
		// internal error
		L.WithError(err).Errorln("failed to find relay:", relay.ClientId)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = relay.Delete()
	if err != nil {
		// internal error
		L.WithError(err).Errorln("failed to delete relay:", relay.ClientId)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// unsubscribe
	_, err = model.FindRelayByClientId(pClientId)
	if err != nil { // not exist
		token := mqtt.Client.Unsubscribe(pClientId + "/oln")
		token.Wait()
		if token.Error() != nil {
			L.WithError(err).Errorln("failed to unsubscribe:" + pClientId + "/oln")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": token.Error().Error(),
			})
		}
		token = mqtt.Client.Unsubscribe(pClientId + "/msg")
		token.Wait()
		if token.Error() != nil {
			L.WithError(err).Errorln("failed to unsubscribe:" + pClientId + "/msg")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": token.Error().Error(),
			})
		}
	}
	L.Debugln(pClientId + " unsubscribed!")
	// ok
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
