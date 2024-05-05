package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/indikay/notification-service/internal/biz"
)

const (
	validateKey = "code:telegram:%s"
)

type validateCode struct {
	data *Data
	log  *log.Helper
}

func NewValidateCodeRepo(data *Data) biz.ValidateCodeRepo {
	return &validateCode{
		data: data,
		log:  log.NewHelper(log.DefaultLogger),
	}
}

func (v *validateCode) Save(ctx context.Context, code, data string, timeoutInSecs int32) error {
	rs, err := v.data.redisCli.Set(ctx, fmt.Sprintf(validateKey, code), data, timeoutInSecs)
	v.log.Infof("Save code status %s", rs)
	return err
}

func (v *validateCode) Get(ctx context.Context, code string) (string, error) {
	rs, err := v.data.redisCli.Get(ctx, fmt.Sprintf(validateKey, code))
	if err != nil {
		return "", err
	}

	return rs, nil
}

func (v *validateCode) SaveToken(ctx context.Context, token, data string, timeoutInSecs int32) error {
	rs, err := v.data.redisCli.Set(ctx, fmt.Sprintf(validateKey, token), data, timeoutInSecs)
	v.log.Infof("Save code status %s", rs)
	return err
}

func (v *validateCode) GetToken(ctx context.Context, token string) (string, error) {
	rs, err := v.data.redisCli.Get(ctx, fmt.Sprintf(validateKey, token))
	if err != nil {
		return "", err
	}

	return rs, nil
}
