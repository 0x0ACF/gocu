/*
gocu is a curl copycat, a CLI HTTP client focused on simplicity and ease of use
Copyright (C) 2025  Andrés C.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RequestError struct {
	e error
	method,
	url string
}

func (err *RequestError) Error() string {
	return fmt.Sprintf("Error performing %s request on %s: %s", err.method, err.url, err.e)
}

type RequestInfo struct {
	Method,
	Url,
	Data string
	Headers map[string]string
}

type Response struct {
	Status string
	Data   string
}

func SendRequest(r *RequestInfo) (*Response, error) {
	switch r.Method {
	case http.MethodGet:
		return Get(r.Url, r.Headers)
	case http.MethodPost:
		return Post(r.Url, r.Headers, r.Data)
	case http.MethodPut:
		return Put(r.Url, r.Headers, r.Data)
	case http.MethodPatch:
		return Patch(r.Url, r.Headers, r.Data)
	case http.MethodDelete:
		return Delete(r.Url, r.Headers, r.Data)
	}

	return nil, errors.New("Invalid HTTP method received: " + r.Method)
}

func Get(url string, headers map[string]string) (*Response, error) {
	req := RequestInfo{
		Method:  http.MethodGet,
		Url:     url,
		Headers: headers,
	}

	return sendRequest(&req)
}

func Post(url string, headers map[string]string, data string) (*Response, error) {
	req := RequestInfo{
		Method:  http.MethodPost,
		Url:     url,
		Data:    data,
		Headers: headers,
	}

	return sendRequest(&req)
}

func Put(url string, headers map[string]string, data string) (*Response, error) {
	req := RequestInfo{
		Method:  http.MethodPut,
		Url:     url,
		Data:    data,
		Headers: headers,
	}

	return sendRequest(&req)
}

func Patch(url string, headers map[string]string, data string) (*Response, error) {
	req := RequestInfo{
		Method:  http.MethodPatch,
		Url:     url,
		Data:    data,
		Headers: headers,
	}

	return sendRequest(&req)
}

func Delete(url string, headers map[string]string, data string) (*Response, error) {
	req := RequestInfo{
		Method:  http.MethodDelete,
		Url:     url,
		Data:    data,
		Headers: headers,
	}

	return sendRequest(&req)
}

func PrettifyJson(data []byte) []byte {
	spaceIndentation := "  "

	var result bytes.Buffer
	json.Indent(&result, data, "", spaceIndentation)

	return result.Bytes()
}

func sendRequest(r *RequestInfo) (*Response, error) {
	client := &http.Client{}
	req, err := setupRequest(r)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, &RequestError{e: err, url: r.Url, method: http.MethodPost}
	}

	return parseResponse(res)
}

func setupRequest(r *RequestInfo) (*http.Request, error) {
	var bodyReader io.Reader = nil

	if r.Data != "" {
		bodyReader = bytes.NewReader([]byte(r.Data))
	}

	req, err := http.NewRequest(r.Method, r.Url, bodyReader)

	if err != nil {
		return nil, fmt.Errorf("Couldn't create request: %s", err)
	}

	setRequestHeaders(req, r.Headers)

	return req, nil
}

func parseResponse(resp *http.Response) (*Response, error) {
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	prettyJson := string(PrettifyJson(data))

	result := Response{Status: resp.Status, Data: prettyJson}

	return &result, nil
}

func setRequestHeaders(r *http.Request, headers map[string]string) {
	for n, v := range headers {
		r.Header.Set(n, v)
	}
}
