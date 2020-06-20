package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/plaurent-dev/nttbm/pkg/proto/site"

	"google.golang.org/grpc"
)

type server struct {
	site.UnimplementedSiteServiceServer
}

func main() {
	log.Println("Server running ...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer()
	site.RegisterSiteServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(lis))
}

func (s *server) Site(ctx context.Context, request *site.SiteRequest) (*site.SiteResponse, error) {
	log.Println(fmt.Sprintf("Request: %g", request.GetUrl()))
	access := true
	urlStr := request.GetUrl()
	fmt.Printf("Get Informations: %v\n", urlStr)

	start := time.Now()

	//creating the proxyURL
	proxyStr := "http://proxy.sbe-online.com:3128"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
		access = false
	}

	//creating the URL to be loaded through the proxy
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		access = false
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	hrequest, herr := http.NewRequest("GET", url.String(), nil)
	if herr != nil {
		log.Println(herr)
		access = false
		fmt.Println("ACCES KO : NOT FOUND")
		elapsed := time.Since(start)
		log.Printf("matb test took %s", elapsed)
		return &site.SiteResponse{Site: urlStr, Access: access, Elapsedtime: elapsed.Milliseconds()}, nil
	}
	fmt.Println(" * REQUEST IN PROGRESS * ")
	//calling the URL
	response, cerr := client.Do(hrequest)
	if cerr != nil {
		log.Println(err)
		access = false
		if strings.Contains(cerr.Error(), "Forbidden") {
			fmt.Println("ACCES KO : Forbidden")
		} else {
			fmt.Println("ACCES KO : UNDETERMINED")
		}
		elapsed := time.Since(start)
		log.Printf("matb test took %s", elapsed)
		return &site.SiteResponse{Site: urlStr, Access: access, Elapsedtime: elapsed.Milliseconds()}, nil
	}

	//getting the response
	fmt.Println(" * LECTURE PAGE * ")
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		access = false
		fmt.Println("ACCES KO : PROBLEM READING SITE")
		elapsed := time.Since(start)
		log.Printf("matb test took %s", elapsed)
		return &site.SiteResponse{Site: urlStr, Access: access, Elapsedtime: elapsed.Milliseconds()}, nil
	}

	arr := make([]string, 2)
	arr[0] = "google"
	arr[1] = "accounts"
	res := findAllOccurrences([]byte(data), arr)
	//fmt.Println(res)
	if len(res) > 0 {
		fmt.Println("ACCES OK")
	} else {
		access = false
	}

	elapsed := time.Since(start)
	log.Printf("matb test took %s", elapsed)
	log.Printf("_____________________________________")

	return &site.SiteResponse{Site: urlStr, Access: access, Elapsedtime: elapsed.Milliseconds()}, nil
}

func findAllOccurrences(data []byte, searches []string) map[string][]int {
	results := make(map[string][]int, 0)

	for _, search := range searches {
		index := len(data)
		tmp := data
		for true {
			match := bytes.LastIndex(tmp[0:index], []byte(search))
			if match == -1 {
				break
			} else {
				index = match
				results[search] = append(results[search], match)
			}
		}
	}

	return results
}
