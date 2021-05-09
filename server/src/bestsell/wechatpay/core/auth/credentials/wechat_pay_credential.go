//Package credentials 微信支付api v3 go http-client authorization生成器
package credentials

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"bestsell/wechatpay/core/auth"
	"bestsell/wechatpay/core/consts"
)

// WechatPayCredentials authorization生成器
type WechatPayCredentials struct {
	Signer              auth.Signer // 签名器
	MchID               string      // 商户号
	CertificateSerialNo string      // 商户证书序列号
}

// GenerateAuthorizationHeader  生成http request header 中的authorization信息
func (c *WechatPayCredentials) GenerateAuthorizationHeader(ctx context.Context,
	method, canonicalURL, signBody string) (authorization string, err error) {
	if c.Signer == nil {
		return "", fmt.Errorf("you must init WechatPayCredentials with signer")
	}
	nonce, err := generateNonceStr()
	if err != nil {
		return "", err
	}
	timestamp := time.Now().Unix()
	message := fmt.Sprintf(consts.FormatMessage, method, canonicalURL, timestamp, nonce, signBody)
	signature, err := c.Signer.Sign(ctx, message)
	if err != nil {
		return "", err
	}
	authorization = fmt.Sprintf(consts.HeaderAuthorization, c.MchID, nonce, timestamp, c.CertificateSerialNo, signature)
	return authorization, nil
}

func (c *WechatPayCredentials) GenerateSignaturePayInfo(ctx context.Context, params *map[string]interface{}, reply *map[string]interface{}) error {
	if c.Signer == nil {
		return fmt.Errorf("you must init WechatPayCredentials with signer")
	}
	nonceStr, err := generateNonceStr()
	if err != nil {
		return err
	}
	timeStamp := time.Now().Unix()
	appId := (*params)["appId"].(string)
	prepay_id := (*params)["prepay_id"].(string)
	_package := "prepay_id="+prepay_id
	message := fmt.Sprintf("%s\n%d\n%s\n%s\n", appId, timeStamp, nonceStr, _package)
	fmt.Println("message:", message)
	signature, err := c.Signer.Sign(ctx, message)
	if err != nil {
		return err
	}
	(*reply)["timeStamp"] = strconv.FormatInt(timeStamp, 10)
	(*reply)["nonceStr"] = nonceStr
	(*reply)["package"] = _package
	(*reply)["signType"] = "RSA"
	(*reply)["paySign"] = signature
	return nil
}

func generateNonceStr() (string, error) {
	bytes := make([]byte, consts.NonceLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	symbolsByteLength := byte(len(consts.Symbols))
	for i, b := range bytes {
		bytes[i] = consts.Symbols[b%symbolsByteLength]
	}
	return string(bytes), nil
}
