package discordapi

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/kensh1ro/willie/config"
)

const api = "https://discord.com/api"

//var proxy, _ = url.Parse("http://127.0.0.1:8080")
var client = &http.Client{}

func Route(method string, param string, form io.Reader, attachment *Attachment) []byte {
	url := api + param
	var content_type string
retry:
	if attachment != nil {
		headers := make(textproto.MIMEHeader)
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		content_type = writer.FormDataContentType()
		headers.Add("Content-Disposition", "form-data; name=file;"+"filename="+attachment.Filename)
		headers.Add("Content-Type", attachment.ContentType)
		part, _ := writer.CreatePart(headers)
		io.Copy(part, form)
		writer.Close()
		form = bytes.NewReader(body.Bytes())
	} else {
		content_type = "application/json"
	}
	req, err := http.NewRequest(method, url, form)

	if err != nil {
		fmt.Println(err.Error())
		goto retry
	}

	req.Header = http.Header{
		"Authorization": []string{config.TOKEN},
		"Content-Type":  []string{content_type},
		"User-Agent":    []string{"Willie v1.0"},
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		goto retry
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	//if method == "POST" {
	//fmt.Println("*****************************")
	//fmt.Printf("----%s----", method)
	//fmt.Println(string(body))
	//fmt.Println("*****************************")
	//}
	if err != nil {
		fmt.Println(err.Error())
		goto retry
	}
	return body
}
