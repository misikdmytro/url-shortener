import http from 'k6/http';
import { check } from 'k6';

export const options = {
    stages: [
        { target: 60, duration: '1m' },
        { target: 60, duration: '3m' },
        { target: 0, duration: '30s' },
    ],
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(95)<500'],
    },
};

const baseuURL = 'https://qkbobxlzzi.execute-api.eu-central-1.amazonaws.com/develop';
export default function () {
    const body = {
        url: 'https://www.google.com',
        duration: 60,
    };

    const res = http.put(`${baseuURL}/shorten`, JSON.stringify(body), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(res, {
        'status is 201': (r) => r.status === 201,
        'response body': (r) => {
            const body = JSON.parse(r.body);
            return body.url.length > 0 && body.key.length > 0;
        },
    });
}