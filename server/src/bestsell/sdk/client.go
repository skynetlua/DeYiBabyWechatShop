package sdk

import (
	"bestsell/wechatpay/core"
	"bestsell/wechatpay/core/auth"
	"bestsell/wechatpay/core/auth/credentials"
	"bestsell/wechatpay/core/auth/signers"
	"bestsell/wechatpay/core/auth/validators"
	"bestsell/wechatpay/core/auth/verifiers"
	"bestsell/wechatpay/core/option"
	"bestsell/wechatpay/utils"
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
)

const (
	testMchID                   = "XXXXXXXXXXX"
	testCertificateSerialNumber = "XXXXXXXXXXXXXXXXXXXXXX"
	testPrivateKey              = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDwPOGPGKFghP0X
XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
TlQP6JYjV1fY76krynW9Ye94v0Xg0fyiT5VDuI9JlRopTbSKbbGMot5shOEo/00y
tba/8vbiflUIrW+cS7CheGs=
-----END PRIVATE KEY-----`

	testWechatCertSerialNumber = "XXXXXXXXXXXXXXXXXXXXXX"
	testWechatCertificateStr   = `-----BEGIN CERTIFICATE-----
MIID/DCCAuSgAwIBAgIUVZMY8oJoYTnHd7HAVRCgkEpvDAAwDQYJKoZIhvcNAQEL
BQAwXjELMAkGA1UEBhMCQ04xEzARBgNVBAoTClRlbnBheS5jb20xHTAbBgNVBAsT
XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
d2cLCoyUcpdM5zGarNeYe41K3ZrtMkOV2P26gxb7lL9FtDkYuptH9iGe/ekYxrws
E6A2joGmoriATnmnbdNbsg==
-----END CERTIFICATE-----`

	apiv3Key = "900"

	filePath  = ""
	fileName  = "picture.jpeg"
	postURL   = "https://api.mch.weixin.qq.com/v3/marketing/favor/users/oHkLxt_htg84TUEbzvlMwQzVDBqo/coupons"
	GetURL    = "https://api.mch.weixin.qq.com/v3/certificates"
	uploadURL = "https://api.mch.weixin.qq.com/v3/merchant/media/upload"
)

var (
	privateKey           *rsa.PrivateKey
	wechatPayCertificate *x509.Certificate
	credential           auth.Credential
	validator            auth.Validator
	ctx                  context.Context
	err                  error
)

func init() {
	privateKey, err = utils.LoadPrivateKey(testPrivateKey)
	if err != nil {
		panic(fmt.Errorf("load private err:%s", err.Error()))
	}
	wechatPayCertificate, err = utils.LoadCertificate(testWechatCertificateStr)
	if err != nil {
		panic(fmt.Errorf("load certificate err:%s", err.Error()))
	}
	ctx = context.Background()
}

func DecryptWeChatMsg(associatedData, nonce, ciphertext string) (string, error) {
	return utils.DecryptToString(apiv3Key, associatedData, nonce, ciphertext)
}

func PostWeChatPay(requestURL string, requestBody interface{}) ([]byte, error) {
	fmt.Println("PostWeChatPay post requestURL:", requestURL, "requestBody:", requestBody)

	opts := []option.ClientOption{
		option.WithMerchant(testMchID, testCertificateSerialNumber, privateKey),
		option.WithWechatPay([]*x509.Certificate{wechatPayCertificate}),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		fmt.Println("PostWeChatPay post init error:", err)
		return nil, err
	}
	response, err := client.Post(ctx, requestURL, requestBody)
	fmt.Println("PostWeChatPay post error:", err)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("PostWeChatPay read error:", err)
		return nil, err
	}
	fmt.Println("PostWeChatPay body:", string(body))
	return body, err
}

func GetWeChatPay(requestURL string) ([]byte, error) {
	fmt.Println("GetWeChatPay get requestURL:", requestURL)
	opts := []option.ClientOption{
		option.WithMerchant(testMchID, testCertificateSerialNumber, privateKey),
		option.WithWechatPay([]*x509.Certificate{wechatPayCertificate}),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		fmt.Println("GetWeChatPay get init error:", err)
		return nil, err
	}
	response, err := client.Get(ctx, requestURL)
	fmt.Println("GetWeChatPay get error:", err)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("GetWeChatPay read error:", err)
		return nil, err
	}
	fmt.Println("GetWeChatPay get:", string(body))
	return body, err
}

func SignatureWeChatPay(params *map[string]interface{}, reply *map[string]interface{}) error {
	opts := []option.ClientOption {
		option.WithMerchant(testMchID, testCertificateSerialNumber, privateKey),
		option.WithWechatPay([]*x509.Certificate{wechatPayCertificate}),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		fmt.Println("SignatureWeChatPay init error:", err)
		return err
	}
	err = client.SignaturePayInfo(ctx, params, reply)
	if err != nil {
		fmt.Println("SignatureWeChatPay signature error:", err)
		return err
	}
	return nil
}

func TestGet() {
	opts := []option.ClientOption{
		option.WithMerchant(testMchID, testCertificateSerialNumber, privateKey),
		option.WithWechatPay([]*x509.Certificate{wechatPayCertificate}),
	}
	client, err := core.NewClient(ctx, opts...)
	fmt.Println(err)
	response, err := client.Get(ctx, GetURL)
	fmt.Println(err)
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(err)
	fmt.Println(string(body))
	fmt.Println(core.CheckResponse(response))
}

type testData struct {
	StockID           string `json:"stock_id"`
	StockCreatorMchID string `json:"stock_creator_mchid"`
	OutRequestNo      string `json:"out_request_no"`
	AppID             string `json:"appid"`
}

func TestPost() {
	opts := []option.ClientOption{
		option.WithMerchant(testMchID, testCertificateSerialNumber, privateKey),
		option.WithWechatPay([]*x509.Certificate{wechatPayCertificate}),
	}
	client, err := core.NewClient(ctx, opts...)
	fmt.Println(err)
	data := &testData{
		StockID:           "xxx",
		StockCreatorMchID: "xxx",
		OutRequestNo:      "xxx",
		AppID:             "xxx",
	}
	response, err := client.Post(ctx, postURL, data)
	fmt.Println(err)
	fmt.Println(core.CheckResponse(response))
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(err)
	fmt.Println(string(body))
}

type meta struct {
	FileName string `json:"filename" binding:"required"` // 商户上传的媒体图片的名称，商户自定义，必须以JPG、BMP、PNG为后缀。
	Sha256   string `json:"sha256" binding:"required"`   // 图片文件的文件摘要，即对图片文件的二进制内容进行sha256计算得到的值。
}

func TestClient_Upload() {
	// 如果你有自定义的Signer或者Verfifer
	credential = &credentials.WechatPayCredentials{
		Signer:              &signers.Sha256WithRSASigner{PrivateKey: privateKey},
		MchID:               testMchID,
		CertificateSerialNo: testCertificateSerialNumber,
	}
	validator = &validators.WechatPayValidator{
		Verifier: &verifiers.WechatPayVerifier{
			Certificates: map[string]*x509.Certificate{
				testWechatCertSerialNumber: wechatPayCertificate,
			},
		},
	}
	client, err := core.NewClient(ctx, option.WithCredential(credential), option.WithValidator(validator))
	fmt.Println(err)
	pictureByes, err := ioutil.ReadFile(filePath)
	fmt.Println(err)
	// 计算文件序列化后的sha256
	h := sha256.New()
	_, err = h.Write(pictureByes)
	fmt.Println(err)
	metaObject := &meta{}
	pictureSha256 := h.Sum(nil)
	metaObject.FileName = fileName
	metaObject.Sha256 = fmt.Sprintf("%x", string(pictureSha256))
	metaByte, _ := json.Marshal(metaObject)
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	err = core.CreateFormField(writer, "meta", "application/json", metaByte)
	fmt.Println(err)
	err = core.CreateFormFile(writer, fileName, "image/jpg", pictureByes)
	fmt.Println(err)
	err = writer.Close()
	fmt.Println(err)
	response, err := client.Upload(ctx, uploadURL, string(metaByte), reqBody.String(), writer.FormDataContentType())
	fmt.Println(err)
	if response.Body != nil {
		defer response.Body.Close()
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(err)
	fmt.Println(string(body))
}
