package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const youdaoAPI  ="http://openapi.youdao.com/api"

const(
	youdaoAppKey string=""
	salt string=""
	password string=""
)
type result struct {
	Query 		string
	From		string
	To			string
	Translation []string
	Explains	[]string
}
type Basic struct {
	UsPhonetic string `json:"us-phonetic"`
	UkPhonetic string `json:"uk-phonetic"`
	Explains  []string `json:"explains"`
}

type YoudaoResponse struct {
	ErrorCode string
	Query string
	Translation []string
	Basic	`json:"basic"`
}
var from=flag.String("from","EN","translation from")
var to=flag.String("to","zh_CHS","translation to")
const temp=
`------------------YoudaoTranslation-----------------------
Query:{{.Query}}		from:{{.From}}			to:{{.To}}
Translation	:{{range $item,$value :=.Translation}}
	{{$item| printf "%d."}}{{$value | printf "%s"}}
{{end}}
Explains	:	{{range $item,$value:=.Explains}}
	{{$item| printf "%d."}}{{$value | printf "%s"}}
{{end}}
------------------YoudaoTranslation-----------------------
`
func main(){
//	resp,_:=http.Get("https://www.baidu.com")
//	println(resp.Body)
	flag.Parse()
	var ret *result
	ret,_=translateViaYoudao(os.Args[1],*from,*to)
	prt,err:=template.New("prt").Parse(temp)
	if err!=nil{
		log.Fatal(err)
	}
	if err:=prt.Execute(os.Stdout,ret);err!=nil{
		log.Fatal(err)
	}
	//fmt.Println(ret)
}
func translateViaYoudao(key string,from string,to string)(*result,error){
	querySign:=md5.New()
	querySign.Write([]byte(youdaoAppKey+key+salt+password))
	queryURL:="q="+key+"&from="+from+"&to="+to+"&appKey="+youdaoAppKey+
		"&salt="+salt+"&sign="+hex.EncodeToString(querySign.Sum(nil))
	queryURL=youdaoAPI+"?"+queryURL
	resp,err:=http.Get(queryURL)
	if err!=nil{
		return nil,err
	}
	if resp.StatusCode!=http.StatusOK {
		resp.Body.Close()
		return nil,fmt.Errorf("search query failed: %s",resp.Status)
	}
	var youdaoRet YoudaoResponse
	if err:=json.NewDecoder(resp.Body).Decode(&youdaoRet);err!=nil{
		resp.Body.Close()
		return nil,err
	}
	resp.Body.Close()
	var ret =result{key,from,to,youdaoRet.Translation,youdaoRet.Explains}
	return &ret,err
}
