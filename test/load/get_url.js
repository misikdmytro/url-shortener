import { check } from "k6";
import http from "k6/http";

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: "constant-arrival-rate",
      rate: 1200,
      timeUnit: "1s",
      duration: "2m",
      preAllocatedVUs: 40,
      maxVUs: 200,
    },
  },
  thresholds: {
    http_req_duration: ["p(99.9)<300"],
  },
  summaryTrendStats: ["avg", "p(90)", "p(95)", "p(99)", "p(99.9)"],
};

const symbols =
  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
function randomString(length) {
  let result = "_";
  for (let i = 0; i < length; i++) {
    result += symbols.charAt(Math.floor(Math.random() * symbols.length));
  }
  return result;
}

export default function () {
  const res = http.get(`${__ENV.BASE_URL}/${randomString(8)}`, null, {
    headers: { "Content-Type": "application/json" },
  });

  check(res, {
    "status is 404": (r) => r.status === 404,
    "response body": (r) => {
      const body = JSON.parse(r.body);
      return body && body.error === "URL not found";
    },
  });
}
