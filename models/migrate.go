package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	var err error

	if err = db.AutoMigrate(&UserInfo{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserAuthInfo{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserLoginLog{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserPostStatus{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserCommentStatus{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(&PostInfo{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(&CommentInfo{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(&ReplyInfo{}); err != nil {
		return err
	}

	return nil
}
