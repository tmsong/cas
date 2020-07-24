/**
 * @note
 * log
 *
 * @author	songtianming
 * @date 	2019-12-03
 */
package cas

import (
	"github.com/sirupsen/logrus"
	"github.com/tmsong/hlog"
	"net/http"
	"net/url"
)

const (
	LOG_OK   string = "ok"
	LOG_TAG  string = "tag"
	LOG_FAIL string = "fail"
)

const (
	LOGCODE_NAME string = "code"
)

const (
	LOGTAG_REQUEST_OK string = "_com_http_success"
)

const (
	LOGTAG_REQUEST_ERR string = "_com_http_failure"
)

func printLogWithHttpCode(code int, l *hlog.Logger, fields ...logrus.Fields) {
	if l == nil {
		return
	}
	var log interface{}
	var tag interface{}
	f := logrus.Fields{}
	if len(fields) >= 1 {
		if item, ok := fields[0][LOG_TAG]; ok {
			tag = item
		}
		e := l.WithFields(fields[0])
		for i := 1; i < len(fields); i++ {
			if item, ok := fields[i][LOG_TAG]; ok {
				tag = item
			}
			e = e.WithFields(fields[i])
		}
		log = e
	} else {
		log = l
	}
	f = logrus.Fields{
		LOGCODE_NAME: code,
	}
	if tag != nil {
		f[LOG_TAG] = tag
	}
	if code != http.StatusOK {
		if _, ok := f[LOG_TAG]; ok {
			f[LOG_TAG] = LOGTAG_REQUEST_ERR
		}
		if _, ok := log.(*logrus.Entry); ok {
			log.(*logrus.Entry).WithFields(f).Errorln(LOG_FAIL)
		} else {
			log.(*hlog.Logger).WithFields(f).Errorln(LOG_FAIL)
		}
	} else {
		if _, ok := log.(*logrus.Entry); ok {
			log.(*logrus.Entry).WithFields(f).Infoln(LOG_OK)
		} else {
			log.(*hlog.Logger).WithFields(f).Infoln(LOG_OK)
		}
	}
}

func printHttpLog(l *hlog.Logger, req *http.Request, res *http.Response, reqBody, resBody string, fields logrus.Fields) {
	u := req.URL
	get, _ := url.ParseQuery(u.RawQuery)
	post, uErr := url.QueryUnescape(reqBody)
	if uErr != nil {
		post = reqBody
	}
	f := logrus.Fields{
		"api":    u.Path,
		"url":    u.String(),
		"out":    resBody,
		"get":    get,
		"post":   post,
		"header": req.Header,
		"method": req.Method,
		"code":   res.StatusCode,
	}
	printLogWithHttpCode(res.StatusCode, l, fields, f)
}
