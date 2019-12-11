package cas

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tmsong/hlog"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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

func PostByJson(realUrl string, reqBodyStr string, l *hlog.Logger) string {
	payload := strings.NewReader(reqBodyStr)
	req, _ := http.NewRequest(http.MethodPost, realUrl, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Golang CAS client")
	l.AddHttpTrace(req)
	res, _ := http.DefaultClient.Do(req)
	defer func() {
		if res.Body != nil {
			_ = res.Body.Close()
		}
	}()
	body, _ := ioutil.ReadAll(res.Body)
	resBodyStr := string(body)
	defer printHttpLog(l, req, res, reqBodyStr, resBodyStr)
	return resBodyStr
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

func StructToMap(data interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	j, err := json.Marshal(data)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(j, &m)
	if err != nil {
		return m, err
	}
	return m, err
}
func MapToStruct(data map[string]interface{}, v interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, &v)
	return err
}

func ListToStructList(data []interface{}, v interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, &v)
	return err
}

func InterfaceToStruct(data interface{}, v interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, &v)
}
