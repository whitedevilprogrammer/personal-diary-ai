<!-- Banner -->
<!--<p align="center">
  <img src="https://your-image-link-here.com/banner.png" alt="Personal Diary App" width="800"/>
</p> -->

# 📝 Personal Diary App

A secure, full-featured diary web app built with **Golang**, **Vue.js**, **MongoDB**, and integrated with **ChatGPT** and **Gemini AI**. Write, reflect, and refine your thoughts with the help of AI in a clean, responsive UI.

---

## 🚀 Features

- ✍️ Add, edit, delete diary entries
- 🔐 JWT-based authentication
- 📱 Mobile API support for attendance modules
- 🤖 AI-powered text refinement using ChatGPT & Gemini
- 📊 MongoDB for persistent diary storage
- 💬 Real-time chat support (optional with WebSocket)
- 🐳 Docker-ready for easy deployment

---

## 🧱 Tech Stack

| Layer       | Technology              |
|------------|--------------------------|
| Backend     | Golang (Gorilla Mux, MVC)|
| Frontend    | Vue.js 3 (Vite)          |
| Database    | MongoDB                  |
| Auth        | JWT                      |
| AI API      | OpenAI ChatGPT + Gemini  |
| Container   | Docker & Docker Compose  |

---

## 📁 Project Structure
```
personal-diary/
├── server/ # Go backend
│   ├── config/
│   ├── controllers/
│   ├── middleware/
│   ├── models/
│   ├── routers/
│   ├── services/
│   └── utils/
├── personal-diary-frontend/ # Vue.js frontend
│   ├── src/
│   └── public/
├── .vite.config.js
├── docker-compose.yml
└── README.md
```


---

## ⚙️ Getting Started

### Prerequisites

- Go >= 1.18
- Node.js >= 16.x
- MongoDB instance running
- OpenAI API key (for ChatGPT)
- Gemini API key (for Gemini)

### 🧪 Local Development

```bash
# Clone the repo
git clone https://github.com/whitedevilprogrammer/personal-diary.git
cd personal-diary

# Backend
cd server
go mod tidy
go run main.go

# Frontend
cd ../personal-diary-frontend
npm install
npm run dev
```

---

## 🛠️ Environment Setup

Create a `.env` file in the `server/` directory with the following variables:

```env
PORT=8080
MONGO_URI=mongodb://localhost:27017/personal-diary
JWT_SECRET=your_jwt_secret_here
OPENAI_API_KEY=your_openai_key_here
GEMINI_API_KEY=your_gemini_api_key_here
```

---

## 📸 Screenshots

A quick preview of the Personal Diary app powered by Golang, Vue.js, MongoDB, ChatGPT, and Gemini AI.

---

### 🏠 Dashboard — Diary Overview  
An elegant dashboard displaying your AI-enhanced diary entries.  
<img src="https://github.com/user-attachments/assets/7ec9a3a7-312a-49ec-80ce-4dbae91ac2c0" width="100%" alt="Diary Dashboard" />

---

### 🔐 Secure Login with Google OAuth (Golang)  
Authenticate securely with Google using Go’s OAuth integration.  
<img src="https://github.com/user-attachments/assets/04d41b95-d6d1-4894-a1b7-4147f84ceefe" width="100%" alt="Login Page" />

---

### 📝 Create a New Diary Entry  
Simple, responsive UI for adding personal thoughts or reflections.  
<img src="https://github.com/user-attachments/assets/c558c7ba-e019-43d6-a64a-2bdbc500eb8e" width="100%" alt="Create Diary Entry" />

---

### 🤖 Refine Thoughts with ChatGPT & Gemini AI  
Use AI to enhance or rewrite selected text inside entries.  
<img src="https://github.com/user-attachments/assets/4304e497-e93b-4e89-9875-7bdccf392075" width="100%" alt="AI Text Refinement" />

---

### 🌄 Auto-Generate Entry Backgrounds (Gemini AI)  
Dynamically generate image backgrounds based on mood or text.  
<img src="https://github.com/user-attachments/assets/58ae90fb-20ae-4d38-984c-64c5e4353db1" width="100%" alt="Gemini AI Background" />

---

### ✏️ Edit Diary Entries  
Easily update your thoughts with AI assistance and Markdown support.  
<img src="https://github.com/user-attachments/assets/dfae6f1e-9e77-4e79-aa77-e3998fdd1110" width="100%" alt="Edit Diary Entry" />

---

## 🔐 AES Encryption – Secure Payload Diary Entry

This diary entry documents my experience with implementing **secure payload encryption using AES (Advanced Encryption Standard)** between the Vue.js frontend and Golang backend.

It demonstrates:
- 🔒 How sensitive data can be encrypted on the frontend
- 🧩 How it is securely decrypted on the backend using the same AES key
- ✅ Real-world use of cryptography in full stack applications

<img src="https://github.com/user-attachments/assets/b69799b6-9cd7-4d7e-89f3-c460e0abaed9" width="100%" alt="AES Encryption Diary Entry" />

> This implementation helps ensure data privacy between client and server in production-ready apps like this personal diary.

---



✅ **👨‍💻 Author
Ellalan Haridoss — Full Stack Developer (Golang | Vue.js)**
- [GitHub: whitedevilprogrammer](https://github.com/whitedevilprogrammer)
- [LinkedIn: ellalan-haridoss](https://www.linkedin.com/in/ellalan-haridoss)

Let me know if you'd like help generating:
- GitHub profile README
- Project banner
- GitHub badges for OpenAI, Vue, Docker, Go, etc.


