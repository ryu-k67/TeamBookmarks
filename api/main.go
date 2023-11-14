package main

import (
	// "fmt"
	"log"
	// "net/http"
	// "net/url"
	// "os"
	// "path"
	"strings"
	// "time"
	"context"
	"encoding/json"
	// "encoding/base64"

	"github.com/PuerkitoBio/goquery"
	// "github.com/gin-contrib/cors"
	// "github.com/gin-gonic/gin"
	// "github.com/temoto/robotstxt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
)


type GetData struct {
	Scheme string `json:"scheme"`
    Url string `json:"url"`
}

type ReturnData struct {
    Title string `json:"title"`
	Description string `json:"description"`
	Image string `json:"image"`
}



// AWS lambda関数用に再構成
// レスポンスの中身を作成
func makeReturnData(inputs string) (ReturnData, error) {
	log.Printf("inputs: %s\n", inputs)
	log.Printf(inputs)
	var event GetData
	// b64String, _ := base64.StdEncoding.DecodeString(inputs)
    // rawIn := json.RawMessage(b64String)
	// err := json.Unmarshal([]byte(inputs), &event)
	// bodyBytes, err := rawIn.MarshalJSON()
	// if err != nil {
    //     return ReturnData{
	// 		Title: "error",
	// 		Description: "error",
	// 		Image: "error",
	// 	}, err
    // }
	// in := "{\"scheme\":\"https\",\"url\":\"www.google.com/\"}"
	// log.Printf("in: %s\n", in)
	jsonMarshalErr := json.Unmarshal([]byte(inputs), &event)
	log.Printf("event: %s\n", event)
    // jsonMarshalErr := json.Unmarshal(bodyBytes, &event)
	if jsonMarshalErr != nil {
		log.Printf("json unmarshal error")
		// if err, ok := err.(*json.SyntaxError); ok {
		// 	log.Println(string(in[err.Offset]))
		// }
		return ReturnData{
			Title: "error",
			Description: "error",
			Image: "error",
		}, jsonMarshalErr
	}

	// recv_url := event.Scheme + "://" + event.Url
	// // log.Printf("scheme: %s\n", event.Scheme)
	// // log.Printf("url: %s\n", event.Url)

	recv_url := event.Url
	log.Printf("url: %s\n", event.Url)

	doc, err := goquery.NewDocument(recv_url)
	if err != nil {
		log.Printf("no document")
		return ReturnData{
			Title: "error",
			Description: "error",
			Image: "error",
		}, err
	}

	var res_title string
	var res_description string
	var res_image string

	res_title = doc.Find("title").Text()
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		property, _ := s.Attr("property")
		name, _ := s.Attr("name")
		content, _ := s.Attr("content")
		if strings.Contains(property, "description") && len(res_description) == 0 {
			res_description = content
		} else if strings.Contains(property, "image") && len(res_image) == 0 {
			res_image = content
		} else if strings.Contains(name, "description") && len(res_description) == 0 {
			res_description = content
		} else if strings.Contains(name, "image") && len(res_image) == 0 {
			res_image = content
		}
	})

	log.Printf("Title: %s\n", res_title)
	log.Printf("Description: %s\n", res_description)
	log.Printf("Image: %s\n", res_image)

	if len(res_title)==0{
		res_title = "no title"
	}
	if len(res_description)==0{
		res_description = "no description"
	}
	if len(res_image)==0{
		res_image = "no image"
	}

	return ReturnData{
		Title: res_title,
		Description: res_description,
		Image: res_image,
	}, nil
}

// 												↓ APIGatewayProxyRequest型で引数にPOSTした内容を受け取り、APIGatewayProxyResponseで変換するのが作法
func bookmarkAPIHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// テキストデータをgolang構造体に変換
	// var req ReturnData
	log.Printf("request: %s\n", request)

	log.Printf("HTTPMethod:%s\n",request.HTTPMethod)
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			Headers:map[string]string{
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "https://supporters2022-vol8.web.app",
				"Access-Control-Allow-Methods": "OPTIONS,POST",
			},
			Body: "ok",
			StatusCode: 200,
		}, nil
	}

	// log.Printf(request)
	req, err := makeReturnData(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:map[string]string{
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "https://supporters2022-vol8.web.app",
				"Access-Control-Allow-Methods": "OPTIONS,POST",
			},
			Body: err.Error(),
			StatusCode: 500,
		}, err
	}

	data, _ := json.Marshal(req)
	
	return events.APIGatewayProxyResponse{
		Headers:map[string]string{
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "https://supporters2022-vol8.web.app",
            "Access-Control-Allow-Methods": "OPTIONS,POST",
		},
		Body: string(data),
		StatusCode: 200,
	}, nil
}

// 本来のmain関数
func main(){
	lambda.Start(bookmarkAPIHandler)
}