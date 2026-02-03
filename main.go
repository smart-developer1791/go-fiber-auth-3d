package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;size:255;not null"`
	Phone     string    `gorm:"size:20"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

var (
	db    *gorm.DB
	store *session.Store
)

func main() {
	initDatabase()
	seedDemoUser()

	store = session.New(session.Config{
		Expiration:     24 * time.Hour,
		CookieSecure:   false,
		CookieHTTPOnly: true,
	})

	app := fiber.New(fiber.Config{
		AppName: "3D Glass Auth",
	})

	app.Use(logger.New())

	app.Get("/", handleIndex)
	app.Get("/login", handleLoginPage)
	app.Post("/login", handleLogin)
	app.Get("/register", handleRegisterPage)
	app.Post("/register", handleRegister)
	app.Get("/dashboard", authRequired, handleDashboard)
	app.Post("/logout", handleLogout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 3D Glass Auth running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	db.AutoMigrate(&User{})
	log.Println("✅ Database initialized")
}

func seedDemoUser() {
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("demo2024"), bcrypt.DefaultCost)
		db.Create(&User{
			Email:    "demo@glassauth.io",
			Phone:    "+1 (555) 987-6543",
			Password: string(hash),
		})
		log.Println("✅ Demo user created: demo@glassauth.io / demo2024")
	}
}

func authRequired(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil || sess.Get("userID") == nil {
		return c.Redirect("/login")
	}
	return c.Next()
}

func handleIndex(c *fiber.Ctx) error {
	return c.Redirect("/login")
}

func handleLoginPage(c *fiber.Ctx) error {
	c.Type("html")
	return c.SendString(renderLoginPage(""))
}

func handleLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.Type("html")
		return c.SendString(renderLoginPage("Invalid credentials"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.Type("html")
		return c.SendString(renderLoginPage("Invalid credentials"))
	}

	sess, _ := store.Get(c)
	sess.Set("userID", user.ID)
	sess.Set("userEmail", user.Email)
	sess.Save()

	return c.Redirect("/dashboard")
}

func handleRegisterPage(c *fiber.Ctx) error {
	c.Type("html")
	return c.SendString(renderRegisterPage(""))
}

func handleRegister(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	if password != confirmPassword {
		c.Type("html")
		return c.SendString(renderRegisterPage("Passwords do not match"))
	}

	if len(password) < 6 {
		c.Type("html")
		return c.SendString(renderRegisterPage("Password must be at least 6 characters"))
	}

	var existing User
	if db.Where("email = ?", email).First(&existing).Error == nil {
		c.Type("html")
		return c.SendString(renderRegisterPage("Email already registered"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Type("html")
		return c.SendString(renderRegisterPage("Registration failed"))
	}

	user := User{
		Email:    email,
		Password: string(hash),
	}

	if err := db.Create(&user).Error; err != nil {
		c.Type("html")
		return c.SendString(renderRegisterPage("Registration failed"))
	}

	sess, _ := store.Get(c)
	sess.Set("userID", user.ID)
	sess.Set("userEmail", user.Email)
	sess.Save()

	return c.Redirect("/dashboard")
}

func handleDashboard(c *fiber.Ctx) error {
	sess, _ := store.Get(c)
	email := sess.Get("userEmail").(string)
	c.Type("html")
	return c.SendString(renderDashboard(email))
}

func handleLogout(c *fiber.Ctx) error {
	sess, _ := store.Get(c)
	sess.Destroy()
	return c.Redirect("/login")
}

func renderLoginPage(errorMsg string) string {
	errorHTML := ""
	if errorMsg != "" {
		errorHTML = fmt.Sprintf(`<div class="error-shake bg-red-500/20 border border-red-500/50 text-red-200 px-4 py-3 rounded-xl mb-6 backdrop-blur-sm">%s</div>`, errorMsg)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login | 3D Glass Auth</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        body {
            min-height: 100vh;
            background: linear-gradient(135deg, #0c0015 0%%, #1a0a2e 25%%, #16213e 50%%, #0f3460 75%%, #1a1a2e 100%%);
            overflow: hidden;
            font-family: 'Segoe UI', system-ui, sans-serif;
        }

        .scene {
            perspective: 1500px;
            perspective-origin: 50%% 50%%;
            position: fixed;
            inset: 0;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .floating-shapes {
            position: fixed;
            inset: 0;
            pointer-events: none;
            overflow: hidden;
        }

        .shape {
            position: absolute;
            border-radius: 50%%;
            background: linear-gradient(135deg, rgba(139, 92, 246, 0.3), rgba(59, 130, 246, 0.3));
            filter: blur(1px);
            animation: float 20s infinite ease-in-out;
        }

        .shape:nth-child(1) { width: 300px; height: 300px; top: -150px; left: 10%%; animation-delay: 0s; }
        .shape:nth-child(2) { width: 200px; height: 200px; top: 60%%; right: -100px; animation-delay: -5s; background: linear-gradient(135deg, rgba(236, 72, 153, 0.3), rgba(239, 68, 68, 0.3)); }
        .shape:nth-child(3) { width: 150px; height: 150px; bottom: -75px; left: 30%%; animation-delay: -10s; background: linear-gradient(135deg, rgba(34, 211, 238, 0.3), rgba(16, 185, 129, 0.3)); }
        .shape:nth-child(4) { width: 250px; height: 250px; top: 20%%; right: 20%%; animation-delay: -15s; }
        .shape:nth-child(5) { width: 180px; height: 180px; bottom: 20%%; left: -90px; animation-delay: -7s; background: linear-gradient(135deg, rgba(251, 191, 36, 0.3), rgba(245, 158, 11, 0.3)); }

        @keyframes float {
            0%%, 100%% { transform: translate(0, 0) rotate(0deg) scale(1); }
            25%% { transform: translate(30px, -30px) rotate(90deg) scale(1.1); }
            50%% { transform: translate(-20px, 20px) rotate(180deg) scale(0.9); }
            75%% { transform: translate(40px, 10px) rotate(270deg) scale(1.05); }
        }

        .geometric-grid {
            position: fixed;
            inset: 0;
            background-image: 
                linear-gradient(rgba(139, 92, 246, 0.03) 1px, transparent 1px),
                linear-gradient(90deg, rgba(139, 92, 246, 0.03) 1px, transparent 1px);
            background-size: 50px 50px;
            transform: perspective(500px) rotateX(60deg);
            transform-origin: center top;
            animation: gridMove 20s linear infinite;
        }

        @keyframes gridMove {
            0%% { background-position: 0 0; }
            100%% { background-position: 50px 50px; }
        }

        .cube-container {
            position: fixed;
            width: 100px;
            height: 100px;
            transform-style: preserve-3d;
            animation: rotateCube 25s infinite linear;
        }

        .cube-container.left { left: 10%%; top: 30%%; }
        .cube-container.right { right: 10%%; bottom: 30%%; animation-direction: reverse; }

        .cube-face {
            position: absolute;
            width: 100px;
            height: 100px;
            border: 2px solid rgba(139, 92, 246, 0.3);
            background: rgba(139, 92, 246, 0.05);
            backdrop-filter: blur(5px);
        }

        .cube-face:nth-child(1) { transform: rotateY(0deg) translateZ(50px); }
        .cube-face:nth-child(2) { transform: rotateY(180deg) translateZ(50px); }
        .cube-face:nth-child(3) { transform: rotateY(90deg) translateZ(50px); }
        .cube-face:nth-child(4) { transform: rotateY(-90deg) translateZ(50px); }
        .cube-face:nth-child(5) { transform: rotateX(90deg) translateZ(50px); }
        .cube-face:nth-child(6) { transform: rotateX(-90deg) translateZ(50px); }

        @keyframes rotateCube {
            0%% { transform: rotateX(0deg) rotateY(0deg); }
            100%% { transform: rotateX(360deg) rotateY(360deg); }
        }

        .glass-card {
            width: 420px;
            padding: 3rem;
            background: rgba(255, 255, 255, 0.03);
            backdrop-filter: blur(20px);
            border-radius: 24px;
            border: 1px solid rgba(255, 255, 255, 0.1);
            box-shadow: 
                0 25px 50px -12px rgba(0, 0, 0, 0.5),
                0 0 0 1px rgba(255, 255, 255, 0.05) inset,
                0 -20px 40px -20px rgba(139, 92, 246, 0.3) inset;
            transform-style: preserve-3d;
            transform: rotateX(5deg) rotateY(0deg);
            transition: transform 0.1s ease-out;
            animation: cardEntrance 1s ease-out;
        }

        @keyframes cardEntrance {
            0%% { opacity: 0; transform: rotateX(20deg) rotateY(-20deg) translateZ(-100px); }
            100%% { opacity: 1; transform: rotateX(5deg) rotateY(0deg) translateZ(0); }
        }

        .card-glow {
            position: absolute;
            inset: -2px;
            background: linear-gradient(135deg, rgba(139, 92, 246, 0.5), rgba(59, 130, 246, 0.5), rgba(236, 72, 153, 0.5));
            border-radius: 26px;
            z-index: -1;
            filter: blur(20px);
            opacity: 0.5;
            animation: glowPulse 3s ease-in-out infinite;
        }

        @keyframes glowPulse {
            0%%, 100%% { opacity: 0.3; transform: scale(1); }
            50%% { opacity: 0.6; transform: scale(1.02); }
        }

        .form-title {
            font-size: 2rem;
            font-weight: 700;
            text-align: center;
            margin-bottom: 0.5rem;
            background: linear-gradient(135deg, #fff 0%%, #a78bfa 50%%, #60a5fa 100%%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            text-shadow: 0 0 40px rgba(139, 92, 246, 0.5);
        }

        .form-subtitle {
            text-align: center;
            color: rgba(255, 255, 255, 0.5);
            margin-bottom: 2rem;
            font-size: 0.9rem;
        }

        .input-group {
            position: relative;
            margin-bottom: 1.5rem;
        }

        .input-group label {
            display: block;
            color: rgba(255, 255, 255, 0.7);
            font-size: 0.85rem;
            margin-bottom: 0.5rem;
            font-weight: 500;
        }

        .input-group input {
            width: 100%%;
            padding: 1rem 1.25rem;
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            color: white;
            font-size: 1rem;
            transition: all 0.3s ease;
            outline: none;
        }

        .input-group input:focus {
            border-color: rgba(139, 92, 246, 0.5);
            background: rgba(255, 255, 255, 0.08);
            box-shadow: 0 0 20px rgba(139, 92, 246, 0.2);
        }

        .input-group input::placeholder {
            color: rgba(255, 255, 255, 0.3);
        }

        .submit-btn {
            width: 100%%;
            padding: 1rem;
            background: linear-gradient(135deg, #8b5cf6 0%%, #6366f1 50%%, #3b82f6 100%%);
            border: none;
            border-radius: 12px;
            color: white;
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
            margin-top: 1rem;
        }

        .submit-btn::before {
            content: '';
            position: absolute;
            inset: 0;
            background: linear-gradient(135deg, transparent, rgba(255, 255, 255, 0.2), transparent);
            transform: translateX(-100%%);
            transition: transform 0.5s ease;
        }

        .submit-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 30px rgba(139, 92, 246, 0.4);
        }

        .submit-btn:hover::before {
            transform: translateX(100%%);
        }

        .submit-btn:active {
            transform: translateY(0);
        }

        .alt-action {
            text-align: center;
            margin-top: 1.5rem;
            color: rgba(255, 255, 255, 0.5);
            font-size: 0.9rem;
        }

        .alt-action a {
            color: #a78bfa;
            text-decoration: none;
            font-weight: 500;
            transition: color 0.3s ease;
        }

        .alt-action a:hover {
            color: #c4b5fd;
            text-decoration: underline;
        }

        .demo-hint {
            margin-top: 1.5rem;
            padding: 1rem;
            background: rgba(139, 92, 246, 0.1);
            border-radius: 12px;
            border: 1px solid rgba(139, 92, 246, 0.2);
        }

        .demo-hint p {
            color: rgba(255, 255, 255, 0.6);
            font-size: 0.8rem;
            margin: 0;
        }

        .demo-hint code {
            color: #a78bfa;
            background: rgba(139, 92, 246, 0.2);
            padding: 0.1rem 0.4rem;
            border-radius: 4px;
            font-size: 0.75rem;
        }

        .particles {
            position: fixed;
            inset: 0;
            pointer-events: none;
        }

        .particle {
            position: absolute;
            width: 4px;
            height: 4px;
            background: rgba(139, 92, 246, 0.6);
            border-radius: 50%%;
            animation: particleFloat 15s infinite linear;
        }

        @keyframes particleFloat {
            0%% { transform: translateY(100vh) rotate(0deg); opacity: 0; }
            10%% { opacity: 1; }
            90%% { opacity: 1; }
            100%% { transform: translateY(-100vh) rotate(720deg); opacity: 0; }
        }

        .error-shake {
            animation: shake 0.5s ease-in-out;
        }

        @keyframes shake {
            0%%, 100%% { transform: translateX(0); }
            20%% { transform: translateX(-10px); }
            40%% { transform: translateX(10px); }
            60%% { transform: translateX(-10px); }
            80%% { transform: translateX(10px); }
        }

        .torus {
            position: fixed;
            width: 200px;
            height: 200px;
            border: 30px solid transparent;
            border-radius: 50%%;
            border-top-color: rgba(139, 92, 246, 0.2);
            border-bottom-color: rgba(59, 130, 246, 0.2);
            animation: spinTorus 10s linear infinite;
        }

        .torus.one { top: 5%%; left: 5%%; }
        .torus.two { bottom: 5%%; right: 5%%; animation-direction: reverse; border-top-color: rgba(236, 72, 153, 0.2); }

        @keyframes spinTorus {
            0%% { transform: rotateX(45deg) rotateZ(0deg); }
            100%% { transform: rotateX(45deg) rotateZ(360deg); }
        }
    </style>
</head>
<body>
    <div class="floating-shapes">
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
    </div>

    <div class="geometric-grid"></div>

    <div class="cube-container left">
        <div class="cube-face"></div>
        <div class="cube-face"></div>
        <div class="cube-face"></div>
        <div class="cube-face"></div>
        <div class="cube-face"></div>
        <div class="cube-face"></div>
    </div>

    <div class="cube-container right">
        <div class="cube-face"></div>
        <div class="cube-face"></div>
        <div class="cube-face"></div>
        <div class="cube-face"></div>
        <div class="cube-face"></div>
        <div class="cube-face"></div>
    </div>

    <div class="torus one"></div>
    <div class="torus two"></div>

    <div class="particles" id="particles"></div>

    <div class="scene">
        <div class="glass-card" id="card">
            <div class="card-glow"></div>
            <h1 class="form-title">Welcome Back</h1>
            <p class="form-subtitle">Enter your credentials to continue</p>

            %s

            <form method="POST" action="/login">
                <div class="input-group">
                    <label>Email Address</label>
                    <input type="email" name="email" placeholder="you@example.com" required>
                </div>

                <div class="input-group">
                    <label>Password</label>
                    <input type="password" name="password" placeholder="••••••••" required>
                </div>

                <button type="submit" class="submit-btn">Sign In</button>
            </form>

            <p class="alt-action">Don't have an account? <a href="/register">Create one</a></p>

            <div class="demo-hint">
                <p>🔐 Demo: <code>demo@glassauth.io</code> / <code>demo2024</code></p>
                <p style="margin-top: 0.5rem">📱 Phone: <code>+1 (555) 987-6543</code></p>
            </div>
        </div>
    </div>

    <script>
        const particlesContainer = document.getElementById('particles');
        for (let i = 0; i < 30; i++) {
            const particle = document.createElement('div');
            particle.className = 'particle';
            particle.style.left = Math.random() * 100 + '%%';
            particle.style.animationDelay = Math.random() * 15 + 's';
            particle.style.animationDuration = (10 + Math.random() * 10) + 's';
            particlesContainer.appendChild(particle);
        }

        const card = document.getElementById('card');
        document.addEventListener('mousemove', (e) => {
            const xAxis = (window.innerWidth / 2 - e.pageX) / 25;
            const yAxis = (window.innerHeight / 2 - e.pageY) / 25;
            card.style.transform = 'rotateY(' + xAxis + 'deg) rotateX(' + yAxis + 'deg)';
        });

        document.addEventListener('mouseleave', () => {
            card.style.transform = 'rotateX(5deg) rotateY(0deg)';
        });
    </script>
</body>
</html>`, errorHTML)
}

func renderRegisterPage(errorMsg string) string {
	errorHTML := ""
	if errorMsg != "" {
		errorHTML = fmt.Sprintf(`<div class="error-shake bg-red-500/20 border border-red-500/50 text-red-200 px-4 py-3 rounded-xl mb-6 backdrop-blur-sm">%s</div>`, errorMsg)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Register | 3D Glass Auth</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        body {
            min-height: 100vh;
            background: linear-gradient(135deg, #0c0015 0%%, #1a0a2e 25%%, #16213e 50%%, #0f3460 75%%, #1a1a2e 100%%);
            overflow: hidden;
            font-family: 'Segoe UI', system-ui, sans-serif;
        }

        .scene {
            perspective: 1500px;
            perspective-origin: 50%% 50%%;
            position: fixed;
            inset: 0;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .floating-shapes {
            position: fixed;
            inset: 0;
            pointer-events: none;
            overflow: hidden;
        }

        .shape {
            position: absolute;
            border-radius: 50%%;
            background: linear-gradient(135deg, rgba(16, 185, 129, 0.3), rgba(34, 211, 238, 0.3));
            filter: blur(1px);
            animation: float 20s infinite ease-in-out;
        }

        .shape:nth-child(1) { width: 300px; height: 300px; top: -150px; left: 10%%; animation-delay: 0s; background: linear-gradient(135deg, rgba(34, 211, 238, 0.3), rgba(59, 130, 246, 0.3)); }
        .shape:nth-child(2) { width: 200px; height: 200px; top: 60%%; right: -100px; animation-delay: -5s; background: linear-gradient(135deg, rgba(16, 185, 129, 0.3), rgba(52, 211, 153, 0.3)); }
        .shape:nth-child(3) { width: 150px; height: 150px; bottom: -75px; left: 30%%; animation-delay: -10s; }
        .shape:nth-child(4) { width: 250px; height: 250px; top: 20%%; right: 20%%; animation-delay: -15s; background: linear-gradient(135deg, rgba(139, 92, 246, 0.3), rgba(168, 85, 247, 0.3)); }
        .shape:nth-child(5) { width: 180px; height: 180px; bottom: 20%%; left: -90px; animation-delay: -7s; background: linear-gradient(135deg, rgba(236, 72, 153, 0.3), rgba(244, 114, 182, 0.3)); }

        @keyframes float {
            0%%, 100%% { transform: translate(0, 0) rotate(0deg) scale(1); }
            25%% { transform: translate(30px, -30px) rotate(90deg) scale(1.1); }
            50%% { transform: translate(-20px, 20px) rotate(180deg) scale(0.9); }
            75%% { transform: translate(40px, 10px) rotate(270deg) scale(1.05); }
        }

        .geometric-grid {
            position: fixed;
            inset: 0;
            background-image: 
                linear-gradient(rgba(16, 185, 129, 0.03) 1px, transparent 1px),
                linear-gradient(90deg, rgba(16, 185, 129, 0.03) 1px, transparent 1px);
            background-size: 50px 50px;
            transform: perspective(500px) rotateX(60deg);
            transform-origin: center top;
            animation: gridMove 20s linear infinite;
        }

        @keyframes gridMove {
            0%% { background-position: 0 0; }
            100%% { background-position: 50px 50px; }
        }

        .pyramid {
            position: fixed;
            width: 0;
            height: 0;
            border-left: 60px solid transparent;
            border-right: 60px solid transparent;
            border-bottom: 100px solid rgba(16, 185, 129, 0.15);
            animation: rotatePyramid 15s linear infinite;
        }

        .pyramid.one { top: 15%%; left: 8%%; }
        .pyramid.two { bottom: 15%%; right: 8%%; animation-direction: reverse; border-bottom-color: rgba(34, 211, 238, 0.15); }

        @keyframes rotatePyramid {
            0%% { transform: rotateY(0deg); }
            100%% { transform: rotateY(360deg); }
        }

        .glass-card {
            width: 420px;
            padding: 3rem;
            background: rgba(255, 255, 255, 0.03);
            backdrop-filter: blur(20px);
            border-radius: 24px;
            border: 1px solid rgba(255, 255, 255, 0.1);
            box-shadow: 
                0 25px 50px -12px rgba(0, 0, 0, 0.5),
                0 0 0 1px rgba(255, 255, 255, 0.05) inset,
                0 -20px 40px -20px rgba(16, 185, 129, 0.3) inset;
            transform-style: preserve-3d;
            transform: rotateX(5deg) rotateY(0deg);
            transition: transform 0.1s ease-out;
            animation: cardEntrance 1s ease-out;
        }

        @keyframes cardEntrance {
            0%% { opacity: 0; transform: rotateX(-20deg) rotateY(20deg) translateZ(-100px); }
            100%% { opacity: 1; transform: rotateX(5deg) rotateY(0deg) translateZ(0); }
        }

        .card-glow {
            position: absolute;
            inset: -2px;
            background: linear-gradient(135deg, rgba(16, 185, 129, 0.5), rgba(34, 211, 238, 0.5), rgba(59, 130, 246, 0.5));
            border-radius: 26px;
            z-index: -1;
            filter: blur(20px);
            opacity: 0.5;
            animation: glowPulse 3s ease-in-out infinite;
        }

        @keyframes glowPulse {
            0%%, 100%% { opacity: 0.3; transform: scale(1); }
            50%% { opacity: 0.6; transform: scale(1.02); }
        }

        .form-title {
            font-size: 2rem;
            font-weight: 700;
            text-align: center;
            margin-bottom: 0.5rem;
            background: linear-gradient(135deg, #fff 0%%, #34d399 50%%, #22d3ee 100%%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            text-shadow: 0 0 40px rgba(16, 185, 129, 0.5);
        }

        .form-subtitle {
            text-align: center;
            color: rgba(255, 255, 255, 0.5);
            margin-bottom: 2rem;
            font-size: 0.9rem;
        }

        .input-group {
            position: relative;
            margin-bottom: 1.25rem;
        }

        .input-group label {
            display: block;
            color: rgba(255, 255, 255, 0.7);
            font-size: 0.85rem;
            margin-bottom: 0.5rem;
            font-weight: 500;
        }

        .input-group input {
            width: 100%%;
            padding: 1rem 1.25rem;
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            color: white;
            font-size: 1rem;
            transition: all 0.3s ease;
            outline: none;
        }

        .input-group input:focus {
            border-color: rgba(16, 185, 129, 0.5);
            background: rgba(255, 255, 255, 0.08);
            box-shadow: 0 0 20px rgba(16, 185, 129, 0.2);
        }

        .input-group input::placeholder {
            color: rgba(255, 255, 255, 0.3);
        }

        .submit-btn {
            width: 100%%;
            padding: 1rem;
            background: linear-gradient(135deg, #10b981 0%%, #14b8a6 50%%, #06b6d4 100%%);
            border: none;
            border-radius: 12px;
            color: white;
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
            margin-top: 0.5rem;
        }

        .submit-btn::before {
            content: '';
            position: absolute;
            inset: 0;
            background: linear-gradient(135deg, transparent, rgba(255, 255, 255, 0.2), transparent);
            transform: translateX(-100%%);
            transition: transform 0.5s ease;
        }

        .submit-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 30px rgba(16, 185, 129, 0.4);
        }

        .submit-btn:hover::before {
            transform: translateX(100%%);
        }

        .submit-btn:active {
            transform: translateY(0);
        }

        .alt-action {
            text-align: center;
            margin-top: 1.5rem;
            color: rgba(255, 255, 255, 0.5);
            font-size: 0.9rem;
        }

        .alt-action a {
            color: #34d399;
            text-decoration: none;
            font-weight: 500;
            transition: color 0.3s ease;
        }

        .alt-action a:hover {
            color: #6ee7b7;
            text-decoration: underline;
        }

        .particles {
            position: fixed;
            inset: 0;
            pointer-events: none;
        }

        .particle {
            position: absolute;
            width: 4px;
            height: 4px;
            background: rgba(16, 185, 129, 0.6);
            border-radius: 50%%;
            animation: particleFloat 15s infinite linear;
        }

        @keyframes particleFloat {
            0%% { transform: translateY(100vh) rotate(0deg); opacity: 0; }
            10%% { opacity: 1; }
            90%% { opacity: 1; }
            100%% { transform: translateY(-100vh) rotate(720deg); opacity: 0; }
        }

        .error-shake {
            animation: shake 0.5s ease-in-out;
        }

        @keyframes shake {
            0%%, 100%% { transform: translateX(0); }
            20%% { transform: translateX(-10px); }
            40%% { transform: translateX(10px); }
            60%% { transform: translateX(-10px); }
            80%% { transform: translateX(10px); }
        }

        .hex-ring {
            position: fixed;
            width: 150px;
            height: 150px;
            border: 3px solid rgba(16, 185, 129, 0.2);
            clip-path: polygon(50%% 0%%, 100%% 25%%, 100%% 75%%, 50%% 100%%, 0%% 75%%, 0%% 25%%);
            animation: spinHex 20s linear infinite;
        }

        .hex-ring.one { top: 10%%; right: 15%%; }
        .hex-ring.two { bottom: 10%%; left: 15%%; animation-direction: reverse; border-color: rgba(34, 211, 238, 0.2); }

        @keyframes spinHex {
            0%% { transform: rotate(0deg); }
            100%% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="floating-shapes">
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
    </div>

    <div class="geometric-grid"></div>

    <div class="pyramid one"></div>
    <div class="pyramid two"></div>

    <div class="hex-ring one"></div>
    <div class="hex-ring two"></div>

    <div class="particles" id="particles"></div>

    <div class="scene">
        <div class="glass-card" id="card">
            <div class="card-glow"></div>
            <h1 class="form-title">Create Account</h1>
            <p class="form-subtitle">Join us and start your journey</p>

            %s

            <form method="POST" action="/register">
                <div class="input-group">
                    <label>Email Address</label>
                    <input type="email" name="email" placeholder="you@example.com" required>
                </div>

                <div class="input-group">
                    <label>Password</label>
                    <input type="password" name="password" placeholder="••••••••" required minlength="6">
                </div>

                <div class="input-group">
                    <label>Confirm Password</label>
                    <input type="password" name="confirm_password" placeholder="••••••••" required minlength="6">
                </div>

                <button type="submit" class="submit-btn">Create Account</button>
            </form>

            <p class="alt-action">Already have an account? <a href="/login">Sign in</a></p>
        </div>
    </div>

    <script>
        const particlesContainer = document.getElementById('particles');
        for (let i = 0; i < 30; i++) {
            const particle = document.createElement('div');
            particle.className = 'particle';
            particle.style.left = Math.random() * 100 + '%%';
            particle.style.animationDelay = Math.random() * 15 + 's';
            particle.style.animationDuration = (10 + Math.random() * 10) + 's';
            particlesContainer.appendChild(particle);
        }

        const card = document.getElementById('card');
        document.addEventListener('mousemove', (e) => {
            const xAxis = (window.innerWidth / 2 - e.pageX) / 25;
            const yAxis = (window.innerHeight / 2 - e.pageY) / 25;
            card.style.transform = 'rotateY(' + xAxis + 'deg) rotateX(' + yAxis + 'deg)';
        });

        document.addEventListener('mouseleave', () => {
            card.style.transform = 'rotateX(5deg) rotateY(0deg)';
        });
    </script>
</body>
</html>`, errorHTML)
}

func renderDashboard(email string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard | 3D Glass Auth</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        body {
            min-height: 100vh;
            background: linear-gradient(135deg, #0c0015 0%%, #1a0a2e 25%%, #16213e 50%%, #0f3460 75%%, #1a1a2e 100%%);
            font-family: 'Segoe UI', system-ui, sans-serif;
        }

        .navbar {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            height: 70px;
            background: rgba(255, 255, 255, 0.03);
            backdrop-filter: blur(20px);
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0 2rem;
            z-index: 100;
        }

        .logo {
            font-size: 1.5rem;
            font-weight: 700;
            background: linear-gradient(135deg, #fff 0%%, #a78bfa 50%%, #60a5fa 100%%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }

        .user-section {
            display: flex;
            align-items: center;
            gap: 1.5rem;
        }

        .user-email {
            color: rgba(255, 255, 255, 0.7);
            font-size: 0.9rem;
        }

        .logout-btn {
            padding: 0.6rem 1.5rem;
            background: rgba(239, 68, 68, 0.2);
            border: 1px solid rgba(239, 68, 68, 0.3);
            border-radius: 10px;
            color: #fca5a5;
            font-size: 0.9rem;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .logout-btn:hover {
            background: rgba(239, 68, 68, 0.3);
            border-color: rgba(239, 68, 68, 0.5);
            transform: translateY(-2px);
            box-shadow: 0 5px 20px rgba(239, 68, 68, 0.2);
        }

        .dashboard-content {
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            padding-top: 70px;
        }

        .empty-state {
            text-align: center;
            color: rgba(255, 255, 255, 0.4);
        }

        .empty-icon {
            width: 120px;
            height: 120px;
            margin: 0 auto 1.5rem;
            background: rgba(255, 255, 255, 0.03);
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            border: 2px dashed rgba(255, 255, 255, 0.1);
        }

        .empty-icon svg {
            width: 50px;
            height: 50px;
            stroke: rgba(255, 255, 255, 0.2);
        }

        .empty-title {
            font-size: 1.5rem;
            margin-bottom: 0.5rem;
            color: rgba(255, 255, 255, 0.6);
        }

        .empty-text {
            font-size: 1rem;
        }

        .floating-orbs {
            position: fixed;
            inset: 0;
            pointer-events: none;
            overflow: hidden;
        }

        .orb {
            position: absolute;
            border-radius: 50%%;
            filter: blur(60px);
            opacity: 0.3;
            animation: orbFloat 30s infinite ease-in-out;
        }

        .orb:nth-child(1) { width: 400px; height: 400px; background: #8b5cf6; top: -200px; left: -200px; }
        .orb:nth-child(2) { width: 300px; height: 300px; background: #3b82f6; bottom: -150px; right: -150px; animation-delay: -10s; }
        .orb:nth-child(3) { width: 350px; height: 350px; background: #ec4899; top: 50%%; right: -175px; animation-delay: -20s; }

        @keyframes orbFloat {
            0%%, 100%% { transform: translate(0, 0); }
            50%% { transform: translate(50px, 50px); }
        }
    </style>
</head>
<body>
    <div class="floating-orbs">
        <div class="orb"></div>
        <div class="orb"></div>
        <div class="orb"></div>
    </div>

    <nav class="navbar">
        <div class="logo">3D Glass Auth</div>
        <div class="user-section">
            <span class="user-email">%s</span>
            <form method="POST" action="/logout" style="margin: 0;">
                <button type="submit" class="logout-btn">Sign Out</button>
            </form>
        </div>
    </nav>

    <main class="dashboard-content">
        <div class="empty-state">
            <div class="empty-icon">
                <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                </svg>
            </div>
            <h2 class="empty-title">Welcome to your Dashboard</h2>
            <p class="empty-text">Your workspace is empty. Start building something amazing!</p>
        </div>
    </main>
</body>
</html>`, email)
}
