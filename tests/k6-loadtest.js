import http from 'k6/http';
import { sleep, check, group } from 'k6';


export const options = {
  vus: 20,
  duration: '30s',
  thresholds: {

    // 90% of requests must finish within 1ms.
    http_req_duration: ['p(90) < 10'],
  },
};

export const baseUrl = 'http://127.0.0.1:8080'

export function setup() {

  console.info("set log level to error ...");
  const url = baseUrl + '/logger';

  const payload = JSON.stringify({"configuredLevel": "error"});

  const params = {
    headers: {
      'Content-Type': 'application/json'
    }
  };

  const req = http.post(url, payload, params);

  check(req, {
    'response code was 200': (res) => res.status == 200,
    'effective log level check': (res) => JSON.parse(res.body)["effectiveLevel"] === "error"
  });
}



export default function () {

  group('json-post', function() {
  
    const jsonTestVal = Math.random().toString(36).slice(2, 7);
    const pathTestVal = Math.random().toString(36).slice(2, 7);
    const headerTestVal = Math.random().toString(36).slice(2, 7);

    const url = baseUrl + '/some/json-xml-req-helpers/' + pathTestVal;

    const payload = JSON.stringify({
      test:{
        "body-var": jsonTestVal
      }
    });

    const params = {
      headers: {
        'Content-Type': 'application/json',
        'header-var': headerTestVal
      },
      tags: {
        name: 'json-post'
      }
    };

    const req = http.post(url, payload, params);

    check(req, {
      'response code was 200': (res) => res.status == 200,
      'content from json value': (res) => JSON.parse(res.body)["from-req-body"] == jsonTestVal,
      'content from header value': (res) => res.headers["From-Req-Header"] == headerTestVal,
      'content from path value': (res) => JSON.parse(res.body)["from-req-path"] == pathTestVal
    });

    sleep(0.01);

  });
}

