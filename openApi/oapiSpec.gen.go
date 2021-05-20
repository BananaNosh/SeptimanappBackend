// Package Openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package Openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	externalRef0 "SeptimanappBackend/openApi/types"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xTQW/UPBD9K9Z839FNsssFcqyEoBKHivZW9WCSSTKQ2K49WRqt/N/RxGx32VaiFwSX",
	"WKOMx++9eW8PjZu8s2g5Qr1PGsh2Duo9tBibQJ7JWajh4+3ttYoYdhhU54LiAdUNeqbJWKM+GSZr1KVp",
	"vqFtQQMTjwg1xEOL9xdf8t+LgdmDhh2GmGdXxaaoIGlwHq3xBDW8KaqiAg3e8CC4oMTdT4jQI8vhPAYj",
	"6K5aqOED8vvcIZeCmZAxRKjvzolQp3raoVXOjovKU5XrFA8U1YImKBNQBeQ5WBQqJLceZgwLaLBmElrS",
	"BxpiM+BkBAw+mskL4221eadhMo80zZOU20rDRPZQbt5q4MXLELKMPQZISb8KJNkMcjS2n02PrwEqvb8C",
	"tYLjDloEDaOBe33ELvUTusiBbA8p3WsIGL2zEVf5t1UlR+Mso103YbwfqVl3UX6NQmB/8uL/ATuooShl",
	"cMzfYjHT+F95tF6Z++NhzSklfaZJZhpX552uDdWTyVZnmnFUcqjosaGOmnWreaJ38QXvXLt4NE/Ahxkj",
	"X7p2+dM0MyZ5kAK2UHOYMf0jat+YHWatRTNs1bE16UMayz216feRvGqfh3I1q6T76FVq4VyNU+c+S81f",
	"MOZLSn0+96X6TjysZc4wtatqKf0IAAD///iwFg1sBQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	pathPrefix := path.Dir(pathToFile)

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "./types/types.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
	}
	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.Swagger, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.SwaggerLoader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadSwaggerFromData(specData)
	if err != nil {
		return
	}
	return
}
