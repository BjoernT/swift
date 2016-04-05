//  Copyright (c) 2015 Rackspace
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
//  implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package proxyserver

import (
	"net/http"

	"github.com/openstack/swift/go/hummingbird"
)

func (server *ProxyServer) AccountGetHandler(writer http.ResponseWriter, request *http.Request) {
	vars := hummingbird.GetVars(request)
	ctx := GetProxyContext(request)
	if ctx == nil {
		hummingbird.StandardResponse(writer, 500)
		return
	}
	if ctx.Authorize != nil && !ctx.Authorize(request) {
		hummingbird.StandardResponse(writer, 401)
		return
	}
	options := map[string]string{
		"format":     request.FormValue("format"),
		"limit":      request.FormValue("limit"),
		"marker":     request.FormValue("marker"),
		"end_marker": request.FormValue("end_marker"),
		"prefix":     request.FormValue("prefix"),
		"delimiter":  request.FormValue("delimiter"),
	}
	r, headers, code := server.C.GetAccount(vars["account"], options, request.Header)
	for k := range headers {
		writer.Header().Set(k, headers.Get(k))
	}
	writer.WriteHeader(code)
	if r != nil {
		defer r.Close()
		hummingbird.Copy(r, writer)
	}
}

func (server *ProxyServer) AccountHeadHandler(writer http.ResponseWriter, request *http.Request) {
	vars := hummingbird.GetVars(request)
	ctx := GetProxyContext(request)
	if ctx == nil {
		hummingbird.StandardResponse(writer, 500)
		return
	}
	if ctx.Authorize != nil && !ctx.Authorize(request) {
		hummingbird.StandardResponse(writer, 401)
		return
	}
	headers, code := server.C.HeadAccount(vars["account"], request.Header)
	for k := range headers {
		writer.Header().Set(k, headers.Get(k))
	}
	writer.WriteHeader(code)
}

func (server *ProxyServer) AccountPutHandler(writer http.ResponseWriter, request *http.Request) {
	vars := hummingbird.GetVars(request)
	ctx := GetProxyContext(request)
	if ctx == nil {
		hummingbird.StandardResponse(writer, 500)
		return
	}
	if ctx.Authorize != nil && !ctx.Authorize(request) {
		hummingbird.StandardResponse(writer, 401)
		return
	}
	request.Header.Set("X-Timestamp", hummingbird.GetTimestamp())
	hummingbird.StandardResponse(writer, server.C.PutAccount(vars["account"], request.Header))
}

func (server *ProxyServer) AccountDeleteHandler(writer http.ResponseWriter, request *http.Request) {
	vars := hummingbird.GetVars(request)
	ctx := GetProxyContext(request)
	if ctx == nil {
		hummingbird.StandardResponse(writer, 500)
		return
	}
	if ctx.Authorize != nil && !ctx.Authorize(request) {
		hummingbird.StandardResponse(writer, 401)
		return
	}
	request.Header.Set("X-Timestamp", hummingbird.GetTimestamp())
	hummingbird.StandardResponse(writer, server.C.DeleteAccount(vars["account"], request.Header))
}
