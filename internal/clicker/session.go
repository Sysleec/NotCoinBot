package clicker

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/Carcraftz/cclient"

	//"crypto/tls"
	"fmt"
	"io"
	"math/rand"

	"net/url"
	"strings"
	"time"

	"compress/gzip"
	"compress/zlib"

	"github.com/andybalholm/brotli"

	http "github.com/Carcraftz/fhttp"

	tls "github.com/Carcraftz/utls"

	"golang.org/x/net/proxy"
)

type Session struct {
	client  *http.Client
	headers http.Header
}

func proxyDialer(proxyURL string) (proxy.ContextDialer, error) {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
	}

	var cntx_dialer proxy.ContextDialer = proxy.Direct

	switch parsedURL.Scheme {
	case "socks5":
		username := parsedURL.User.Username()
		passwd, _ := parsedURL.User.Password()
		auth := &proxy.Auth{
			User:     username,
			Password: passwd,
		}
		dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, auth, proxy.Direct)
		if err != nil {
			ErrorLogger.Printf("failed to create socks5 proxy dialer: %v", err.Error())
			return nil, err
		}
		cntx_dialer = dialer.(proxy.ContextDialer)
	case "http", "https":
		dialer, err := proxy.FromURL(parsedURL, proxy.Direct)
		if err != nil {
			ErrorLogger.Printf("failed to create http proxy dialer: %v", err.Error())
			return nil, err
		}
		cntx_dialer = dialer.(proxy.ContextDialer)
	default:
		return nil, fmt.Errorf("unsupported proxy scheme: %s", parsedURL.Scheme)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy dialer: %w", err)
	}
	return cntx_dialer, nil
}

type Response struct {
	body   []byte
	status int
	err    error
}

func (res *Response) String() string {
	return string(res.body)
}
func (res *Response) Error() string {
	return res.err.Error()
}

func custm_err(text string, err error) error {
	return fmt.Errorf("%v: %v", text, err.Error())
}

func CreateSession() Session {
	session := Session{}

	allowRedirect := true

	// Change JA3
	var tlsClient tls.ClientHelloID

	tlsClient = tls.HelloIOS_Auto

	client, err := cclient.NewClient(tlsClient, "", allowRedirect, time.Duration(6))
	if err != nil {
		log.Fatal(err)
	}

	session.client = &client
	session.headers = generate_headers()
	return session
}

func readBody(resp *http.Response) ([]byte, error) {
	encoding := resp.Header["Content-Encoding"]
	body, err := ioutil.ReadAll(resp.Body)

	var clearBody []byte
	finalres := ""

	if err != nil {
		panic(err)
	}
	finalres = string(body)
	if len(encoding) > 0 {
		if encoding[0] == "gzip" {
			unz, err := gUnzipData(body)
			if err != nil {
				panic(err)
			}
			clearBody = unz
			finalres = string(unz)
		} else if encoding[0] == "deflate" {
			unz, err := enflateData(body)
			if err != nil {
				panic(err)
			}

			clearBody = unz
			finalres = string(unz)
		} else if encoding[0] == "br" {
			unz, err := unBrotliData(body)
			if err != nil {
				panic(err)
			}
			clearBody = unz
			finalres = string(unz)
		} else {
			fmt.Println("UNKNOWN ENCODING: " + encoding[0])
			clearBody = body
			finalres = string(body)
		}
	} else {
		clearBody = body
		finalres = string(body)
	}

	// fmt.Printf("RESPONSE: %v\n", finalres)
	_ = finalres

	return clearBody, nil
}

func gUnzipData(data []byte) (resData []byte, err error) {
	gz, _ := gzip.NewReader(bytes.NewReader(data))
	defer gz.Close()
	respBody, err := ioutil.ReadAll(gz)
	return respBody, err
}
func enflateData(data []byte) (resData []byte, err error) {
	zr, _ := zlib.NewReader(bytes.NewReader(data))
	defer zr.Close()
	enflated, err := ioutil.ReadAll(zr)
	return enflated, err
}
func unBrotliData(data []byte) (resData []byte, err error) {
	br := brotli.NewReader(bytes.NewReader(data))
	respBody, err := ioutil.ReadAll(br)
	return respBody, err
}

func get_reader(datastr string) *bytes.Reader {
	data := []byte(datastr)
	return bytes.NewReader(data)
}

func Check_localhost(proxy string) *tls.Config {
	if strings.Contains(proxy, "127.0.0.1") || strings.Contains(proxy, "localhost") {
		return &tls.Config{InsecureSkipVerify: true}
	} else {
		return nil
	}
}

func (session *Session) Set_proxy(proxy string) error {
	proxy = strings.Replace(proxy, "\r", "", -1)
	dialer, err := proxyDialer(proxy)
	if err != nil {
		return err
	}
	tr := &http.Transport{
		DialContext:     dialer.DialContext,
		TLSClientConfig: Check_localhost(proxy)}
	session.client.Transport = tr
	session.client.Timeout = 15 * time.Second
	return nil
}

func (session *Session) Getreq(url string) *Response {
	return session.send_req(url, "GET", nil)
}

func (session *Session) Postreq(url string, data_str string) *Response {
	return session.send_req(url, "POST", get_reader(data_str))
}

func (session *Session) Patchreq(url string, data_str string) *Response {
	return session.send_req(url, "PATCH", get_reader(data_str))
}

func (session *Session) send_req(url string, method string, reader io.Reader) *Response {
	result := &Response{}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		result.err = custm_err("Error on create request", err)
		return result
	}

	session.headers.Del("Cookie")
	req.Header = session.headers

	resp, err := session.client.Do(req)

	if err != nil {
		result.err = custm_err("Error on send requests", err)
		return result
	}

	body, err := readBody(resp)

	if err != nil {
		result.err = custm_err("Error on read result response", err)
		return result
	}
	result.body = body
	result.status = resp.StatusCode
	return result
}

func generate_agent() string {
	agents := []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_1_1; like Mac OS X) AppleWebKit/600.36 (KHTML, like Gecko) Chrome/49.0.3676.327 Mobile Safari/603.9",
		"Mozilla/5.0 (Linux; Android 6.0.1; HTC One M9 Build/MRA58K) AppleWebKit/534.42 (KHTML, like Gecko) Chrome/48.0.3842.338 Mobile Safari/601.3",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 7_7_3; like Mac OS X) AppleWebKit/603.26 (KHTML, like Gecko) Chrome/55.0.3188.205 Mobile Safari/601.5",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 7_6_2; like Mac OS X) AppleWebKit/534.40 (KHTML, like Gecko) Chrome/49.0.2438.379 Mobile Safari/600.2",
		"Mozilla/5.0 (iPad; CPU iPad OS 7_5_4 like Mac OS X) AppleWebKit/603.16 (KHTML, like Gecko) Chrome/54.0.3013.211 Mobile Safari/603.8",
		"Mozilla/5.0 (Android; Android 7.1.1; Pixel C Build/NME91E) AppleWebKit/600.28 (KHTML, like Gecko) Chrome/49.0.3014.203 Mobile Safari/602.4",
		"Mozilla/5.0 (Android; Android 4.4.1; SM-J110G Build/KTU84P) AppleWebKit/600.40 (KHTML, like Gecko) Chrome/54.0.3057.201 Mobile Safari/601.5",
		"Mozilla/5.0 (Android; Android 4.4.4; LG Optimus G Build/KRT16M) AppleWebKit/601.46 (KHTML, like Gecko) Chrome/54.0.1097.125 Mobile Safari/602.4",
		"Mozilla/5.0 (Android; Android 7.0; Nexus 7 Build/NME91E) AppleWebKit/534.10 (KHTML, like Gecko) Chrome/50.0.1471.157 Mobile Safari/536.2",
		"Mozilla/5.0 (Linux; U; Android 4.4.4; XT1070 Build/SU6-7.3) AppleWebKit/537.46 (KHTML, like Gecko) Chrome/55.0.1148.223 Mobile Safari/537.6",
		"Mozilla/5.0 (Android; Android 7.0; Pixel C Build/NME91E) AppleWebKit/603.21 (KHTML, like Gecko) Chrome/50.0.2850.139 Mobile Safari/537.0",
		"Mozilla/5.0 (iPod; CPU iPod OS 7_7_9; like Mac OS X) AppleWebKit/600.24 (KHTML, like Gecko) Chrome/49.0.1878.176 Mobile Safari/534.5",
		"Mozilla/5.0 (Linux; Android 4.4.4; Elephone P2000 Build/KTU84P) AppleWebKit/535.29 (KHTML, like Gecko) Chrome/47.0.3020.286 Mobile Safari/603.9",
		"Mozilla/5.0 (Linux; U; Android 5.1; MOTOROLA MOTO X PURE XT1575 Build/LPK23) AppleWebKit/534.45 (KHTML, like Gecko) Chrome/47.0.1851.345 Mobile Safari/600.7",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_0_9; like Mac OS X) AppleWebKit/601.22 (KHTML, like Gecko) Chrome/47.0.1645.112 Mobile Safari/535.0",
		"Mozilla/5.0 (Linux; Android 4.4.1; XT1051 Build/[KXB20.9|KXC21.5]) AppleWebKit/603.29 (KHTML, like Gecko) Chrome/55.0.1553.138 Mobile Safari/602.7",
		"Mozilla/5.0 (Android; Android 4.3.1; HTC One 801e Build/JLS36C) AppleWebKit/534.41 (KHTML, like Gecko) Chrome/54.0.1790.313 Mobile Safari/536.6",
		"Mozilla/5.0 (Linux; U; Android 4.4.1; SAMSUNG SM-N9006 Build/KOT49H) AppleWebKit/603.16 (KHTML, like Gecko) Chrome/47.0.2539.118 Mobile Safari/533.9",
		"Mozilla/5.0 (Android; Android 5.0.2; Nokia 1100 LTE Build/GRK39F) AppleWebKit/535.47 (KHTML, like Gecko) Chrome/48.0.1702.219 Mobile Safari/534.1",
		"Mozilla/5.0 (Android; Android 7.1.1; SAMSUNG GT-I9500 Build/KTU84P) AppleWebKit/535.38 (KHTML, like Gecko) Chrome/55.0.2457.147 Mobile Safari/535.0",
		"Mozilla/5.0 (iPad; CPU iPad OS 7_1_3 like Mac OS X) AppleWebKit/534.12 (KHTML, like Gecko) Chrome/53.0.1883.283 Mobile Safari/533.9",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 8_3_1; like Mac OS X) AppleWebKit/534.1 (KHTML, like Gecko) Chrome/54.0.2369.370 Mobile Safari/603.4",
		"Mozilla/5.0 (iPad; CPU iPad OS 8_5_1 like Mac OS X) AppleWebKit/537.10 (KHTML, like Gecko) Chrome/52.0.2805.398 Mobile Safari/602.0",
		"Mozilla/5.0 (iPad; CPU iPad OS 11_2_2 like Mac OS X) AppleWebKit/537.9 (KHTML, like Gecko) Chrome/47.0.2162.161 Mobile Safari/602.7",
		"Mozilla/5.0 (Android; Android 5.0.1; LG-D334 Build/LRX22G) AppleWebKit/536.16 (KHTML, like Gecko) Chrome/55.0.1464.134 Mobile Safari/537.7",
		"Mozilla/5.0 (Android; Android 4.4.4; IQ4504 Quad Build/KOT49H) AppleWebKit/603.30 (KHTML, like Gecko) Chrome/52.0.3534.240 Mobile Safari/602.0",
		"Mozilla/5.0 (iPod; CPU iPod OS 8_0_5; like Mac OS X) AppleWebKit/600.5 (KHTML, like Gecko) Chrome/54.0.1899.226 Mobile Safari/537.3",
		"Mozilla/5.0 (iPod; CPU iPod OS 10_6_7; like Mac OS X) AppleWebKit/537.30 (KHTML, like Gecko) Chrome/51.0.1128.264 Mobile Safari/600.3",
		"Mozilla/5.0 (Android; Android 5.0; SM-G830K Build/LRX22G) AppleWebKit/534.39 (KHTML, like Gecko) Chrome/52.0.2039.328 Mobile Safari/603.3",
		"Mozilla/5.0 (Android; Android 5.0; LG-D326 Build/LRX22G) AppleWebKit/537.44 (KHTML, like Gecko) Chrome/55.0.1023.234 Mobile Safari/601.7",
		"Mozilla/5.0 (iPad; CPU iPad OS 10_6_5 like Mac OS X) AppleWebKit/537.39 (KHTML, like Gecko) Chrome/50.0.2473.310 Mobile Safari/601.2",
		"Mozilla/5.0 (iPod; CPU iPod OS 10_1_6; like Mac OS X) AppleWebKit/600.37 (KHTML, like Gecko) Chrome/50.0.1356.368 Mobile Safari/534.8",
		"Mozilla/5.0 (Linux; U; Android 4.3.1; SGH-N075S Build/JSS15J) AppleWebKit/536.33 (KHTML, like Gecko) Chrome/48.0.3076.325 Mobile Safari/601.7",
		"Mozilla/5.0 (iPad; CPU iPad OS 7_1_4 like Mac OS X) AppleWebKit/601.37 (KHTML, like Gecko) Chrome/52.0.1203.209 Mobile Safari/533.2",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_4_3; like Mac OS X) AppleWebKit/600.9 (KHTML, like Gecko) Chrome/51.0.3690.162 Mobile Safari/601.9",
		"Mozilla/5.0 (Linux; Android 5.1.1; SM-G9350S Build/MMB29M) AppleWebKit/603.27 (KHTML, like Gecko) Chrome/54.0.3658.284 Mobile Safari/603.7",
		"Mozilla/5.0 (iPad; CPU iPad OS 9_7_1 like Mac OS X) AppleWebKit/603.10 (KHTML, like Gecko) Chrome/48.0.1204.394 Mobile Safari/600.7",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 8_9_4; like Mac OS X) AppleWebKit/603.47 (KHTML, like Gecko) Chrome/55.0.3291.196 Mobile Safari/603.3",
		"Mozilla/5.0 (Linux; U; Android 5.0; HTC Butterfly S 901 Build/LRX22G) AppleWebKit/533.7 (KHTML, like Gecko) Chrome/54.0.2746.242 Mobile Safari/602.5",
		"Mozilla/5.0 (iPad; CPU iPad OS 11_8_1 like Mac OS X) AppleWebKit/536.34 (KHTML, like Gecko) Chrome/53.0.2346.234 Mobile Safari/536.0",
		"Mozilla/5.0 (iPad; CPU iPad OS 9_4_2 like Mac OS X) AppleWebKit/534.48 (KHTML, like Gecko) Chrome/49.0.2645.121 Mobile Safari/535.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 8_3_2; like Mac OS X) AppleWebKit/533.39 (KHTML, like Gecko) Chrome/54.0.2375.105 Mobile Safari/603.3",
		"Mozilla/5.0 (Linux; U; Android 5.1; Nexus 9 Build/LRX22C) AppleWebKit/537.39 (KHTML, like Gecko) Chrome/55.0.2118.304 Mobile Safari/534.4",
		"Mozilla/5.0 (iPad; CPU iPad OS 10_7_2 like Mac OS X) AppleWebKit/601.6 (KHTML, like Gecko) Chrome/52.0.1037.181 Mobile Safari/536.5",
		"Mozilla/5.0 (Linux; U; Android 7.1.1; Nexus 8P Build/NME91E) AppleWebKit/535.26 (KHTML, like Gecko) Chrome/55.0.2968.158 Mobile Safari/534.8",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_4_5; like Mac OS X) AppleWebKit/536.43 (KHTML, like Gecko) Chrome/51.0.1067.110 Mobile Safari/536.4",
		"Mozilla/5.0 (Linux; Android 5.1; Nexus 7 Build/LMY48B) AppleWebKit/536.23 (KHTML, like Gecko) Chrome/52.0.1133.387 Mobile Safari/535.5",
		"Mozilla/5.0 (iPad; CPU iPad OS 8_1_1 like Mac OS X) AppleWebKit/600.47 (KHTML, like Gecko) Chrome/55.0.2336.345 Mobile Safari/601.9",
		"Mozilla/5.0 (Linux; U; Android 5.0.1; SM-T805 Build/LRX22G) AppleWebKit/533.10 (KHTML, like Gecko) Chrome/52.0.3996.299 Mobile Safari/536.1",
		"Mozilla/5.0 (Linux; U; Android 6.0.1; SM-G920S Build/MDB08I) AppleWebKit/536.47 (KHTML, like Gecko) Chrome/48.0.1817.270 Mobile Safari/601.6",
		"Mozilla/5.0 (iPod; CPU iPod OS 7_3_6; like Mac OS X) AppleWebKit/537.22 (KHTML, like Gecko) Chrome/48.0.3619.367 Mobile Safari/602.2",
	}
	randomIndex := rand.Intn(len(agents))
	agent := agents[randomIndex]
	return agent
}

func generate_headers() http.Header {
	headers := http.Header{}
	agent := generate_agent()
	headers.Set("User-Agent", agent)
	//headers.Set("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
	headers.Set("Content-Type", "application/json")
	return headers
}
