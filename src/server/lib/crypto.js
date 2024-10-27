const { mkdir, rm, readFile } = require('node:fs/promises');
const { execSync } = require('node:child_process');

// create private key and self-signed X509 cert for https server
exports.createWebKeyAndCert = async function () {
    await mkdir('tmp-cert');
    execSync(`openssl req -newkey rsa:2048 -nodes -keyout tmp-cert/ssl.key -out tmp-cert/ssl.crt -x509 -days 365`);
    const key = await readFile('tmp-cert/ssl.key');
    const cert = await readFile('tmp-cert/ssl.crt');
    await rm('tmp-cert', { recursive: true });

    return {
      key: key.toString(),
      cert: cert.toString(),
      date: Date.now(),
    };
}

// create a pub / private key pair for JWT token signing and verification
exports.createJwtKey = function (passphrase = '') {
  const { generateKeyPairSync } = require('crypto');

  return new Promise((resolve, reject) => {
    generateKeyPairSync('rsa', {
      modulusLength: 4096,
      publicExponent: 0x10001,
      publicKeyEncoding: {
        type:'pkcs1',
        format: 'pem'
      },
      privateKeyEncoding: {
        type: 'pkcs8',
        format: 'pem',
        cipher: 'aes-256-cbc',
        passphrase
      }
    }, ( err, publicKey, privateKey ) => {
         if (err) {
            reject(err);
          } else {
            resolve({ 
              public: publicKey.toString(),
              private: privateKey.toString()
          });
        }
    });
  });
};

// Maybe add the 3-Way Handshake Syn - Ack - Syn, - Ack  - Fin, idk
exports.sign = function (message, key, passphrase = '') {
  const signer = createSign('sha256');
  signer.update(message);
  signer.end();
  return signer.sign({ key, passphrase }, 'base64');
}

exports.verify = function (message, signature, publicKey) { 
  const verifier = createVerify('sha256');
  verifier.update(message);
  verifier.end();
  return verifier.verify(publicKey, signature, 'base64');
}