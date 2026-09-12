package sungrow

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/imroc/req/v3"
)

type Credentials struct {
	AppID     string
	AppSecret string

	UserAccount  string
	UserPassword string

	BaseURL string
}

func (c Credentials) Validate() error {
	if strings.TrimSpace(c.AppID) == "" {
		return errors.New("[Sungrow]: AppID 是空值")
	}

	if strings.TrimSpace(c.AppSecret) == "" {
		return errors.New("[Sungrow]: AppSecret 是空值")
	}

	if strings.TrimSpace(c.UserAccount) == "" {
		return errors.New("[Sungrow]: UserAccount 是空值")
	}

	if strings.TrimSpace(c.UserPassword) == "" {
		return errors.New("[Sungrow]: UserPassword 是空值")
	}

	return nil
}

var ErrNotLoggedIn = errors.New("[Sungrow]: 未获取 Token 令牌")

type SungrowSDK struct {
	creds  Credentials
	client *req.Client

	mu           sync.RWMutex
	token        string
	refreshToken string
	expiry       time.Time
	uid          int64

	loginMu sync.Mutex
}
