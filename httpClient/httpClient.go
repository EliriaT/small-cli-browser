package httpClient

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"hi/cache"
	"io"
	"jaytaylor.com/html2text"
	"log"
	"net"
	"net/textproto"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	htmlContentDef  = "text/html;"
	plainContentDef = "text/plain;"
	jsonContentDef  = "application/json;"
	CacheBucket     = "cacheBucket"
)

type Client struct {
	Url *url.URL
}

func (c *Client) CheckIfInCache() (bool, cache.CacheValue) {
	cachedS := cache.CacheValue{}

	cached := cache.Get([]byte(CacheBucket), []byte(c.Url.String()))
	json.Unmarshal(cached, &cachedS)

	if cachedS.Etag == "" {
		return false, cachedS
	}
	return true, cachedS
}

func (c *Client) MakeHTTPRequest(count int) (textproto.MIMEHeader, string) {
	present, cachedS := c.CheckIfInCache()
	var request string
	if present {
		request = fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\nIf-None-Match: %s\r\nAccept: text/plain,text/html,application/json;q=0.8,*/*;q=0.8\r\n\r\n", c.Url.RequestURI(), c.Url.Host, cachedS.Etag)

	} else {
		request = fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\nAccept: text/plain,text/html,application/json;q=0.8,*/*;q=0.8\r\n\r\n", c.Url.RequestURI(), c.Url.Host)

	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", c.Url.Host+":443")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	TLSConn := tls.Client(conn, &tls.Config{
		ServerName: c.Url.Hostname(),
	})
	defer TLSConn.Close()

	err = TLSConn.Handshake()
	checkError(err)

	_, err = TLSConn.Write([]byte(request))
	checkError(err)

	result, err := io.ReadAll(TLSConn)
	checkError(err)

	headers, statusCode, body := c.parseHTTPResponse(string(result))

	if statusCode >= 300 && statusCode <= 303 || statusCode == 307 || statusCode == 308 {
		if count == 10 {
			fmt.Println("Sorry this  address can not be accessed due to infinite redirection")
			os.Exit(0)
		}
		redirectURL := headers.Get("Location")
		urlStruct, err := url.Parse(redirectURL)
		checkError(err)
		c.Url = urlStruct
		newHeader, newBody := c.MakeHTTPRequest(count + 1)
		return newHeader, newBody
	}
	if statusCode == 304 {
		headers.Add("Content-Type", cachedS.ContentType)
		return headers, cachedS.Content
	}

	return headers, body

}

func (c *Client) HandleCache(headers textproto.MIMEHeader, body string) {
	//check etag field
	// if exist cache in file [url] = {etag, content}
	// on each request, send if-none-match
	// if response if 304, print cache, else parse body

	cacheHeaders := headers["Cache-Control"]

	// checking  wether we must not to store any version of the resource under any circumstances;
	for _, v := range cacheHeaders {
		if v == "no-store" {
			return
		}
	}
	// no etag value
	if headers.Get("Etag") == "" {
		return
	}

	contentType := strings.Split(headers.Get("Content-Type"), " ")[0]
	etagValue := headers["Etag"][0]

	byteStruct, _ := json.Marshal(cache.CacheValue{Etag: etagValue, Content: body, ContentType: contentType})
	err := cache.Put([]byte(CacheBucket), []byte(c.Url.String()), byteStruct)
	if err != nil {
		log.Fatal(err)
	}

}

func (c *Client) HandleBody(headers textproto.MIMEHeader, body string) {
	contentType := strings.Split(headers.Get("Content-Type"), " ")[0]

	//because of the ; at the end
	alternativeContentTypeConst := contentType + ";"

	if contentType == htmlContentDef || alternativeContentTypeConst == htmlContentDef {
		c.printHTMLBody(body)
	} else if contentType == plainContentDef || alternativeContentTypeConst == plainContentDef {
		c.printPlainTextBody(body)
	} else if contentType == jsonContentDef || alternativeContentTypeConst == jsonContentDef {
		c.printJsonBody(body)
	}
}
func (c *Client) printPlainTextBody(body string) {
	width := 140
	chunks := splitString(body, width)
	formatted := strings.Join(chunks, "\n")
	formatted = strings.ReplaceAll(formatted, "\n", "\n\t\t\t\t\t")
	//formatted = justify.Center(100, formatted)
	fmt.Println("\n\t\t\t\t\t" + formatted)

}
func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func (c *Client) printJsonBody(body string) {
	list := strings.Split(body, "\r\n")
	val, err := PrettyString(list[1])
	if err != nil {
		log.Fatal("Sorry, error beautifying json response")
	}
	fmt.Println(val)
}
func (c *Client) printHTMLBody(body string) {
	plain, err := html2text.FromString(body, html2text.Options{PrettyTables: true, OmitLinks: true, TextOnly: true})
	checkError(err)
	width := 140
	chunks := splitString(plain, width)
	formatted := strings.Join(chunks, "\n")
	formatted = strings.ReplaceAll(formatted, "\n", "\n\t\t\t\t\t")
	//formatted = justify.Center(100, formatted)
	fmt.Println("\n\t\t\t\t\t" + formatted)
}

func splitString(s string, width int) []string {
	var chunks []string

	for len(s) > 0 {
		// If the string is shorter than the maximum width, add it as is
		if len(s) <= width {
			chunks = append(chunks, s)
			break
		}

		// Find the maximum index to split the string at
		maxIndex := strings.LastIndex(s[:width+1], " ")

		// If there are no spaces in the first `width` characters, split at `width`
		if maxIndex == -1 {
			maxIndex = width
		}

		// Add the chunk to the list
		chunks = append(chunks, s[:maxIndex])

		// Remove the processed chunk from the string
		s = s[maxIndex+1:]
	}

	return chunks
}

func (c *Client) parseHTTPResponse(response string) (mapHeaders textproto.MIMEHeader, statusCode int, body string) {

	header, body, _ := strings.Cut(response, "\r\n\r\n")

	header = "First-Header: " + header

	reader := bufio.NewReader(strings.NewReader(header + "\r\n\r\n"))
	tp := textproto.NewReader(reader)

	mapHeaders, err := tp.ReadMIMEHeader()

	checkError(err)

	firstHeader := mapHeaders.Get("First-Header")

	listFirstHeader := strings.Split(firstHeader, " ")
	statusCode, err = strconv.Atoi(listFirstHeader[1])
	checkError(err)

	return mapHeaders, statusCode, body
}

func NewClient(urlS string) Client {
	if !strings.HasPrefix(urlS, "https://") && !strings.HasPrefix(urlS, "http://") {
		urlS = "https://" + urlS
	}
	urlStruct, err := url.Parse(urlS)
	checkError(err)

	if err != nil {
		log.Fatal("Can't init boltDB")
	}

	return Client{Url: urlStruct}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
