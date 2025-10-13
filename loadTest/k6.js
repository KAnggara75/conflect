import http from 'k6/http';
import { check, sleep } from 'k6';

// Konfigurasi load test
export const options = {
    vus: 100,        // jumlah virtual users
    duration: '30s', // durasi test
};

// Token auth bisa diganti sesuai kebutuhan
const AUTH_TOKEN = 'Bearer iniRahasiaWebhook';

// daftar environment yang mungkin
const ENVS = ['dev', 'prd', 'aws'];

// Endpoint dasar
const BASE_URL = 'http://localhost:8080/pakaiwa';

export default function () {
    // Pilih environment acak: dev atau prd
    const env = ENVS[Math.floor(Math.random() * ENVS.length)];
    const url = `${BASE_URL}/${env}`;

    // Set header authorization
    const headers = {
        headers: {
            Authorization: AUTH_TOKEN,
            'Content-Type': 'application/json',
        },
    };

    // Kirim request
    const res = http.get(url, headers);

    // Validasi hasil response
    check(res, {
        'status 200': (r) => r.status === 200,
        'body not empty': (r) => r.body && r.body.length > 0,
    });

    // Tidur 1 detik antar request
    sleep(1);
}
