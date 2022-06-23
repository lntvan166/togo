package utils

import (
	"fmt"
	"lntvan166/togo/model"
	"net/http"
	"time"
)

const MySQLTimeFormat = "2006-01-02 15:04:05"

func GetCurrentTime() string {
	return time.Now().Format(MySQLTimeFormat)
}

func CheckAccessPermission(w http.ResponseWriter, username string, taskUserID int) error {
	userID, err := model.GetUserIDByUsername(username)
	if err != nil {
		ERROR(w, http.StatusInternalServerError, fmt.Errorf(err.Error()))
		return err
	}

	if userID != taskUserID {
		ERROR(w, http.StatusBadRequest, fmt.Errorf("you are not allowed to access this task"))
		return err
	}

	return nil
}