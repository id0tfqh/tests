package main

/*
* ToDo
* $ go build parser.go
* $ ./parser <URL> -t a -a href -tp -ua
* $ ./parser <URL> --tag a --attr href --tor-proxy --user-agent
* Example:
* $ ./parser ya.ru -t a -a href -tp -ua
 */

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

import "golang.org/x/net/proxy"

const (
	BUFF = 512
	TORS_PROXY = "socks5://127.0.0.1:9050"
	USER_AGENT = "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.1547.2 Safari/537.36"
)

var (
	TAG_NAME = ""
	ATR_NAME = ""

	GET_TORS_PROXY = false
	GET_USER_AGENT = false
)

func main() {
	check_args(os.Args)
	check_func(urlopen("http://" + os.Args[1]))
}

func check_func(html string) {
	if TAG_NAME != "" {
		for _, result := range get_object(html) {
			fmt.Println(result)
		}
	} else {
		fmt.Println(html)
	}
}

func check_args(args []string) {
	var (
		flag_tag bool = false
		flag_atr bool = false
	)

	switch len(args) {
	case 1:
		get_error("No URL specified")
	case 2:
		if args[1] == "-i" || args[1] == "--info" {
			get_info()
		} else {
			return
		}
	default:
		for _, value := range args[2:] {

			switch value {
			case "-tp", "--tor-proxy":
				GET_TORS_PROXY = true
				continue

			case "-ua", "--user-agent":
				GET_USER_AGENT = true
				continue

			case "-t", "--tag":
				flag_tag = true
				continue

			case "-a", "--attr":
				flag_atr = true
				continue
			}

			if flag_tag {
				TAG_NAME = value
				flag_tag = false

			} else if flag_atr {
				ATR_NAME = value
				flag_atr = false
			}
		}
	}
}

func get_info() {
	// Пример: ./parser -i
	// или: ./parser --info
	fmt.Println(
		`Modules:
    -tp || --tor-proxy	-> Turn on tor proxy
    -us || --user-agent -> Turn on user-agent
    -t  || --tag	-> Tag name
    -a  || --attr	-> Attribute name

Example:
    $ ./parser ya.ru --tag a --attr href -tp -ua`)
	os.Exit(0)
}

func urlopen(url_str string) string {
	var (
		html_page string
		buffer []byte
		lenght int
		err error
	)

	var (
		client *http.Client
		req *http.Request
		resp *http.Response
	)

	if GET_TORS_PROXY {
		torProxyUrl, err := url.Parse(TORS_PROXY)
		check_error(err)

		torDialer, err := proxy.FromURL(torProxyUrl, proxy.Direct)
		check_error(err)

		torTransport := &http.Transport{
			Dial: torDialer.Dial,
		}

		client = &http.Client{
			Transport: torTransport,
			Timeout: time.Second * 5,
		}
	} else {
		client = &http.Client{}
	}

	req, err = http.NewRequest("GET", url_str, nil)
	check_error(err)
	req.Header.Add("Accept", "text/html")

	if GET_USER_AGENT {
		req.Header.Add("User-Agent", USER_AGENT)
	}

	resp, err = client.Do(req)
	check_error(err)
	defer resp.Body.Close()

	buffer = make([]byte, BUFF)
	for {
		lenght, err = resp.Body.Read(buffer)
		html_page += string(buffer[:lenght])
		if lenght == 0 || err != nil {
			break
		}
	}

	return html_page
}

func get_object(html string) []string {
	var (
		result [][]string
		slise_result []string
		regular *regexp.Regexp
	)

	if TAG_NAME != "" {
		if ATR_NAME != "" {
			regular = regexp.MustCompile(`<` + TAG_NAME + `.?` + ATR_NAME + `\s*=\s*['"]([^\s'"]+)[\s'"]`) // <a.*?href\s*=\s*['"]([^\s'"]+)[\s'"]
			result = regular.FindAllStringSubmatch(html, -1)
			for _, slise := range result {
				slise_result = append(slise_result, slise[1])
			}
			TAG_NAME = ""
			ATR_NAME = ""
		} else {
			regular = regexp.MustCompile(`<` + TAG_NAME + `[^>]*>.+</` + TAG_NAME + `>`) // <a[^>]*>.+</a>
			result = regular.FindAllStringSubmatch(html, -1)
			for _, slise := range result {
				slise_result = append(slise_result, slise[0])
			}
			TAG_NAME = ""
		}
	}
	return slise_result

}

func check_error(err error) {
	if err != nil {
		fmt.Println("Error is:", err)
		os.Exit(1)
	}
}

func get_error(err string) {
	fmt.Println("Get error is:", err)
	os.Exit(1)
}
