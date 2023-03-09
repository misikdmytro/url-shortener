import { check } from 'k6';
import http from 'k6/http';

export const options = {
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            rate: 100,
            timeUnit: '1s',
            duration: '2m',
            preAllocatedVUs: 10,
            maxVUs: 100,
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(99.9)<600'],
    },
    summaryTrendStats: ['avg', 'p(90)', 'p(95)', 'p(99)', 'p(99.9)'],
};

export default function () {
    const url = 'https://www.google.com';
    const body = {
        url,
        duration: 60,
    };

    const save = http.put(`${__ENV.BASE_URL}/shorten`, JSON.stringify(body), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(save, {
        'status is 201': (r) => r.status === 201,
        'response body': (r) => {
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

    if (save.status === 201) {
        const body = JSON.parse(save.body);

        const get = http.get(
            `${__ENV.BASE_URL}/${body.key}`,
            {
                headers: { 'Content-Type': 'application/json' },
                redirects: 0
            }
        );

        check(get, {
            'status is 301': (r) => r.status === 301,
            'response location': (r) => r.headers.Location === url,
        });
    }
}
