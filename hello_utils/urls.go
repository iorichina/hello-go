package hello_utils

import (
	"bytes"
	"net/url"
)

func HttpBuildQuery(params *map[string]string) (reqParam string) {
	if len(*params) <= 0 {
		return
	}

	var buf bytes.Buffer
	for k, v := range *params {
		rowK := url.QueryEscape(k)
		buf.WriteString(rowK)
		buf.WriteByte('=')
		rowV := url.QueryEscape(v)
		buf.WriteString(rowV)
		buf.WriteByte('&')
	}
	reqParam = buf.String()
	reqParam = reqParam[0 : len(reqParam)-1]
	return
}

func HttpBuildQueryAny(params *map[string]interface{}) (reqParam string) {
	if len(*params) <= 0 {
		return
	}

	var buf bytes.Buffer
	for k, v := range *params {
		rowK := url.QueryEscape(k)
		buf.WriteString(rowK)
		buf.WriteByte('=')
		rowV := url.QueryEscape(ToString(v))
		buf.WriteString(rowV)
		buf.WriteByte('&')
	}
	reqParam = buf.String()
	reqParam = reqParam[0 : len(reqParam)-1]
	return
}

func HttpBuildQueryArr(params *map[string][]string) (reqParam string, originParams *map[string]string) {
	if len(*params) <= 0 {
		return
	}

	var buf bytes.Buffer
	for k, v := range *params {
		rowK := url.QueryEscape(k)
		for _, vv := range v {
			buf.WriteString(rowK)
			buf.WriteByte('=')
			rowV := url.QueryEscape(vv)
			buf.WriteString(rowV)
			buf.WriteByte('&')
			(*originParams)[k] = vv
		}
	}
	reqParam = buf.String()
	reqParam = reqParam[0 : len(reqParam)-1]
	return
}
