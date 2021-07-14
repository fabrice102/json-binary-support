package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ugorji/go/codec"
)

var JSONHandle *codec.JsonHandle

func init() {
	JSONHandle = new(codec.JsonHandle)
	JSONHandle.ErrorIfNoField = true
	JSONHandle.ErrorIfNoArrayExpand = true
	JSONHandle.Canonical = true
	JSONHandle.RecursiveEmptyCheck = true
	JSONHandle.Indent = 2
	JSONHandle.HTMLCharsAsIs = true
	// Enabling StringToRaw encode in base64
	// JSONHandle.StringToRaw = true
}

type TestObject struct {
	Data string
}

func CodecJSONEncode(obj TestObject) ([]byte, error) {
	var output []byte
	enc := codec.NewEncoderBytes(&output, JSONHandle)

	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to encode object: %v", err)
	}
	return output, nil
}

func CodecJSONDecode(data []byte) (TestObject, error) {
	enc := codec.NewDecoderBytes(data, JSONHandle)

	var v TestObject
	err := enc.Decode(&v)
	if err != nil {
		return TestObject{}, fmt.Errorf("failed to decode object: %v", err)
	}
	return v, nil
}

func main() {
	var err error

	data := []byte{0x00, 0x01, 0x80, 0x81, 0xff}

	//// make data the binary string with all the characters from 0x00 to 0xff
	//// concretely, it is in hexadecimal: 0x000102...ff
	//data := make([]byte, 256)
	//for i := 0; i < 256; i++ {
	//	data[i] = byte(i)
	//}

	// The string version of the data as should be escaped in JSON
	var dataEscapeStrBuilder strings.Builder
	for _, c := range data {
		dataEscapeStrBuilder.WriteString(fmt.Sprintf("\\u00%02x", c))
	}

	// Original object to encode
	orig := TestObject{
		Data: string(data),
	}

	// One possible way to encode it in JSON
	origJSON := fmt.Sprintf(`{
  "Data": "%s" 
}`, dataEscapeStrBuilder.String())

	fmt.Printf("Original object:\n%#v\n\n", orig)
	fmt.Printf("Data as bytes: %v\n\n", []byte(orig.Data))
	fmt.Printf("Possible JSON:\n%s\n\n", origJSON)

	var goStdDecObj TestObject
	err = json.Unmarshal([]byte(origJSON), &goStdDecObj)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded original JSON using Go standard library:\n%#v\n", goStdDecObj)
	fmt.Printf("Data as bytes: %v\n", []byte(goStdDecObj.Data))
	fmt.Printf("Does it match original data? %v\n\n", orig.Data == goStdDecObj.Data)

	codecDecObj, err := CodecJSONDecode([]byte(origJSON))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded original JSON using Go standard library:\n%#v\n", codecDecObj)
	fmt.Printf("Data as bytes: %v\n", []byte(codecDecObj.Data))
	fmt.Printf("Does it match original data? %v\n\n", orig.Data == codecDecObj.Data)

	goStdEncObj, err := json.Marshal(orig)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encoded JSON using Go standard library:\n%s\n\n", goStdEncObj)

	goStdEncObj = []byte(`{
	"Data": "\u0000\u0001\u0080\u0081\u00ff"
}`)

	var goStdEncDecObj TestObject
	err = json.Unmarshal(goStdEncObj, &goStdEncDecObj)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded JSON above using Go standard library:\n%#v\n", goStdEncDecObj)
	fmt.Printf("Does it match original data? %v\n\n", orig.Data == goStdEncDecObj.Data)

	codecEncObj, err := CodecJSONEncode(orig)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encoded JSON using Codec library:\n%s\n\n", codecEncObj)

	codecEncDecObj, err := CodecJSONDecode(codecEncObj)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded JSON above using Codec library:\n%#v\n", codecEncDecObj)
	fmt.Printf("Does it match with original data? %v\n", orig.Data == codecEncDecObj.Data)
}
