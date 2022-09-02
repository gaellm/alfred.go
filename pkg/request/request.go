/*
 * Copyright The Alfred.go Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package request

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const DEFAULT_TIMEOUT = 60 * time.Second

type Request struct {
	method  string
	url     *url.URL
	Body    []byte
	Query   map[string]string
	Headers map[string]string
	timeout time.Duration
}

type Req struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Body    string            `json:"body"`
	Query   map[string]string `json:"query"`
	Headers map[string]string `json:"headers"`
}

func (r *Req) SetHeaders(headers http.Header) {

	r.Headers = map[string]string{}

	if len(headers) > 0 {

		for header, headerArray := range headers {

			r.Headers[header] = strings.Join(headerArray, ",")
		}
	}
}

func (r *Req) SetQuery(q url.Values) {

	r.Query = map[string]string{}

	if len(q) > 0 {

		for k, v := range q {

			r.Query[k] = strings.Join(v, ",")

		}
	}
}

func (r *Request) GetBodyBytesBuffer() *bytes.Buffer {

	if len(r.Body) == 0 {
		return bytes.NewBuffer([]byte{})
	}

	return bytes.NewBuffer(r.Body)

}

func (r *Request) Send(ctx context.Context) (Response, error) {

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	client.Timeout = r.GetTimeout()

	// add client span (dns, send, etc ...)
	//ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	req, err := http.NewRequestWithContext(ctx, r.GetMethod(), r.GetBaseUrl(), r.GetBodyBytesBuffer())
	if err != nil {
		log.Fatal(err)
	}

	// appending to existing query args
	q := req.URL.Query()
	for queryArgName, queryArgValue := range r.GetQueryArgs() {

		q.Add(queryArgName, queryArgValue)
	}

	// assign encoded query string to http request
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("errored when sending request to the server: " + err.Error())
		return Response{}, nil
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response Response

	response.Body = string(responseBody)
	response.SetHeaders(resp.Header.Clone())
	response.Status = resp.Status

	return response, nil
}

func (r *Request) SetUrl(rawUrl string) error {

	u, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return err
	}

	//get url query
	if r.Query == nil {
		r.Query = make(map[string]string)

		urlQueryArgsMap := UrlValues2Map(u.Query())

		for k, v := range urlQueryArgsMap {

			r.Query[k] = v

		}
	}

	if u.Scheme == "" || u.Host == "" {
		return errors.New("request url not valid")
	}

	r.url = u

	return nil

}

func (r *Request) GetBaseUrl() string {

	url := r.GetUrl()

	return url.Scheme + "://" + url.Host + url.Path

}

func (r *Request) GetUrl() *url.URL {

	return r.url
}

func (r *Request) SetMethod(method string) error {

	method = strings.ToUpper(method)

	methods := [9]string{
		http.MethodPost,
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodConnect,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, v := range methods {

		if v == method {
			r.method = v
			return nil
		}

	}

	return errors.New("request method " + method + " not exists")

}

func (r *Request) GetMethod() string {

	return r.method

}

func (r *Request) SetTimeout(timeoutDuration string) error {

	duration, err := time.ParseDuration(timeoutDuration)

	r.timeout = duration

	return err

}

func (r *Request) GetTimeout() time.Duration {

	if r.timeout == 0 {

		return DEFAULT_TIMEOUT
	}

	return r.timeout
}

func (r *Request) GetQueryArgs() map[string]string {

	if r.Query == nil {
		r.Query = make(map[string]string)
	}

	return r.Query
}

func UrlValues2Map(urlValues url.Values) map[string]string {

	result := make(map[string]string)

	for k, v := range urlValues {

		result[k] = v[0]

	}

	return result
}

func (r *Request) AddQueryArgs(queryArgs map[string]string) {

	if r.Query == nil {
		r.Query = make(map[string]string)
	}

	for k, v := range queryArgs {

		r.Query[k] = v
	}
}
