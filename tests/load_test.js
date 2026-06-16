import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 200,
  duration: '1m',

  thresholds: {
    http_req_duration: ['p(95)<100'],
    checks: ['rate>0.99'],
  },
};

export default function () {
  const ip = `10.${__VU % 255}.${__ITER % 255}.${(__VU + __ITER) % 255}`;

  const res = http.get('http://localhost:8080/', {
    headers: {
      'X-Forwarded-For': ip,
    },
  });

  check(res, {
    'status is 200 or 429': (r) =>
      r.status === 200 || r.status === 429,
  });
}