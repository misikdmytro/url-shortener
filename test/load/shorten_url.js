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
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(99.9)<300"],
  },
  summaryTrendStats: ["avg", "p(90)", "p(95)", "p(99)", "p(99.9)"],
};

export default function () {
  const body = {
    url: "https://www.google.com",
    duration: 60,
  };

  const res = http.put(`${__ENV.BASE_URL}/shorten`, JSON.stringify(body), {
    headers: { "Content-Type": "application/json" },
  });

  check(res, {
    "status is 201": (r) => r.status === 201,
    "response body": (r) => {
      const body = JSON.parse(r.body);
      return (
        body &&
        body.url &&
        body.key &&
        body.url.length > 0 &&
        body.key.length > 0
      );
    },
  });
}
