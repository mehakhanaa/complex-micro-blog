package serializers

import (
	"strings"

	"github.com/mehakhanaa/complex-micro-blog/models"
)

type UserProfileData struct {
	UID      uint64  `json:"uid"`
	Username string  `json:"username"`
	Nickname string  `json:"nickname"`
	Avatar   string  `json:"avatar_url"`
	Birth    *int64  `json:"birth"`
	Gender   *string `json:"gender"`
	Level    uint64  `json:"level"`
}

func NewUserProfileData(model *models.UserInfo) *UserProfileData {

	profile := new(UserProfileData)
	profile.UID = uint64(model.ID)
	profile.Username = model.UserName

	if model.NickName != nil {
		profile.Nickname = *model.NickName
	} else {

		profile.Nickname = model.UserName
	}

	var sb strings.Builder
	sb.WriteString("/resources/avatar/")
	sb.WriteString(model.Avatar)
	profile.Avatar = sb.String()

	if model.Birth != nil {
		birth := model.Birth.Unix()
		profile.Birth = &birth
	} else {
		profile.Birth = nil
	}
	if model.Gender != nil {
		profile.Gender = model.Gender
	} else {
		profile.Gender = nil
	}
	profile.Level = model.Level

	return profile
}

type UserToken struct {
	Token string `json:"token"`
}

func NewUserToken(token string) *UserToken {
	return &UserToken{Token: token}
}
