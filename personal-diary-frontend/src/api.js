// import axios from 'axios';

// const API_URL = 'http://localhost:8080/api'; // Replace with your backend URL
// const api = axios.create({
//   baseURL: API_URL,
// });

// // Set up interceptors to add the token to every request
// api.interceptors.request.use(config => {
//     const token = localStorage.getItem('token');
//     if (token) {
//         config.headers['Authorization'] = `Bearer ${token}`;
//     }
//     return config;
//     }, error => {
//     return Promise.reject(error); 
//     }
// );

// export default api;


// src/api/api.js
import axios from "axios";
import { encryptAES, decryptAES } from "./utils/aes"; // Adjust the import path as necessary

// API URL using Vite environment variable
const API_URL = import.meta.env.VITE_BACKEND_URL;
const ENCRYPTION_AND_DECRYPTION = true; // Set to true if you want to encrypt the data

const api = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Add Bearer token (already implemented)
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }

    // Encrypt request data if exists
    if (config.data && ENCRYPTION_AND_DECRYPTION) {
      config.data = encryptAES(config.data); // Send encrypted string
    }else {
      config.data = JSON.stringify(config.data); // Send as JSON string
    }

    return config;
  },
  (error) => Promise.reject(error)
);

// Decrypt all responses
api.interceptors.response.use(
  (response) => {
    try {
      // Decrypt AES only if data is a string (encrypted)
      if (typeof response.data === "string" && ENCRYPTION_AND_DECRYPTION) {
        return decryptAES(response.data);
      } else if (typeof response.data === "object") {
        // If it's an object, return it as is
        return response.data;
      } 
      return response.data;
    } catch (err) {
      console.error("Decryption failed", err);
      return Promise.reject({ error: "Failed to decrypt server response." });
    }
  },
  (error) => Promise.reject(error)
);

export default api;
