package model

import (
	"errors"
)

type TokenChannelOverride struct {
	Id          int    `json:"id"`
	TokenId     int    `json:"token_id" gorm:"uniqueIndex:idx_token_channel;not null"`
	ChannelId   int    `json:"channel_id" gorm:"uniqueIndex:idx_token_channel;not null"`
	OverrideKey string `json:"override_key" gorm:"not null"`
	TokenName   string `json:"token_name" gorm:"-"`
	ChannelName string `json:"channel_name" gorm:"-"`
}

func (TokenChannelOverride) TableName() string {
	return "token_channel_overrides"
}

func GetTokenChannelOverridesByTokenId(tokenId int) ([]*TokenChannelOverride, error) {
	var overrides []*TokenChannelOverride
	err := DB.Where("token_id = ?", tokenId).Find(&overrides).Error
	if err != nil {
		return nil, err
	}
	for _, o := range overrides {
		token, err := GetTokenById(o.TokenId)
		if err == nil {
			o.TokenName = token.Name
		}
		channel, err := GetChannelById(o.ChannelId, false)
		if err == nil {
			o.ChannelName = channel.Name
		}
	}
	return overrides, nil
}

func GetTokenChannelOverride(tokenId int, channelId int) (*TokenChannelOverride, error) {
	var override TokenChannelOverride
	err := DB.Where("token_id = ? AND channel_id = ?", tokenId, channelId).First(&override).Error
	if err != nil {
		return nil, err
	}
	return &override, nil
}

func GetOverrideKey(tokenId int, channelId int) (string, error) {
	override, err := GetTokenChannelOverride(tokenId, channelId)
	if err != nil {
		return "", err
	}
	return override.OverrideKey, nil
}

func (o *TokenChannelOverride) Insert() error {
	if o.TokenId == 0 || o.ChannelId == 0 {
		return errors.New("token_id 和 channel_id 不能为空")
	}
	if o.OverrideKey == "" {
		return errors.New("override_key 不能为空")
	}
	token, err := GetTokenById(o.TokenId)
	if err != nil {
		return errors.New("令牌不存在")
	}
	_, err = GetChannelById(o.ChannelId, false)
	if err != nil {
		return errors.New("渠道不存在")
	}
	o.TokenName = token.Name
	return DB.Create(o).Error
}

func (o *TokenChannelOverride) Update() error {
	if o.Id == 0 {
		return errors.New("id 不能为空")
	}
	if o.OverrideKey == "" {
		return errors.New("override_key 不能为空")
	}
	return DB.Model(o).Update("override_key", o.OverrideKey).Error
}

func DeleteTokenChannelOverride(id int, tokenId int) error {
	if id == 0 || tokenId == 0 {
		return errors.New("id 或 token_id 不能为空")
	}
	return DB.Where("id = ? AND token_id = ?", id, tokenId).Delete(&TokenChannelOverride{}).Error
}

func DeleteTokenChannelOverridesByTokenId(tokenId int) error {
	if tokenId == 0 {
		return errors.New("token_id 不能为空")
	}
	return DB.Where("token_id = ?", tokenId).Delete(&TokenChannelOverride{}).Error
}

func DeleteTokenChannelOverridesByChannelId(channelId int) error {
	if channelId == 0 {
		return errors.New("channel_id 不能为空")
	}
	return DB.Where("channel_id = ?", channelId).Delete(&TokenChannelOverride{}).Error
}

func GetTokenChannelOverridesByUserId(userId int) ([]*TokenChannelOverride, error) {
	var overrides []*TokenChannelOverride
	tokenIds, err := getTokenIdsByUserId(userId)
	if err != nil {
		return nil, err
	}
	if len(tokenIds) == 0 {
		return []*TokenChannelOverride{}, nil
	}
	err = DB.Where("token_id IN ?", tokenIds).Find(&overrides).Error
	if err != nil {
		return nil, err
	}
	for _, o := range overrides {
		token, err := GetTokenById(o.TokenId)
		if err == nil {
			o.TokenName = token.Name
		}
		channel, err := GetChannelById(o.ChannelId, false)
		if err == nil {
			o.ChannelName = channel.Name
		}
	}
	return overrides, nil
}

func getTokenIdsByUserId(userId int) ([]int, error) {
	var ids []int
	err := DB.Model(&Token{}).Where("user_id = ?", userId).Pluck("id", &ids).Error
	return ids, err
}

func GetTokenChannelOverridesForUser(tokenId int, userId int) ([]*TokenChannelOverride, error) {
	token, err := GetTokenByIds(tokenId, userId)
	if err != nil {
		return nil, err
	}
	return GetTokenChannelOverridesByTokenId(token.Id)
}

func (o *TokenChannelOverride) Clean() {
	o.OverrideKey = MaskTokenKey(o.OverrideKey)
}
