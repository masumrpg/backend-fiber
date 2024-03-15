import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  let startTime = new Date().getTime(); // Waktu awal permintaan

  let response = http.get('http://127.0.0.1:3000/api/users');

  let endTime = new Date().getTime(); // Waktu setelah menerima respons

  // Menghitung kecepatan fetch (fetch rate)
  let fetchTime = endTime - startTime;
  let fetchRate = (1 / fetchTime) * 1000; // Hitung fetch rate dalam requests per detik

  console.log(`Kecepatan fetch: ${fetchRate.toFixed(2)} requests/s`);

  check(response, {
    'Status is 200': (r) => r.status === 200
  });

  sleep(1);
}
