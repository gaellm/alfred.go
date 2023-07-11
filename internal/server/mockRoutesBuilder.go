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

package server

import (
	"alfred/internal/action"
	"alfred/internal/function"
	"alfred/internal/helper"
	"alfred/internal/log"
	"alfred/internal/mock"
	"alfred/internal/tracing"
	"alfred/pkg/detachcontext"
	"alfred/pkg/request"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func AddMocksRoutes(mux *http.ServeMux, mockCollection mock.MockCollection, functions function.FunctionCollection, alfredGlobalDelay *time.Duration) {

	ctx := context.Background()
	for _, m := range mockCollection.Mocks {

		m := m

		log.Debug(ctx, "Creating route for mock '"+m.GetName()+"'", zap.String("mock-url", m.GetRequestUrl()), zap.String("mock-conf", string(m.GetJsonBytes())))

		mux.HandleFunc(m.GetRequestUrl(), func(w http.ResponseWriter, r *http.Request) {

			requestRecover(w, r)
			ctx := r.Context()

			span := tracing.GetSpanFromContext(ctx)
			span.SetAttributes(attribute.String("mockUsed", m.GetName()))
			tracer := span.TracerProvider().Tracer(tracing.TracerName, trace.WithInstrumentationVersion(tracing.TracerVersion))

			helpersPopulated := []helper.Helper{}
			var req request.Req
			var res request.Res

			ctxReqDetailsSpan, reqDetailsSpan := tracer.Start(ctx, "get request details")
			data, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error(ctxReqDetailsSpan, "failed to read request body", err,
					zap.String("mock-name", m.GetName()),
					zap.String("request-path", r.RequestURI),
				)
			}

			//req
			{
				req.Body = string(data)
				req.Method = r.Method
				req.SetHeaders(r.Header)
				req.Url = r.RequestURI
				req.SetQuery(r.URL.Query())
			}

			reqDetailsStr, _ := json.Marshal(req)
			span.SetAttributes(attribute.String("requestDetails", string(reqDetailsStr)))

			reqDetailsSpan.End()

			log.Debug(ctx, "received a mock request, gona use mock '"+m.GetName()+"'",
				zap.String("request-details", string(reqDetailsStr)),
				zap.String("mock-conf", string(m.GetJsonBytes())),
			)

			res.Body = m.GetResponseBody()

			if m.HasHelper() {
				span.SetAttributes(attribute.Bool("useHelper", true))
				ctxHelper, helperSpan := tracer.Start(ctx, "manage helper(s)")

				if m.HasRequestHelper() {

					ctxReqHelperSpan, reqHelperSpan := tracer.Start(ctxHelper, "populate request helper(s)")

					log.Debug(ctxReqHelperSpan, "start to populate request helper(s)",
						zap.String("mock-name", m.GetName()),
						zap.String("request-details", string(reqDetailsStr)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
					)

					// Populate mock request helpers
					pathHelpersPopulated, err := helper.PathHelperWatcher(r, m.GetPathRegexHelpers())
					if err != nil {
						log.Warn(ctxReqHelperSpan, "helpers path watcher in error", err,
							zap.String("mock-name", m.GetName()),
							zap.String("request-details", string(reqDetailsStr)),
							zap.String("mock-conf", string(m.GetJsonBytes())),
						)
					}

					helpersPopulated = append(helpersPopulated, pathHelpersPopulated...)

					// Populate mock request helpers
					requestHelpersPopulated, err := helper.RequestHelperWatcher([]byte(req.Body), r, m.GetRequestHelpers())
					if err != nil {
						log.Warn(ctxReqHelperSpan, "helpers request watcher in error", err,
							zap.String("mock-name", m.GetName()),
							zap.String("request-details", string(reqDetailsStr)),
							zap.String("mock-conf", string(m.GetJsonBytes())),
						)
					}

					helpersPopulated = append(helpersPopulated, requestHelpersPopulated...)

					log.Debug(ctxReqHelperSpan, "request helper(s) populated",
						zap.String("mock-name", m.GetName()),
						zap.String("request-details", string(reqDetailsStr)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
						zap.String("helpers", helper.StringifyHelpers(helpersPopulated)),
					)

					reqHelperSpan.SetAttributes(attribute.String("helpers", helper.StringifyHelpers(requestHelpersPopulated)))
					reqHelperSpan.End()
				}

				if m.HasDatetHelper() {

					ctxDateHelperSpan, dateHelperSpan := tracer.Start(ctxHelper, "populate date helper(s)")

					log.Debug(ctxDateHelperSpan, "start to populate date helper(s)",
						zap.String("mock-name", m.GetName()),
						zap.String("request-details", string(reqDetailsStr)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
					)

					// Populate mock date helpers
					dateHelpersPopulated, err := helper.DateWatcher(m.GetDateHelpers())
					if err != nil {
						log.Warn(ctxDateHelperSpan, "helpers date watcher in error", err,
							zap.String("mock-name", m.GetName()),
							zap.String("request-details", string(reqDetailsStr)),
							zap.String("mock-conf", string(m.GetJsonBytes())),
						)
					}

					helpersPopulated = append(helpersPopulated, dateHelpersPopulated...)

					log.Debug(ctxDateHelperSpan, "date helper(s) populated",
						zap.String("mock-name", m.GetName()),
						zap.String("request-details", string(reqDetailsStr)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
						zap.String("helpers", helper.StringifyHelpers(helpersPopulated)),
					)

					dateHelperSpan.SetAttributes(attribute.String("helpers", helper.StringifyHelpers(dateHelpersPopulated)))
					dateHelperSpan.End()
				}

				if m.HasRandomHelper() {

					ctxRandomHelperSpan, randomHelperSpan := tracer.Start(ctxHelper, "populate random helper(s)")

					log.Debug(ctxRandomHelperSpan, "start to populate random helper(s)",
						zap.String("mock-name", m.GetName()),
						zap.String("request-details", string(reqDetailsStr)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
					)

					// Populate mock random helpers
					randomHelpersPopulated, err := helper.RandomWatcher(m.GetRandomHelpers())
					if err != nil {
						log.Warn(ctxRandomHelperSpan, "helpers random watcher in error", err,
							zap.String("mock-name", m.GetName()),
							zap.String("request-details", string(reqDetailsStr)),
							zap.String("mock-conf", string(m.GetJsonBytes())),
						)
					}

					helpersPopulated = append(helpersPopulated, randomHelpersPopulated...)

					log.Debug(ctxRandomHelperSpan, "random helper(s) populated",
						zap.String("mock-name", m.GetName()),
						zap.String("request-details", string(reqDetailsStr)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
						zap.String("helpers", helper.StringifyHelpers(helpersPopulated)),
					)

					randomHelperSpan.SetAttributes(attribute.String("helpers", helper.StringifyHelpers(randomHelpersPopulated)))
					randomHelperSpan.End()
				}

				//function JS
				if m.HasFunctionFile() {
					span.SetAttributes(attribute.Bool("useJsFunction", true))

					ctxFuncFileHelperSpan, funcFileHelperSpan := tracer.Start(ctxHelper, "helper updater javascript function")
					funcFileHelperSpan.SetAttributes(attribute.String("helpersBefore", helper.StringifyHelpers(helpersPopulated)))

					f, _ := functions.GetFunction(m.FunctionFile)
					if f.HasFuncUpdateHelpers {

						helpersPopulated, err = f.UpdateHelpersListener(helpersPopulated)
						if err != nil {
							log.Error(ctxFuncFileHelperSpan, "error using user js update helper function", err)
						}
						log.Debug(ctxFuncFileHelperSpan, "update helper(s) populated with user js function",
							zap.String("mock-name", m.GetName()),
							zap.String("request-details", string(reqDetailsStr)),
							zap.String("mock-conf", string(m.GetJsonBytes())),
							zap.String("helpers", helper.StringifyHelpers(helpersPopulated)),
						)
					}
					funcFileHelperSpan.SetAttributes(attribute.String("helpersAfter", helper.StringifyHelpers(helpersPopulated)))
					funcFileHelperSpan.End()
				}

				ctxReplaceHelperSpan, replaceHelperSpan := tracer.Start(ctxHelper, "build response with helper(s) value(s)")

				// Replace helpers inside mock response body
				res.Body, err = helper.HelperReplacement(m.GetResponseBody(), helpersPopulated)

				// Set helpers inside mock response headers and set
				for k, v := range m.GetResponseHeaders() {

					v, err = helper.HelperReplacement(v, helpersPopulated)
					res.SetHeader(k, v)
				}

				if err != nil {
					log.Warn(ctxReplaceHelperSpan, "error during helpers replacement", err,
						zap.String("request-details", string(reqDetailsStr)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
						zap.String("response-body", m.GetResponseBody()),
						zap.String("helpers", helper.StringifyHelpers(helpersPopulated)))
				}

				replaceHelperSpan.End()
				helperSpan.End()

			} else {

				//set headers
				res.Headers = m.GetResponseHeaders()
			}

			res.Status = m.GetResponseStatus()

			//delay the request
			ctxDelaySpan, delaySpan := tracer.Start(ctx, "delay response")
			{
				delay := m.GetDelay() + *alfredGlobalDelay

				log.Debug(ctxDelaySpan, "request delayed with an offset of "+fmt.Sprint(int64(delay/time.Millisecond))+" millisecond(s)",
					zap.String("mock-name", m.GetName()),
					zap.String("request-details", string(reqDetailsStr)),
					zap.String("mock-conf", string(m.GetJsonBytes())),
				)
				time.Sleep(delay)
			}
			delaySpan.End()

			//function JS
			if m.HasFunctionFile() {

				span.SetAttributes(attribute.Bool("useJsFunction", true))
				ctxAlfredJsFuncSpan, alfredJsFuncSpan := tracer.Start(ctx, "alfred javascript function")

				f, _ := functions.GetFunction(m.FunctionFile)
				if f.HasFuncAlfred {

					res, err = f.AlfredFunc(*m, helpersPopulated, req, res)
					if err != nil {
						log.Error(ctx, "error using user js alfred function", err)
					}
					log.Debug(ctxAlfredJsFuncSpan, "use user js alfred function",
						zap.String("mock-name", m.GetName()),
						zap.String("request-details", string(reqDetailsStr)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
						zap.String("helpers", helper.StringifyHelpers(helpersPopulated)),
						zap.String("response", res.Stringify()),
					)
				}
				alfredJsFuncSpan.End()
			}

			//set response headers
			{
				for k, v := range res.Headers {
					w.Header().Set(k, v)
				}
			}

			//req context end with c.String call, so save it for actions
			detachedCtx := detachcontext.Detach(ctx)

			//set status and body to end response
			w.WriteHeader(res.Status)
			_, err = w.Write([]byte(res.Body))
			if err != nil {
				log.Error(r.Context(), "failed to write", err)
			}

			//handle actions
			{

				for _, act := range m.GetActions() {

					//gourtouine
					go func(act mock.MockAction) {

						ctx, alfredActionsSpan := tracer.Start(detachedCtx, "action")

						if act.Type == action.SEND_REQUEST_TYPE {

							_, alfredDelayActionsSpan := tracer.Start(ctx, "delay action")

							delay := action.GetDelayDuration(act)
							log.Debug(ctx, "action delayed for "+fmt.Sprint(delay),
								zap.String("mock-name", m.GetName()),
								zap.String("request-details", string(reqDetailsStr)),
								zap.String("action-type", act.Type),
							)
							time.Sleep(delay)
							alfredDelayActionsSpan.End()

							req, err := action.CreateRequestFromMockAction(act, helpersPopulated)
							if err != nil {
								log.Error(ctx, "create request from mock action failed", err)
							}

							resp, err := req.Send(ctx)
							if err != nil {
								log.Error(ctx, "", err)
							}

							log.Debug(ctx, "action ended",
								zap.String("mock-name", m.GetName()),
								zap.String("request-details", string(reqDetailsStr)),
								zap.String("action-type", act.Type),
								zap.String("action-reqMethod", req.GetMethod()),
								zap.String("action-reqHeaders", fmt.Sprint(req.Headers)),
								zap.String("action-reqTargeturl", req.GetBaseUrl()),
								zap.String("action-reqBody", string(req.Body)),
								zap.String("action-responseStatus", resp.Status),
								zap.String("action-responseBody", resp.Body),
								zap.String("action-responseHeaders", fmt.Sprint(resp.Headers)),
							)
						}

						alfredActionsSpan.End()

					}(act)

				}
			}
		})
	}
}
