import http from 'k6/http';
import { check } from 'k6';

// Node Pool 1-3 Tidak Dikalkulasi oleh KubeEP (Skenario HPA tidak bisa mencapai ke minimum pada saat persiapan)

export const options = {
    tags : {
        task: "scenario_1"
    },
    scenarios: {
        scenario_1: {
            executor: "ramping-vus",
            startVUs: 20,
            stages: [
                { duration: '55m', target: 20 },
                { duration: '4m', target: 100 },
                { duration: '30s', target: 250 },
                { duration: '30s', target: 380 },
                { duration: '60m', target: 400 },
                { duration: '5m', target: 20 },
            ],
        }
    }
};

export default function () {
    const res = http.get('http://35.197.144.132/factorial?number=5000');
    check(res, { 'success': (r) => r.status == 200});
}
