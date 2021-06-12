package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestRune(t *testing.T) {
	str := fmt.Sprintf("%s", string(os.PathSeparator))
	fmt.Println(str)
}

func TestGet(t *testing.T) {


		url := "https://contract.u.360.cn//cpiv3//mopen_company?reg_addr=%5Cu5e7f%5Cu4e1c%5Cu7701%2520%5Cu5e7f%5Cu5dde%5Cu5e02%2520%5Cu5929%2520%5Cu6cb3%5Cu533a%5Cu68e0%5Cu4e0b%5Cu8377%5Cu5149%5Cu4e09%5Cu6a2a%5Cu8def9%5Cu53f7605%5Cu623f&qid=31194%252036994&cp_name=%5Cu5e7f%5Cu5dde%5Cu98de%5Cu62d3%5Cu7f51%5Cu7edc%5Cu79d1%5Cu6280%5Cu6709%5Cu9650%5Cu516c%2520%5Cu53f8&qq=180714868&linkman=%5Cu80e1%5Cu946b&addr=%5Cu5929%5Cu6cb3%5Cu533a%5Cu6021%5Cu7965%5Cu76%2520db%5Cu8fbe%5Cu7535%5Cu5b50%5Cu4fe1%5Cu606f%5Cu521b%5Cu65b0%5Cu56edC%5Cu5ea76%5Cu697c605&email=180714%2520868@qq.com&contect=18998319984&license_no=91440101MA5AUH6388&platform=mgame"
		url = "https://www.baidu.com"
		method := "GET"

		client := &http.Client {
			Transport: &http.Transport{
				TLSClientConfig : &tls.Config{InsecureSkipVerify: true},
			},
		}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))

}

func TestExists(t *testing.T) {

}