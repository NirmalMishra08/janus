import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 },   // Ramp up to 50 users
    { duration: '1m',  target: 100 },  // Stay at 100 users
    { duration: '30s', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],   // 95% of requests under 500ms
    http_req_failed: ['rate<0.01'],     // Less than 1% failures
  },
};

const BASE_URL = 'http://localhost:8080';

// Test data
const payload = JSON.stringify({
  name: "Test User",
  email: "test@example.com"
});

const params = {
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer YOUR_JWT_TOKEN_HERE'  // Replace if needed
  },
};

export default function () {
  // 1. Health Check
  let res = http.get(`${BASE_URL}/health`);
  check(res, {
    'health check is 200': (r) => r.status === 200,
  });

  // 2. Main API endpoints
  res = http.get(`${BASE_URL}/users`);
  check(res, {
    'GET /users status is 200': (r) => r.status === 200,
  });

  // 3. POST request example
  res = http.post(`${BASE_URL}/users`, payload, params);
  check(res, {
    'POST /users status is 2xx': (r) => r.status >= 200 && r.status < 300,
  });

  sleep(1); // Think time
}