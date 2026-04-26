package controller

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
)

func GetTokenChannelOverrides(c *gin.Context) {
	tokenIdStr := c.Param("id")
	tokenId, err := strconv.Atoi(tokenIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的令牌ID",
		})
		return
	}
	userId := c.GetInt("id")
	overrides, err := model.GetTokenChannelOverridesForUser(tokenId, userId)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	for _, o := range overrides {
		o.Clean()
	}
	common.ApiSuccess(c, overrides)
}

func AddTokenChannelOverride(c *gin.Context) {
	override := model.TokenChannelOverride{}
	err := c.ShouldBindJSON(&override)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	userId := c.GetInt("id")
	token, err := model.GetTokenByIds(override.TokenId, userId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权操作此令牌",
		})
		return
	}
	override.TokenId = token.Id
	if err := override.Insert(); err != nil {
		common.ApiError(c, err)
		return
	}
	override.Clean()
	common.ApiSuccess(c, override)
}

func UpdateTokenChannelOverride(c *gin.Context) {
	override := model.TokenChannelOverride{}
	err := c.ShouldBindJSON(&override)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	userId := c.GetInt("id")
	existing, err := model.GetTokenChannelOverride(override.TokenId, override.ChannelId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "覆盖记录不存在",
		})
		return
	}
	_, err = model.GetTokenByIds(existing.TokenId, userId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权操作此令牌",
		})
		return
	}
	existing.OverrideKey = override.OverrideKey
	if err := existing.Update(); err != nil {
		common.ApiError(c, err)
		return
	}
	existing.Clean()
	common.ApiSuccess(c, existing)
}

func DeleteTokenChannelOverride(c *gin.Context) {
	tokenIdStr := c.Param("id")
	tokenId, err := strconv.Atoi(tokenIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的令牌ID",
		})
		return
	}
	overrideIdStr := c.Param("override_id")
	overrideId, err := strconv.Atoi(overrideIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的覆盖记录ID",
		})
		return
	}
	userId := c.GetInt("id")
	_, err = model.GetTokenByIds(tokenId, userId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权操作此令牌",
		})
		return
	}
	if err := model.DeleteTokenChannelOverride(overrideId, tokenId); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, nil)
}

func GetAllUserChannelOverrides(c *gin.Context) {
	userId := c.GetInt("id")
	overrides, err := model.GetTokenChannelOverridesByUserId(userId)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	for _, o := range overrides {
		o.Clean()
	}
	common.ApiSuccess(c, overrides)
}
