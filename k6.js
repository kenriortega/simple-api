import http from 'k6/http';
import { check } from 'k6';
export const options = {
    vus: 1000, // virtual users
    duration: '30s',
};

export default function () {

    const res = http.post('http://localhost:8080/all', fd.body(), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
        'is status 200': (r) => r.status === 200,
    });
}
