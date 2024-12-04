package services

import (
	"errors"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/models"
	"github.com/mehakhanaa/complex-micro-blog/stores"
	"github.com/mehakhanaa/complex-micro-blog/types"
	"github.com/mehakhanaa/complex-micro-blog/utils/converters"
	"github.com/mehakhanaa/complex-micro-blog/utils/encryptors"
	"github.com/mehakhanaa/complex-micro-blog/utils/generators"
	"github.com/mehakhanaa/complex-micro-blog/utils/validers"
)

type UserService struct {
	userStore *stores.UserStore
}

func (factory *Factory) NewUserService() *UserService {
	return &UserService{
		userStore: factory.storeFactory.NewUserStore(),
	}
}

func (service *UserService) GetUserInfoByUID(uid uint64) (*models.UserInfo, error) {
	user, err := service.userStore.GetUserByUID(uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) GetUserInfoByUsername(username string) (*models.UserInfo, error) {
	user, err := service.userStore.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) RegisterUser(username string, password string) error {

	if !validers.IsValidUsername(username) {
		return errors.New("invalid username")
	}
	if !validers.IsValidPassword(password) {
		return errors.New("invalid password")
	}

	_, err := service.userStore.GetUserByUsername(username)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("username already exists")
	}

	salt, err := generators.GenerateSalt(consts.SALT_LENGTH)
	if err != nil {
		return err
	}
	hashedPassword, err := encryptors.HashPassword(password, salt)
	if err != nil {
		return err
	}

	err = service.userStore.RegisterUserByUsername(username, salt, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) LoginUser(username string, password string, ip string, app string, device string) (string, error) {

	userAuthInfo, err := service.userStore.GetUserAuthInfoByUsername(username)
	if err != nil {
		return "", err
	}

	userLoginLog := &models.UserLoginLog{
		UID:         userAuthInfo.UID,
		LoginTime:   time.Now(),
		LoginIP:     ip,
		Application: app,
		Device:      device,
		IsSucceed:   false,
		IfChecked:   false,
	}

	err = encryptors.CompareHashPassword(userAuthInfo.PasswordHash, password, userAuthInfo.Salt)
	if err != nil {
		userLoginLog.Reason = "password error"
		inner_err := service.userStore.CreateUserLoginLog(userLoginLog)
		if inner_err != nil {
			return "", errors.Join(err, inner_err)
		}
		return "", errors.New("password error")
	}

	token, claims, err := generators.GenerateToken(userAuthInfo.UID, username)
	if err != nil {
		userLoginLog.Reason = "token generation error"
		inner_err := service.userStore.CreateUserLoginLog(userLoginLog)
		if inner_err != nil {
			return "", errors.Join(err, inner_err)
		}
		return "", err
	}

	err = service.userStore.CreateUserAvaliableToken(token, claims)
	if err != nil {
		userLoginLog.Reason = "token creation error"
		inner_err := service.userStore.CreateUserLoginLog(userLoginLog)
		if inner_err != nil {
			return "", errors.Join(err, inner_err)
		}
		return "", err
	}

	userLoginLog.IsSucceed = true
	userLoginLog.BearerToken = token
	err = service.userStore.CreateUserLoginLog(userLoginLog)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *UserService) UserUploadAvatar(uid uint64, fileHeader *multipart.FileHeader) error {

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	fileType, err := validers.ValidImageFile(
		fileHeader,
		&file,
		consts.MIN_AVATAR_SIZE,
		consts.MIN_AVATAR_SIZE,
		consts.MAX_AVATAR_FILE_SIZE,
	)
	if err != nil {
		return err
	}

	resizedAvatar, err := converters.ResizeAvatar(fileType, &file)
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString(strconv.FormatUint(uid, 10))
	sb.WriteRune('_')
	sb.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	sb.WriteString(".webp")

	return service.userStore.SaveUserAvatarByUID(uid, sb.String(), resizedAvatar)
}

func (service *UserService) UserUpdatePassword(username string, password string, newPassword string) error {

	userAuthInfo, err := service.userStore.GetUserAuthInfoByUsername(username)
	if err != nil {
		return err
	}

	err = encryptors.CompareHashPassword(userAuthInfo.PasswordHash, password, userAuthInfo.Salt)
	if err != nil {

		return errors.New("incorrect password")
	}

	hashedNewPassword, err := encryptors.HashPassword(newPassword, userAuthInfo.Salt)
	if err != nil {
		return err
	}

	err = service.userStore.UpdateUserPasswordByUsername(userAuthInfo.UserName, hashedNewPassword)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) UpdateUserInfo(uid uint64, reqBody *types.UserUpdateProfileBody) error {

	updatedProfile := &models.UserInfo{
		NickName: reqBody.NickName,
	}

	if reqBody.Birth != nil {
		birth := time.Unix(int64(*reqBody.Birth), 0)
		updatedProfile.Birth = &birth
	} else {
		updatedProfile.Birth = nil
	}
	if reqBody.Gender != nil {
		if *reqBody.Gender != "male" && *reqBody.Gender != "female" {
			updatedProfile.Gender = nil
		} else {
			updatedProfile.Gender = reqBody.Gender
		}
	}

	err := service.userStore.UpdateUserInfoByUID(uid, updatedProfile)
	if err != nil {
		return err
	}

	return nil
}
