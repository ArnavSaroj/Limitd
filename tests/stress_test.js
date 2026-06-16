import http from 'k6/http';

export const options = {
  stages: [
    { duration: '30s', target: 200 },
    { duration: '30s', target: 500 },
    { duration: '30s', target: 1000 },
    { duration: '30s', target: 2000 },
    { duration: '30s', target: 0 },
  ],
};

export default function () {
  const ip = `10.${__VU % 255}.${__ITER % 255}.${(__VU + __ITER) % 255}`;

  http.get('http://localhost:8080/', {
    headers: {
      'X-Forwarded-For': ip,
    },
  });
}


