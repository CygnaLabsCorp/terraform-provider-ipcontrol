package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	// "reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"

	cc "terraform-provider-ipcontrol/ipcontrol/entities"
)

type HostConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type TransportConfig struct {
	SslVerify          bool
	certPool           *x509.CertPool // not exposed
	HttpRequestTimeout time.Duration  // in seconds
}

func NewTransportConfig(sslVerify string, httpRequestTimeout int) (cfg TransportConfig) {
	switch {
	case "false" == strings.ToLower(sslVerify):
		cfg.SslVerify = false
	case "true" == strings.ToLower(sslVerify):
		cfg.SslVerify = true
	default:
		caPool := x509.NewCertPool()
		cert, err := ioutil.ReadFile(sslVerify)
		if err != nil {
			log.Printf("Cannot load certificate file '%s'", sslVerify)
			return
		}
		if !caPool.AppendCertsFromPEM(cert) {
			err = fmt.Errorf("Cannot append certificate from file '%s'", sslVerify)
			return
		}
		cfg.certPool = caPool
		cfg.SslVerify = true
	}

	cfg.HttpRequestTimeout = time.Duration(httpRequestTimeout)
	return
}

type HttpRequestBuilder interface {
	Init(HostConfig)
	BuildUrl(r RequestType, obj cc.IpamObject, ref string) (urlStr string)
	BuildBody(r RequestType, obj cc.IpamObject) (jsonStr []byte)
	BuildRequest(r RequestType, obj cc.IpamObject, ref string) (req *http.Request, err error)
}

type HttpRequestor interface {
	Init(TransportConfig)
	SendRequest(*http.Request) ([]byte, error)
}

type CaaRequestBuilder struct {
	HostConfig HostConfig
}

type CaaHttpRequestor struct {
	client http.Client
}

type CAAConnector interface {
	CreateObject(obj cc.IpamObject, ref string) (id string, err error)
	GetObject(obj cc.IpamObject, ref string, res interface{}) error
	ExportObjects(obj cc.IpamObject, res interface{}) (err error)
	DeleteObject(obj cc.IpamObject, ref string) (refRes string, err error)
	UpdateObject(obj cc.IpamObject, ref string) (refRes string, err error)
}

type Connector struct {
	HostConfig      HostConfig
	TransportConfig TransportConfig
	RequestBuilder  HttpRequestBuilder
	Requestor       HttpRequestor
}

type RequestType int

const (
	CREATE RequestType = iota
	GET
	DELETE
	UPDATE
	EXPORT
	LOGIN
)

func (r RequestType) toMethod() string {
	switch r {
	case CREATE:
		return "POST"
	case GET:
		return "GET"
	case DELETE:
		return "DELETE"
	case UPDATE:
		return "PUT"
	case EXPORT:
		return "POST"
	case LOGIN:
		return "POST"
	}

	return ""
}

func getHTTPResponseError(resp *http.Response) error {
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	msg := fmt.Sprintf("CAA request error: %d('%s')\nContents:\n%s\n", resp.StatusCode, resp.Status, content)
	log.Printf(msg)
	return errors.New(msg)
}

func (whr *CaaHttpRequestor) Init(cfg TransportConfig) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !cfg.SslVerify,
			RootCAs:       cfg.certPool,
			Renegotiation: tls.RenegotiateOnceAsClient},
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	whr.client = http.Client{Jar: jar, Transport: tr, Timeout: cfg.HttpRequestTimeout * time.Second}
}

func (whr *CaaHttpRequestor) SendRequest(req *http.Request) (res []byte, err error) {
	var resp *http.Response
	resp, err = whr.client.Do(req)
	if err != nil {
		return
	} else if !(resp.StatusCode == http.StatusOK ||
		(resp.StatusCode == http.StatusCreated &&
			req.Method == RequestType(CREATE).toMethod())) {
		err := getHTTPResponseError(resp)
		return nil, err
	}
	defer resp.Body.Close()
	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Http Reponse ioutil.ReadAll() Error: '%s'", err)
		return
	}

	return
}

func (wrb *CaaRequestBuilder) Init(cfg HostConfig) {
	wrb.HostConfig = cfg
}

func (wrb *CaaRequestBuilder) BuildUrl(t RequestType, obj cc.IpamObject, ref string) (urlStr string) {
	path := []string{"workflow"}
	if len(ref) > 0 {
		path = append(path, ref)
	}

	var objJSON []byte
	var err error
	objJSON, err = json.Marshal(obj)
	if err != nil {
		log.Printf("Cannot marshal object '%s': %s", obj, err)
		// return path
	}
	var dataMap map[string]interface{}
	if err := json.Unmarshal(objJSON, &dataMap); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	vals := url.Values{}
	for key, value := range dataMap {
		vals.Set(key, fmt.Sprintf("%v", value))
	}

	qry := ""
	if t == GET {
		qry = vals.Encode()
	}

	u := url.URL{
		Scheme:   "https",
		Host:     wrb.HostConfig.Host + ":" + wrb.HostConfig.Port,
		Path:     strings.Join(path, "/"),
		RawQuery: qry,
	}

	return u.String()
}

func (wrb *CaaRequestBuilder) BuildBody(t RequestType, obj cc.IpamObject) []byte {
	var objJSON []byte
	var err error

	objJSON, err = json.Marshal(obj)
	if err != nil {
		log.Printf("Cannot marshal object '%s': %s", obj, err)
		return nil
	}

	log.Println("[DEBUG] BuildBody objJSON: " + fmt.Sprintln(string(objJSON)))

	// // append the 'Params' object which includes the selectors attributes
	// params := obj.Params()

	// log.Println("[DEBUG] BuildBody params: " + fmt.Sprintf("%v", params))

	// if len(params) > 0 {
	// 	paramsJSON, err := json.Marshal(params)

	// 	if err != nil {
	// 		log.Printf("Cannot marshal Search attributes. '%s'\n", err)
	// 		return nil
	// 	}

	// 	log.Printf("[DEBUG] len (objJSON): %v", len(objJSON))

	// 	if len(objJSON) > 2 { // if it's empty it's '{}'
	// 		objJSON = append(append(objJSON[:len(objJSON)-1], byte(',')), paramsJSON[1:]...)
	// 	} else {
	// 		// if it's empty = {}, then just assign the full paramsJSON to objJSON, the append above will generate a ',' empty beginning element in the json otherwise
	// 		objJSON = paramsJSON
	// 	}

	// 	// this should shows that the params obj was appended to the body
	// 	log.Println("[DEBUG] BuildBody paramsJSON: " + fmt.Sprintln(string(paramsJSON)))
	// }

	return objJSON
}

func (wrb *CaaRequestBuilder) BuildRequest(t RequestType, obj cc.IpamObject, ref string) (req *http.Request, err error) {

	urlStr := wrb.BuildUrl(t, obj, ref)

	var bodyStr []byte
	if obj != nil {
		bodyStr = wrb.BuildBody(t, obj)
	}

	log.Println("[DEBUG] BuildRequest bodyStr: " + fmt.Sprintf(string(bodyStr)))

	req, err = http.NewRequest(t.toMethod(), urlStr, bytes.NewBuffer(bodyStr))
	if err != nil {
		log.Printf("err1: '%s'", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(wrb.HostConfig.Username, wrb.HostConfig.Password)

	return
}

func (c *Connector) makeRequest(t RequestType, obj cc.IpamObject, ref string) (res []byte, err error) {
	var req *http.Request
	req, err = c.RequestBuilder.BuildRequest(t, obj, ref)
	res, err = c.Requestor.SendRequest(req)

	// tries twice ...
	if err != nil {
		req, err = c.RequestBuilder.BuildRequest(t, obj, ref)
		res, err = c.Requestor.SendRequest(req)
	}

	return
}

func NewConnector(hostConfig HostConfig, transportConfig TransportConfig,
	requestBuilder HttpRequestBuilder, requestor HttpRequestor) (res *Connector, err error) {
	res = nil

	connector := &Connector{
		HostConfig:      hostConfig,
		TransportConfig: transportConfig,
	}

	connector.RequestBuilder = requestBuilder
	connector.RequestBuilder.Init(connector.HostConfig)

	connector.Requestor = requestor
	connector.Requestor.Init(connector.TransportConfig)

	res = connector
	return
}

// -----------------------------

/* Just the ID is produced as output
 * then TF getSubnet should call getSubnetById to retrieve it at the end of the create execution to return the block information
 */
func (c *Connector) CreateObject(obj cc.IpamObject, ref string) (id string, err error) {

	resp, err := c.makeRequest(CREATE, obj, ref)
	if err != nil || len(resp) == 0 {
		log.Printf("CreateObject request error: '%s'\n", err)
		return
	}

	// expects a string literal as result
	// so in case not provided in the response from the CAA just append them before being unamashalled
	s := string(resp[:])
	if !strings.HasPrefix(s, "\"") {
		s = strconv.Quote(s)
	}
	b := []byte(s)

	err = json.Unmarshal(b, &id)
	if err != nil {
		log.Printf("CreateObject Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return
	}

	return
}

/* the GetObject expects a JS object as res interface{} */
func (c *Connector) GetObject(obj cc.IpamObject, ref string, res interface{}) (err error) {
	resp, err := c.makeRequest(GET, obj, ref)
	if err != nil {
		log.Printf("GetObject request error: '%s'\n", err)
		return err
	}

	//to check empty underlying value of interface
	err = json.Unmarshal(resp, res)
	if err != nil {
		log.Printf("GetObject Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return err
	}

	return
}

func (c *Connector) DeleteObject(obj cc.IpamObject, ref string) (refRes string, err error) {
	refRes = ""
	resp, err := c.makeRequest(DELETE, obj, ref)
	if err != nil {
		log.Printf("DeleteObject request error: '%s'\n", err)
		return
	}
	refRes = string(resp)

	return
}

func (c *Connector) UpdateObject(obj cc.IpamObject, ref string) (refRes string, err error) {
	refRes = ""
	resp, err := c.makeRequest(UPDATE, obj, ref)
	if err != nil {
		log.Printf("Failed to update object %s: %s", obj.ObjectType(), err)
		return
	}

	// expects a string literal as result
	// so in case not provided in the response from the CAA just append them before being unamashalled
	s := string(resp[:])
	if !strings.HasPrefix(s, "\"") {
		s = strconv.Quote(s)
	}
	b := []byte(s)

	err = json.Unmarshal(b, &refRes)
	if err != nil {
		log.Printf("Cannot unmarshall update object response'%s', err: '%s'\n", string(resp), err)
		return
	}
	return
}

/*
store params (in addition to objType) into obj

	return an array of ipmaObjects
*/
func (c *Connector) ExportObjects(obj cc.IpamObject, res interface{}) (err error) {

	// API End point will become: https://<ip>:1880/workflow/tf/export/<objType>

	resp, err := c.makeRequest(EXPORT, obj, "/ipcaddsubnet")
	if err != nil {
		log.Printf("ExportObjects request error: '%s'\n", err)
		return err
	}

	err = json.Unmarshal(resp, res)
	if err != nil {
		log.Printf("ExportObjects: Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return err
	}

	if len(resp) == 0 {
		return
	}

	log.Printf("[DEBUG] ExportObjects JSON unmarshalled '%s'", string(resp))

	return
}

func RunDebug() {
	hostConfig := HostConfig{
		Host:     "192.168.89.155",
		Port:     "1880",
		Username: "incadmin",
		Password: "incadmin",
	}
	requestBuilder := CaaRequestBuilder{}
	requestor := CaaHttpRequestor{}
	transportConfig := TransportConfig{}
	connector, _ := NewConnector(hostConfig, transportConfig, &requestBuilder, &requestor)

	objMgr := new(ObjectManager)
	objMgr.connector = connector

	//result, err := objMgr.CreateSubnet("incadmin", "incadmin", "/InControl/phong", "138.0.0.0", "Any", "24")
	result, err := objMgr.GetIPAddress("23.0.0.2", "/InControl/phong")
	fmt.Print(result, err)
}
