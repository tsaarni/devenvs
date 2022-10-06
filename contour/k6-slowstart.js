import http from 'k6/http';
import exec from 'k6/execution';

export default function () {
    const res = http.get('http://echoserver.127-0-0-101.nip.io');
   if (res.status === 200) {
       exec.vu.tags['pod'] = res.json()['pod'];
   }
};
