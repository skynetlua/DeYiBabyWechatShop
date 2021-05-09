//Package setting 微信支付api v3 go http-client 配置
package setting

import (
	"errors"
	"net/http"
	"time"

	"bestsell/wechatpay/core/auth"
)

// DialSettings 微信支付apiv3 go http-client需要的配置信息
type DialSettings struct {
	HTTPClient *http.Client
	Request    *http.Request
	UserAgent  string
	Credential auth.Credential // authorization生成器
	Validator  auth.Validator  // 校验器
	Timeout    time.Duration   // 超时时间
}

// Validate 校验请求配置是否有效
func (ds *DialSettings) Validate() error {
	if ds.Credential == nil {
		return errors.New("you must set credential with option.WithCredential or option.WithMerchant")
	}
	if ds.Validator == nil {
		return errors.New("you must set validator with option.WithValidator or option.WithWechatPay")
	}
	return nil
}
