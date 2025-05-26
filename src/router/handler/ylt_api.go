package handler

import (
	"eshop_server/src/utils/log"
	"fmt"
	"io"
	"net/http"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"go.uber.org/zap"
)

// @Title	构建公共请求头
// @Return	map[string][]string
var YltRequestHeaders map[string][]string = map[string][]string{
	"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"},
	"content-type":       {"application/json;chartset=utf-8"},
	"accept":             {"application/json"},
	"Accept-Language":    {"zh-CN,zh;q=0.9"},
	"Cache-Control":      {"no-cache"},
	"Connection":         {"keep-alive"},
	"sec-ch-ua":          {`"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`},
	"sec-ch-ua-mobile":   {"?0"},
	"sec-ch-ua-platform": {`"Windows"`},
	"Sec-Fetch-Dest":     {"empty"},
	"Sec-Fetch-Mode":     {"cors"},
	"Sec-Fetch-Site":     {"same-origin"},
	"Origin":             {"https://yuanlitui.com"},
	"gt-token":           {""},
}

// @Title	用户登陆，获取token
// @Method	POST
// @Return	gt_token, cookie, error
func YltUserLogin(phone string, password string) (string, string, error) {
	client := &http.Client{}
	reqbody := fmt.Sprintf(`{"phoneNumber":"%s","password":"%s"}`, phone, password)
	var data = strings.NewReader(reqbody)
	req, err := http.NewRequest("POST", "https://yuanlitui.com/api/login/phoneNumberPassLogin", data)
	if err != nil {
		return "", "", fmt.Errorf("创建请求失败: %v", err)
	}
	YltRequestHeaders["Referer"] = []string{"https://yuanlitui.com/login?redirect=/"}
	req.Header = YltRequestHeaders
	// req.Header.Set("Cookie", "Hm_lvt_d489914e2e65c946cf3060f9534688f6=1745828422; HMACCOUNT=3165FE343E1CD2F4; Hm_lpvt_d489914e2e65c946cf3060f9534688f6=1745829833")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", "", fmt.Errorf("YltUserLogin请求接口失败, %v", err)
	}
	defer resp.Body.Close()

	// 获取Set-Cookie头信息
	cookies := resp.Header.Get("set-cookie")
	if cookies == "" {
		return "", "", fmt.Errorf("登陆响应请求头中未包含set-cookie信息")
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("读取响应体失败: %v", err)
	}
	log.Infof("YltUserLogin登陆成功, response body: %s\n", string(bodyText))
	/* Success Response Body:
	{
		"s": true,
		"c": 200,
		"m": null,
		"d": {
			"token": "user gt_token"
		}
	}
	*/
	loginResp, err := simplejson.NewJson(bodyText)
	if err != nil {
		fmt.Println(err)
		return "", "", fmt.Errorf("生成simplejson失败, %v", err)
	}
	gt_token, err := loginResp.Get("d").Get("token").String()
	if err != nil {
		fmt.Println(err)
		return "", "", fmt.Errorf("提取用户Token失败, %v", err)
	}

	return gt_token, cookies, nil
}

// @Title	获取登陆用户信息
// @Headers	gt-token + cookie
// @Method	GET
// @Return	nil
func YltGetUserInfo(gt_token string, cookie string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://yuanlitui.com/api/user/userInfo", nil)
	if err != nil {
		log.Error("YltGetUserInfo创建Request失败", zap.Error(err))
		return err
	}
	YltRequestHeaders["Referer"] = []string{"https://yuanlitui.com/"}
	YltRequestHeaders["Cookie"] = []string{cookie}
	YltRequestHeaders["gt-token"] = []string{gt_token}
	req.Header = YltRequestHeaders
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("YltGetUserInfo请求接口失败, %v", err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("", zap.Error(err))
		return err
	}
	fmt.Printf("YltGetUserInfo获取用户信息成功, response body: %s\n", string(bodyText))
	/* Success Response Body:
	{
		"s": true,
		"c": 200,
		"m": null,
		"d": {
			"userId": "1916441581532758017",
			"phoneNumber": "178****2032",
			"nickName": "原崽2032_250427",
			"avatarUrl": "https://puss.gt-it.cn/default-avatar",
			"brief": "原崽250427",
			"creatorStatus": "0",
			"customizeUrl": null,
			"userBgImgUrl": "https://puss.gt-it.cn/default-user-bg",
			"followingCount": 1,
			"fansCount": 0
		}
	}
	*/

	return nil
}

// @Title	创建订单
// @Method	POST
// @Return	orderId, qrcodeBase64, error
func YltCreateOrder(gt_token string, cookie string, productId string, customerPrice float64) (string, string, error) {
	client := &http.Client{}
	// var data = strings.NewReader(`{"productId":"5517","payType":"alipay","affiliate":"","sceneType":"pc","customerPrice":0.5}`)
	var data = strings.NewReader(fmt.Sprintf(`{"productId":"%s","payType":"alipay","affiliate":"","sceneType":"pc","customerPrice":%v}`, productId, customerPrice))
	req, err := http.NewRequest("POST", "https://yuanlitui.com/api/order/createOrder", data)
	if err != nil {
		log.Error("", zap.Error(err))
		return "", "", err
	}
	YltRequestHeaders["Cookie"] = []string{cookie}
	YltRequestHeaders["gt-token"] = []string{gt_token}
	req.Header = YltRequestHeaders
	resp, err := client.Do(req)
	if err != nil {
		log.Error("YltCreateOrder 请求接口失败", zap.Error(err))
		return "", "", fmt.Errorf("YltCreateOrder 请求接口失败, %v", err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("YltCreateOrder 读取响应Body失败", zap.Error(err))
		return "", "", fmt.Errorf("YltCreateOrder 读取响应Body失败, %v", err)
	}
	log.Infof("YltCreateOrder 创建订单成功, response body: %s\n", string(bodyText))
	/* Success Response Body:
	{
		"s": true,
		"c": 200,
		"m": null,
		"d": {
			"payObj": "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",
			"orderNo": "202504281757541916793993162407937",
			"price": 0.5
		}
	}
	*/

	// d.payObj 支付二维码: data:image/jpeg;base64,xxxx
	createOrderResp, err := simplejson.NewJson(bodyText)
	if err != nil {
		fmt.Println(err)
		return "", "", fmt.Errorf("生成simplejson失败, %v", err)
	}
	qrcodeBase64, err := createOrderResp.Get("d").Get("payObj").String()
	if err != nil {
		fmt.Println(err)
		return "", "", fmt.Errorf("提取qrcode信息失败, %v", err)
	}
	orderId, err := createOrderResp.Get("d").Get("orderNo").String()
	if err != nil {
		fmt.Println(err)
		return "", "", fmt.Errorf("提取orderId信息失败, %v", err)
	}
	return orderId, qrcodeBase64, nil
}

// @Title	轮询订单购买状态
// @Method	POST
func YltCheckOrder(gt_token string, cookie string, yltOrderId string) (payOk bool, err error) {
	if yltOrderId == "" {
		log.Error("YltCheckOrder订单ID为空")
		return false, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://yuanlitui.com/api/order/checkProductOrder?orderNo=%v", yltOrderId), nil)
	if err != nil {
		log.Error("YltCheckOrder构建Request失败", zap.Error(err))
		return false, err
	}
	YltRequestHeaders["Cookie"] = []string{cookie}
	YltRequestHeaders["gt-token"] = []string{gt_token}
	req.Header = YltRequestHeaders
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, fmt.Errorf("请求接口失败, %v", err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("", zap.Error(err))
		return false, err
	}
	fmt.Printf("YltCheckOrder获取订单状态成功， response body: %s\n", string(bodyText))
	/*
		Topay Response Body:
			{"s":true,"c":200,"m":null,"d":false}
		Success Response Body:
			{"s":true,"c":200,"m":null,"d":true}
	*/

	// d bool 支付是否成功
	checkOrderResp, err := simplejson.NewJson(bodyText)
	if err != nil {
		fmt.Println(err)
		return false, fmt.Errorf("生成simplejson失败, %v", err)
	}
	payOk, err = checkOrderResp.Get("d").Bool()
	if err != nil {
		fmt.Println(err)
		return false, fmt.Errorf("提取qrcode base64失败, %v", err)
	}
	return payOk, nil
}

// TODO
// @Title	获取购买的商品信息
// @Method	POST
func _YltGetPurchasedProductInfo(gt_token string, cookie string, productId string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://yuanlitui.com/api/product/getProductById?id=%s", productId), nil)
	// req, err := http.NewRequest("GET", "https://yuanlitui.com/api/product/getProductById?id=ar55", nil)
	if err != nil {
		log.Error("", zap.Error(err))
		return err
	}
	YltRequestHeaders["Referer"] = []string{fmt.Sprintf("https://yuanlitui.com/a/%s", productId)}
	YltRequestHeaders["Cookie"] = []string{cookie}
	YltRequestHeaders["gt-token"] = []string{gt_token}
	req.Header = YltRequestHeaders
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("请求接口失败, %v", err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("", zap.Error(err))
		return err
	}
	fmt.Printf("YltGetPurchasedProductInfo获取购买商品详情信息成功, response body: %s\n", string(bodyText))
	/*
		Success Response Body:
			{
				"s": true,
				"c": 200,
				"m": null,
				"d": {
					"id": "5517",
					"title": "测试",
					"status": "0",
					"paidFlag": true,
					"ownerFlag": false,
					"content": "<p>测试</p>",
					"paidContent": "<p>测试</p>",
					"price": 0.50,
					"customerPriceFlag": false,
					"coverUrl": "https://puss.gt-it.cn/FmgfP2ZGXlfVEcDfmoSj6YW2bWT2",
					"protectFlag": false,
					"imageCount": null,
					"wordCount": null,
					"fileList": [],
					"collectionList": [],
					"creatorInfo": {
						"id": "1915408799566147585",
						"uid": "jiangkwkmia",
						"nickName": "江夫人_KWKmia",
						"avatarUrl": "https://puss.gt-it.cn/Fql7ME3_uNKSKs6HAsdqM2d39KIe",
						"brief": "音频作品｜剧情｜熟女｜虚拟UP主\n谢谢你喜欢妾身的作品～\n版权所有禁止二传、盗卖。",
						"customizeUrl": "jiangkwkmia"
					},
					"affiliateFlag": false,
					"affiliateRatio": null,
					"affiliateUrl": "https://yuanlitui.com/a/ar55",
					"affiliateImg": "iVBORw0KGgoAAAANSUhEUgAAASwAAAEsCAYAAAB5fY51AAAICUlEQVR42u3dQa7iQBBEQd//0nACdsiqzIwn/a3HuLvCm4Z5PpIU0uMRSAKWJAFLErAkCViSBCxJwJIkYEkSsCQBS5KAJUnAkgQsSQKWJAFLErAkCViSBCxJwPpxseeZ/fvX81lei3PD8eK6L88FsIBlcwMLWBYGWMACFrCABSxgAQtYwAIWsIBlYYAFLGABC1jAAhawgAUsYAELWBYGWMAyF8VgJXZt416758T90/qiSpwLYAELWMACloUBFrCABSxgAQtYwAIWsIAFLGBZGGABC1jAAhawgAUsYAELWMACloUBFrDMxThYTqg7WX7hcyW+PNohBhawgAUsYFkYYAHLXAALWMACFrCABSxgAQtYFgZYwDIXwAIWsIAFLGABC1jAApaFARawzAWw5k6Ntw5S4svDXAALWMACFrCABRpgAQtYFsZ1gGUugAUsYAELWMACDbCABSwL4zrAMhfAAhawgAUsYIEGWMACloWpGZLWdfd8gGVhgAUscwEsYLlnYAELWMAClucDLAsDLGCZC2ABC1jAApaFARawPB9gWRhgActcAAtYwAIWsI4szLWufa7EwW7FyFwAy8IAC1jmAljAAhawgGVhgAUscwEsYAELWOYCWMACFrCAZWGABSxzASxgAQtY5gJYwAIWsIA1+PfmArtO33WW5wJYwHIdYAHLwhhI1wEWsIDlOsACFrCA5TrAApaFMZCuAyxgAct1gAUsYAHLdYAFLAtjIF3HXISCpc2BvDa01z67gAUsYAELWAIWsIAFLAELWAIWsIAFLGAJWMASsIAFLGAJWMACFrCAJWABS5VgXTsZ3DpIy6fzrwGa+GK49nyABSxgAQtYwHI/wAIWsIAFLGABC1jAAhawgAUs9wMsYAELWMACFrCABSxgAQtYwHI/wALW+C+OLm+maxsl+ed0kwBdngtgAQtYwAKWhQEWsMwFsIAFLGABC1jAAhawgGVhgAUscwEsYAELWMACFrCABSxgWRhgActchILVuilbT8Mngr4MBLCABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrCABSxgAQtYwAIWsIAFLGABC1hxYDk5rQtAwBpYwAIWsIAFLGABC1j2M7CABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrCABSxgDT4IJ4yzoFl+Pl6KwAIWsIAFLGAZSM8HWMACFrCABSxgAQtYwAIWsAyk5wMsYAELWMACFrCABSxgAQtYBtLzAVYcWODrG9rEnwlePjGfDB+wgAUsYAELWMACFrCABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrCABSxgASsMLNBsDvYyRr5NASxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrCABSxgAQtYwAIWsIAFLGABC1jAmgUrcaO0Yn3tc13DEcT3TswDC1jAAhawgAUsYAELWMACFrCABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrAOgeV0ddamTFz35W9KLN8zsIAFLGABC1jAApZ7BhawgGX4gQUsYAELWMACFrCA5Z6BBSxgGX5gAQtYwAIWsIAFLGC554NgXXtYiRslseWT7onPJ3k/AwtYwAIWsIAFLGABC1gWGFjAAhawgAUsYAELWMAClv0MLAsMLGABC1jAAhawgAUsYAHLfj4EVuuJ3mVoWr8J0Dr87c8QWMACFrCABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrCABSxgAQtYwAIWsIB1CKxE+Foxuoba8kvITyQDC1jAAhawgAUsYAELWMACFrCABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrD8V/VBp2yXUbt2cvrav7V8yt9vugMLWMACFrCABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrCABSxgAQtYwAIWsAS+KYxa94+T7gIWsIAFLGABC1jAErCABSxgAQtYwAIWsIAFLGABS8ACFrCABSxgAQtYwAIWsIAFrFMbpfWEcesgJWKdeLI8+RQ7sIAFLGABC1jAAhawgAUsYAELWMACFrCABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAGrEKzEWjdBIvrWq++FByxgAct6AQtYBgBYwAIWsIAFLGABywAAy3oBC1gGAFjAAhawgGW9gAUsAwAs6wUsYBkAYAFrHKzEE+qtz3A537i4t3+ABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrCABSxgAQtYwAIWsIAFLGABC1jAigFr+WeCExFJnIv2lxmwgAUsYAELWMACFrCABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMACFrCABSxgAQtYVWAtQ5w4SImn4ZNfwMDyuYAFLGAZbJ8LWMACFrCABSxgAcvnAhawgGWwfS5gAQtYwAIWsIAFLJ8LWMAClsEGFrCAVQjWtZxQ74OmNWABC1jAAhawgAUsYAELWMACFrCABSxgAQtYwAIWsIAFLGABC1jAAhawgAUsYAELWMASsE6B5WRwHyKJ3wSwxzJf0sACFrDsMWDZTMACFrCAZTMByx4DFrCABSx7DFg2E7CABSxgAQtY9hiwgAUsYAELWDYTsOwfYB0CS5KAJUnAkgQsSQKWJGBJErAkCViSgCVJwJIkYEkCliQBS5KAJQlYkgQsSQKWJGBJErAkCViSyvoC+zqL2hxNBJMAAAAASUVORK5CYII="
				}
			}
	*/

	return nil
}
