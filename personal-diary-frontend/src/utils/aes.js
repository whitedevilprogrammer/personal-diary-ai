// src/utils/aes.js
import CryptoJS from "crypto-js";

const key = CryptoJS.enc.Utf8.parse("thisis32bitlongpassphraseimusing"); // 32 bytes
const iv = CryptoJS.enc.Utf8.parse("thisis16bytesiv!"); // 16 bytes

export function encryptAES(data) {
  const encrypted = CryptoJS.AES.encrypt(JSON.stringify(data), key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
  });
  console.log("Encrypted data:", encrypted.toString().trim());
  return encrypted.toString(); // base64
}

export function decryptAES(cipherText) {
  const decrypted = CryptoJS.AES.decrypt(cipherText, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
  });

  const decryptedText = decrypted.toString(CryptoJS.enc.Utf8);
  return JSON.parse(decryptedText);
}
