package config

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/ghodss/yaml"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// Important, have to be specified
//
//

type ConfigVars struct {
	ChannelConfigPath string
	OrgId             string
	UserName          string
}

var (
	Vars ConfigVars

	Sdk     *fabsdk.FabricSDK
	SdkFile *simplejson.Json

	ChaincodeId string
	ChannelId   string
	PrivateKey  string

	User msp.SigningIdentity
	Cert []byte
)

const (
	parseFilePath = "/home/tkgoh/Sandbox/block-event-parses/"
)

func LoadConfig() {
	cwd, _ := os.Getwd()
	varsJson, err := os.Open(cwd + "/config/vars.json")
	if err != nil {
		panic(err)
	}
	defer varsJson.Close()

	varsByte, err := ioutil.ReadAll(varsJson)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(varsByte, &Vars)

	data, err := ReadFile(Vars.ChannelConfigPath)
	if err != nil {
		panic(err)
	}
	data, err = yaml.YAMLToJSON(data)
	if err != nil {
		panic(err)
	}
	SdkFile, err = simplejson.NewJson(data)
	ChannelId = GetDefaultChannelId()
	ChaincodeId = GetDefaultChaincodeId()
}

func GetDefaultChannelId() string {
	channels := SdkFile.Get("channels").MustMap()
	for k := range channels {
		return k
	}
	return ""
}

func GetDefaultChaincodeId() string {
	chaincodes := SdkFile.Get("channels").Get(ChannelId).Get("chaincodes").MustArray()
	if str, ok := chaincodes[0].(string); ok {
		return strings.Split(str, ":")[0]
	}
	return ""
}

// ReadFile reads the file named by filename and returns the contents.
func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// It's a good but not certain bet that FileInfo will tell us exactly how much to
	// read, so let's try it but be prepared for the answer to be wrong.
	var n int64 = bytes.MinRead

	if fi, err := f.Stat(); err == nil {
		if size := fi.Size() + bytes.MinRead; size > n {
			n = size
		}
	}
	return readAll(f, n)
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	var buf bytes.Buffer
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.

	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if errors, ok := e.(error); ok && errors == bytes.ErrTooLarge {
			err = errors
		} else {
			panic(e)
		}
	}()
	if int64(int(capacity)) == capacity {
		buf.Grow(int(capacity))
	}
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
