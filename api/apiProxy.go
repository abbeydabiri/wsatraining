package api

import (
	// "bytes"
	// "compress/gzip"
	// "crypto/tls"
	// "io"
	"encoding/json"
	"net/http"
	"strings"

	"wsatraining/config"
	"wsatraining/utils"
)

type apiResponse struct {
	Code    int
	Message string
	Body    interface{}
}

//apiProxy ...
func apiProxy(httpRes http.ResponseWriter, httpReq *http.Request) {
	apiResp := new(apiResponse)
	apiResp.Code = http.StatusInternalServerError
	httpRes.Header().Set("Content-Type", "application/json")

	if httpReq.URL.Path == "/wp/logout" {
		config.Get().Cookie = ""
		apiResp.Body = map[string]string{
			"Redirect": "/",
		}
		json.NewEncoder(httpRes).Encode(apiResp)
		return
	}

	if httpReq.URL.Path == "/wp/isloggedin" && config.Get().Cookie == "" {
		apiResp.Body = map[string]string{
			"Redirect": "/",
		}
		json.NewEncoder(httpRes).Encode(apiResp)
		return
	}

	if config.Get().Cookie != "" {
		httpReq.Header.Add("Cookie", config.Get().Cookie)
	}
	bodyResp, proxyResp := utils.CURL(config.Get().Proxy, httpReq)

	if proxyResp != nil {
		for _, sValueList := range proxyResp.Header {
			for _, sValue := range sValueList {
				if strings.HasPrefix(sValue, "wordpress_logged_in_") {
					config.Get().Cookie = sValue
				}
			}
			// httpRes.Header().Add(sKey, strings.Join(sValueList, " "))
		}
	}

	switch {
	case httpReq.URL.Path == "/wp/wp-login.php":
		if config.Get().Cookie == "" {
			apiResp.Message = "Invalid Login"
		} else {
			apiResp.Code = http.StatusOK
			apiResp.Message = "User Verified"
			apiResp.Body = map[string]string{
				"Redirect": "/#/dashboard",
			}
		}
	default:
		apiResp.Code = http.StatusOK
		apiResp.Body = map[string]string{
			"html": string(bodyResp),
		}
	}
	json.NewEncoder(httpRes).Encode(apiResp)
}
