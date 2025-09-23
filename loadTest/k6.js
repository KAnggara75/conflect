import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 10,        // jumlah virtual users
  duration: '10s', // durasi test
};

export default function () {
  const res = http.get('http://localhost:8080/pakaiwa/dev');

  check(res, {
    'status 200': (r) => r.status === 200,
    'body not empty': (r) => r.body && r.body.length > 0,
  });

  sleep(1); // biar simulasi realistis, tiap user tunggu 1s antar request
}
