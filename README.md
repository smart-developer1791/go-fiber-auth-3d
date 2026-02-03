# 🔐 3D Glass Auth

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-v2.52-00ACD7?style=for-the-badge&logo=go&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-SQLite-336791?style=for-the-badge&logo=sqlite&logoColor=white)
![Tailwind](https://img.shields.io/badge/Tailwind-CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)
![CGO](https://img.shields.io/badge/CGO-Not_Required-success?style=for-the-badge)
![Render](https://img.shields.io/badge/Render-Deployed-46E3B7?style=for-the-badge&logo=render&logoColor=white)

**Stunning 3D glassmorphism authentication system with immersive visual effects.** Built with Go Fiber, GORM, pure Go SQLite, and cutting-edge CSS 3D transforms. Zero CGO dependencies — runs anywhere!

## ✨ Features

- 🎨 **Immersive 3D Effects** — Real-time mouse-tracking parallax on glass cards
- 🌀 **Animated Geometry** — Floating shapes, rotating cubes, spinning torus rings
- 💎 **Glassmorphism UI** — Frosted glass cards with depth and glow effects
- ✨ **Particle System** — Dynamic floating particles throughout the scene
- 🔐 **Secure Auth** — Bcrypt password hashing with session management
- 📱 **Responsive Design** — Works beautifully on all devices
- 🚀 **Zero Config** — SQLite database auto-created on first run
- ⚡ **No CGO** — Pure Go SQLite driver, cross-compile anywhere

## 🎭 Visual Effects

| Effect | Description |
|--------|-------------|
| 🎴 **3D Card Tilt** | Cards follow mouse movement with realistic perspective |
| 🔮 **Glassmorphism** | Frosted glass with backdrop blur and inner glow |
| 🎪 **Rotating Cubes** | CSS 3D transformed cubes with wireframe edges |
| 🌈 **Gradient Glow** | Animated gradient halos behind cards |
| ⭐ **Particles** | Floating luminescent orbs rising through scene |
| 🌐 **Perspective Grid** | 3D grid floor with infinite animation |

## 🚀 Quick Start

Clone the repository:

```bash
git clone https://github.com/smart-developer1791/go-fiber-auth-3d
cd go-fiber-auth-3d
```

Initialize dependencies and run:

```bash
go mod tidy
go run .
```

Open browser:

```text
http://localhost:3000
```

## 🔑 Demo Credentials

| Field | Value |
|-------|-------|
| 📧 Email | `demo@glassauth.io` |
| 🔐 Password | `demo2024` |
| 📱 Phone | `+1 (555) 987-6543` |

## 🛠️ Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.21+** | Backend runtime |
| **Fiber v2** | High-performance web framework |
| **GORM** | ORM with auto-migrations |
| **glebarez/sqlite** | Pure Go SQLite driver (no CGO!) |
| **Bcrypt** | Secure password hashing |
| **Tailwind CSS** | Utility-first styling |
| **CSS 3D** | Hardware-accelerated transforms |

## 📁 Project Structure

```text
go-fiber-auth-3d/
├── main.go          # Server, routes, handlers, templates
├── auth.db          # SQLite database (auto-created)
├── render.yaml      # Render.com deployment config
├── .gitignore       # Git ignore rules
└── README.md        # Documentation
```

## 🌐 API Routes

| Method | Route | Description |
|--------|-------|-------------|
| `GET` | `/` | Redirect to login |
| `GET` | `/login` | Login page with 3D effects |
| `POST` | `/login` | Authenticate user |
| `GET` | `/register` | Registration page |
| `POST` | `/register` | Create new account |
| `GET` | `/dashboard` | Protected dashboard |
| `POST` | `/logout` | End session |

## 🎨 Customization

### Change Color Theme

Modify gradient colors in the CSS:

```css
/* Login theme - Purple/Blue */
background: linear-gradient(135deg, #8b5cf6, #3b82f6);

/* Register theme - Green/Cyan */
background: linear-gradient(135deg, #10b981, #06b6d4);
```

### Adjust 3D Intensity

Change parallax sensitivity in JavaScript:

```javascript
// Lower = more sensitive, Higher = less sensitive
const xAxis = (window.innerWidth / 2 - e.pageX) / 25;
const yAxis = (window.innerHeight / 2 - e.pageY) / 25;
```

## 🔒 Security Features

- ✅ Bcrypt password hashing (cost factor 10)
- ✅ HTTP-only session cookies
- ✅ Protected route middleware
- ✅ Input validation
- ✅ SQL injection prevention via GORM

## 📊 Database Schema

```text
┌─────────────────────────────────────┐
│              users                  │
├─────────────────────────────────────┤
│ id         INTEGER PRIMARY KEY      │
│ email      TEXT UNIQUE NOT NULL     │
│ phone      TEXT                     │
│ password   TEXT NOT NULL            │
│ created_at DATETIME                 │
└─────────────────────────────────────┘
```

## 🌟 Why This Project?

Most auth UIs are boring forms. This project proves authentication can be:

- 🎭 **Visually stunning** without sacrificing UX
- ⚡ **Fast** using native CSS transforms (GPU accelerated)
- 🛡️ **Secure** with industry-standard practices
- 📦 **Minimal** — single file, zero external assets
- 🔧 **Portable** — no CGO, runs on any platform

## 🔧 Why Pure Go SQLite?

This project uses `github.com/glebarez/sqlite` instead of the traditional `gorm.io/driver/sqlite`:

| Feature | gorm.io/driver/sqlite | glebarez/sqlite |
|---------|----------------------|-----------------|
| CGO Required | ✅ Yes | ❌ No |
| Cross-compile | ❌ Complex | ✅ Easy |
| Windows build | ❌ Needs GCC | ✅ Just works |
| Performance | Faster | Slightly slower |
| Compatibility | Full | 99%+ |

---

## Deploy in 10 seconds

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)
