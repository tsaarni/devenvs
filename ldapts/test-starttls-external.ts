import { Client } from './src';
import * as fs from 'fs';

async function main(): Promise<void> {
  const client = new Client({
    url: 'ldap://localhost:389',
  });

  await client.startTLS(
    {
      ca: fs.readFileSync('dist/test-data/certs/server-ca.pem'),
      cert: fs.readFileSync('dist/test-data/certs/user.pem'),
      key: fs.readFileSync('dist/test-data/certs/user-key.pem'),
      maxVersion: 'TLSv1.2',
    }
  );
  await client.bind('EXTERNAL');

  const res = await client.search('dc=example,dc=org');
  console.log(res);

  return client.unbind();
}

main().catch((err) => {
  console.error(err);
});
