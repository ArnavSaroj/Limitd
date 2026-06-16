import http from 'k6/http';
import { check } from 'k6';

//this runs a multiple ip load test
export const options = {
  vus: 200,
  duration: '1m',
};

export default function () {
  // Simulate 1000 unique users
  const ip = `10.${__VU % 255}.${__ITER % 255}.${(__VU + __ITER) % 255}`;

  const res = http.get('http://localhost:8080/', {
    headers: {
      'X-Forwarded-For': ip,
    },
  });

  check(res, {
    'status is 200 or 429': (r) => r.status === 200 || r.status === 429,
  });
}