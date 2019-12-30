package remote

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
	"github.com/pefish/go-json"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-reflect"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Remote struct {
	BaseUrl   string
	ApiKey    string
	ApiSecret string
	PubKey    string
}

var RemoteInstance *Remote

func (this *Remote) sign(method string, apiPath string, params map[string]interface{}) (sig string, apiNonce string) {
	sortedStr := ``
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sortedStr += k + `=` + go_reflect.Reflect.MustToString(params[k]) + `&`
	}
	sortedStr = strings.TrimSuffix(sortedStr, `&`)
	nonce := go_reflect.Reflect.MustToString(time.Now().UnixNano() / 1e6)
	toSignStr := method + `|` + apiPath + `|` + nonce + `|` + sortedStr
	go_logger.Logger.Debug(`to sign str is: `, toSignStr)
	return this.signEcc(toSignStr, this.ApiSecret), nonce
}

func (this *Remote) hash256(s string) string {
	hash_result := sha256.Sum256([]byte(s))
	hash_string := string(hash_result[:])
	return hash_string
}

func (this *Remote) signEcc(message, api_secret string) string {
	secret, _ := hex.DecodeString(api_secret)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), secret)
	sig, _ := privKey.Sign([]byte(this.hash256(this.hash256(message))))
	return fmt.Sprintf("%x", sig.Serialize())
}

func (this *Remote) verifyEcc(message, signature, apiKey string) bool {
	api_key, _ := hex.DecodeString(apiKey)

	pubKey, _ := btcec.ParsePubKey(api_key, btcec.S256())

	sigBytes, _ := hex.DecodeString(signature)
	sigObj, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		fmt.Println(`parse signature error`, err)
		return false
	}
	verified := sigObj.Verify([]byte(this.hash256(this.hash256(message))), pubKey)
	return verified
}

func (this *Remote) post(apiPath string, params map[string]interface{}) interface{} {
	sig, apiNonce := this.sign(`POST`, apiPath, params)
	resp, body := go_http.Http.MustPost(go_http.RequestParam{
		Url: this.BaseUrl + apiPath,
		Headers: map[string]interface{}{
			`BIZ-API-KEY`:       this.ApiKey,
			`BIZ-API-SIGNATURE`: sig,
			`BIZ-API-NONCE`:     apiNonce,
			`Content-Type`:      "application/x-www-form-urlencoded",
		},
		Params: params,
	})
	isValidRequest := this.verifyResponse(resp, body)
	if !isValidRequest {
		go_error.ThrowInternalError(`cobo response signature verify error`, nil)
	}
	go_logger.Logger.Debug(`verify post request sucess, param: `, body)
	result := go_json.Json.Parse(body).(map[string]interface{})
	if result[`success`].(bool) != true {
		go_error.Throw(result[`error_message`].(string)+result[`error_description`].(string), go_reflect.Reflect.MustToUint64(result[`error_code`].(float64)))
	}
	return result[`result`]
}

func (this *Remote) get(apiPath string, params map[string]interface{}) interface{} {
	sig, apiNonce := this.sign(`GET`, apiPath, params)
	go_logger.Logger.Debug(apiPath + ` sig is: ` + sig)
	resp, body := go_http.Http.MustGet(go_http.RequestParam{
		Url: this.BaseUrl + apiPath,
		Headers: map[string]interface{}{
			`BIZ-API-KEY`:       this.ApiKey,
			`BIZ-API-SIGNATURE`: sig,
			`BIZ-API-NONCE`:     apiNonce,
			`Content-Type`:      "application/x-www-form-urlencoded",
		},
		Params: params,
	})
	isValidRequest := this.verifyResponse(resp, body)
	if !isValidRequest {
		go_error.ThrowInternalError(`cobo response signature verify error`, nil)
	}
	go_logger.Logger.Debug(`verify get request sucess, param: `, body)
	result := go_json.Json.Parse(body).(map[string]interface{})
	if result[`success`].(bool) != true {
		go_error.Throw(`cobo response error: `+result[`error_message`].(string)+result[`error_description`].(string), go_reflect.Reflect.MustToUint64(result[`error_code`].(float64)))
	}
	return result[`result`]
}

func (this *Remote) verifyResponse(resp *http.Response, body string) bool {
	var timeStamp, signature, content string
	timeStamp = resp.Header.Get(`BIZ_TIMESTAMP`)
	signature = resp.Header.Get(`BIZ_RESP_SIGNATURE`)
	content = body + "|" + timeStamp
	nowTimestamp := time.Now().UnixNano() / 1e6
	if nowTimestamp-go_reflect.Reflect.MustToInt64(timeStamp) > 30*1000 {
		go_error.ThrowInternalWithInternalMsg(`request expired`, `timestamp is: `+timeStamp+`now time is: `+go_reflect.Reflect.MustToString(nowTimestamp))
	}
	go_logger.Logger.Info(`verify content is: `, content)
	return this.verifyEcc(content, signature, this.PubKey)
}
