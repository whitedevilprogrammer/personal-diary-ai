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

personal-diary/
├─ server/ # Go backend
│ ├─ config/
│ ├─ controllers/
│ ├─ middleware/
│ ├─ models/
│ ├─ routers/
│ ├─ services/
│ └─ utils/
├─ personal-diary-frontend/ # Vue.js frontend
│ ├─ src/
│ ├─ public/
│ └─ vite.config.js
├─ docker-compose.yml
└─ README.md


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

## ✅ Tips to Maximize Professionalism:
| Element            | Suggestion                                                                 |
|--------------------|-----------------------------------------------------------------------------|
| 🖼 Banner           | Add a custom **project banner** at the top (Canva / Figma)                 |
| 📸 Screenshots     | Upload `login`, `create diary`, `refine with AI`, and `profile` pages      |
| 📝 README Styling   | Use Markdown bullets, emojis, headings, and spacing                        |
| 🔗 GitHub Profile  | Make sure your profile `README.md` is updated too                          |

---

If you'd like, I can generate a **custom banner image** or markdown **badges** for GitHub, Vue, Go, Docker, etc.

Let me know if you'd like to:
- Auto-generate a `CONTRIBUTING.md` file
- Add GitHub Action CI workflows
- Add badge previews (OpenAI, Docker, Vue, etc.)

You're on the path to building an amazing portfolio 👨‍💻✨!


📸 Screenshots (Coming Soon)
Add images of your Login, Diary Editor, AI Suggestions, etc.

🤝 Contributing
Want to contribute? Awesome!

1. Fork the repo

2. Create a new branch: git checkout -b feature-name

3. Make your changes and commit

4. Push to your branch: git push origin feature-name

5. Open a Pull Request


---

✅ **👨‍💻 Author
Ellalan Haridoss — Full Stack Developer (Golang | Vue.js)**
- [GitHub: whitedevilprogrammer](https://github.com/whitedevilprogrammer)
- [LinkedIn: ellalan-haridoss](https://www.linkedin.com/in/ellalan-haridoss)

Let me know if you'd like help generating:
- GitHub profile README
- Project banner
- GitHub badges for OpenAI, Vue, Docker, Go, etc.


