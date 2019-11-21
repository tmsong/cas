package cas

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/**
 * @note
 * 获取md5 hash串
 * @param string text 源串
 *
 * @return string
 */
func MD5Hex(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func CreateBaseParams(appId int64, appKey string) map[string]interface{} {
	r := map[string]interface{}{}
	timestamp := time.Now().Unix()
	r["appId"] = appId
	r["sign"] = MD5Hex(fmt.Sprintf("%s%d", appKey, timestamp))
	r["timestamp"] = timestamp
	return r
}

func PostByJson(url string, bodyJson string) string {
	payload := strings.NewReader(bodyJson)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Golang CAS client gopkg.in/cas")
	res, _ := http.DefaultClient.Do(req)
	defer func() {
		if res.Body != nil {
			_ = res.Body.Close()
		}
	}()
	body, _ := ioutil.ReadAll(res.Body)
	if glog.V(2) {
		glog.Infof("post json url %v , body %v , response %v", url, bodyJson, string(body))
	}
	return string(body)
}

func JsonEncode(data interface{}) string {
	s, e := json.Marshal(data)
	if e != nil {
		return ""
	}
	return string(s)
}

func JsonDecode(data string, inter interface{}) error {
	return json.Unmarshal([]byte(data), inter)
}

func StructToMap(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	j, _ := json.Marshal(data)
	_ = json.Unmarshal(j, &m)
	return m
}
func MapToStruct(data map[string]interface{}, v interface{}) {
	j, _ := json.Marshal(data)
	_ = json.Unmarshal(j, &v)
}

func ListToStructList(data []interface{}, v interface{}) {
	j, _ := json.Marshal(data)
	_ = json.Unmarshal(j, &v)
}

func InterfaceToStruct(data interface{}, v interface{}) {
	j, _ := json.Marshal(data)
	_ = json.Unmarshal(j, &v)
}
