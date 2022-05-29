import http from 'k6/http';
import { check } from 'k6';

// Node Pool 1-13 Dikalkulasikan oleh KubeEP (Skenario HPA leluasa menambahkan pod)

export const options = {
    tags : {
        task: "scenario_3"
    },
    scenarios: {
        scenario_3: {
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
