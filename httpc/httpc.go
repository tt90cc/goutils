package httpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/tt90cc/utils/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

func DeleteEmptyValue(src map[string]interface{}) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for key, value := range src {
		if key != "" && value != nil && value != "" {
			resultMap[key] = value
		}
	}
	return resultMap
}

func FormatSignSrcText(method string, paramMap map[string]interface{}) (string, error) {
	validParamMap := DeleteEmptyValue(paramMap)
	if strings.EqualFold(method, "GET") {
		keys := make([]string, 0)
		for k := range validParamMap {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		tmpList := make([]string, 0)
		for i := range keys {
			switch tmpValue := validParamMap[keys[i]]; tmpValue.(type) {
			case string:
				tmpList = append(tmpList, fmt.Sprintf("%s=%s", keys[i], tmpValue))
				continue
			case interface{}:
				rs, err := json.Marshal(tmpValue)
				if err != nil || rs == nil || string(rs) == "" {
					continue
				}
				tmpList = append(tmpList, fmt.Sprintf("%s=%s", keys[i], string(rs)))
				continue
			}
		}
		return strings.Join(tmpList, "&"), nil

	} else if strings.EqualFold(method, "POST") {
		postResult, postErr := json.Marshal(validParamMap)
		return string(postResult), postErr
	} else {
		return "", errors.New("Unknow Method: \"" + method + "\"")
	}
}

func BaseResponse(ctx context.Context, url string, data interface{}, header ...http.Header) (map[string]interface{}, error) {
	logger := logx.WithContext(ctx)

	var reader io.Reader
	if data != nil {
		b, _ := json.Marshal(data)
		reader = bytes.NewReader(b)
	}

	r, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	for _, h := range header {
		for k, listStr := range h {
			for _, s := range listStr {
				r.Header.Set(k, s)
			}
		}
	}

	resp, err := httpc.DoRequest(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var j map[string]interface{}
	err = json.Unmarshal(b, &j)

	logger.Infof("url:%s respData:%s", url, string(b))

	return j, err
}

func Post(ctx context.Context, url string, data interface{}, header ...http.Header) (interface{}, error) {
	logger := logx.WithContext(ctx)

	j, err := BaseResponse(ctx, url, data, header...)
	if err != nil {
		logger.Errorf("request failed. err:%v url:%s data:%+v j:%v", err, url, data, j)
		return nil, err
	}

	if _, ok := j["code"]; !ok {
		logger.Errorf("not found the key of code. url:%s data:%+v j:%v", url, data, j)
		return nil, errors.New("not found the key of code")
	}

	if j["code"].(float64) != errorx.OK {
		logger.Errorf("code validate failed. url:%s data:%+v j:%v", url, data, j)
		return nil, errorx.NewCodeError(cast.ToInt(j["code"]), cast.ToString(j["message"]))
	}

	return j["data"], nil
}

func Get(ctx context.Context, url string, header ...http.Header) (interface{}, error) {
	logger := logx.WithContext(ctx)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for _, h := range header {
		for k, listStr := range h {
			for _, s := range listStr {
				r.Header.Set(k, s)
			}
		}
	}

	resp, err := httpc.DoRequest(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var j map[string]interface{}
	err = json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	if _, ok := j["code"]; !ok {
		logger.Errorf("not found the key of code. url:%s j:%v", url, j)
		return nil, errors.New("not found the key of code")
	}

	if j["code"].(float64) != errorx.OK {
		logger.Errorf("code validate failed. url:%s j:%v", url, j)
		return nil, errorx.NewCodeError(cast.ToInt(j["code"]), cast.ToString(j["message"]))
	}

	return j["data"], nil
}

type CustomizeConfig struct {
	URL    string
	Data   interface{}
	Header http.Header
}

func CustomizePost(ctx context.Context, conf *CustomizeConfig, fn func(resp map[string]interface{}) interface{}) (interface{}, error) {
	logger := logx.WithContext(ctx)

	h := make([]http.Header, 0)
	if conf.Header != nil {
		h = append(h, conf.Header)
	}
	j, err := BaseResponse(ctx, conf.URL, conf.Data, h...)
	if err != nil {
		logger.Errorf("request failed. err:%v url:%s data:%+v j:%v", err, conf.URL, conf.Data, j)
		return nil, err
	}

	return fn(j), nil
}

func CustomizeGet(ctx context.Context, conf *CustomizeConfig, fn func(resp map[string]interface{}) interface{}) (interface{}, error) {
	logger := logx.WithContext(ctx)
	r, err := http.NewRequest(http.MethodGet, conf.URL, nil)
	if err != nil {
		return nil, err
	}

	for k, listStr := range conf.Header {
		for _, s := range listStr {
			r.Header.Set(k, s)
		}
	}

	resp, err := httpc.DoRequest(r)
	if err != nil {
		logger.Errorf("request failed. err:%v", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var j map[string]interface{}
	err = json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	return fn(j), nil
}
