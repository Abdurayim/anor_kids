# 🚀 Pre-Deployment Checklist for Anor Kids Bot

## ✅ Files Ready for Git Commit

### **Before Your First Commit:**

1. **Check what will be committed:**
   ```bash
   git status
   ```

2. **Make sure these are NOT showing up:**
   - ❌ `.env` file (contains secrets!)
   - ❌ `*.db` files (your local database)
   - ❌ `anor-kids*` binaries (executables)
   - ❌ `temp_docs/` directory
   - ❌ `.DS_Store` files

3. **Make sure these ARE included:**
   - ✅ All `.go` source files in `cmd/` and `internal/`
   - ✅ `go.mod` and `go.sum`
   - ✅ `.gitignore`
   - ✅ `.env.example`
   - ✅ `README.md` and documentation
   - ✅ Database migration files (`internal/database/migrations/*.sql`)
   - ✅ Fonts (`fonts/*.ttf`)

---

## 📝 Git Commands to Use

### **First Time Setup:**
```bash
cd /Users/abdurayim/Desktop/PROJECTS/anor-kids

# Initialize git (if not already done)
git init

# Add files
git add .

# Check what's being added (IMPORTANT!)
git status

# Create commit
git commit -m "Initial commit: Anor Kids complaint bot"

# Add remote (replace with your repo URL)
git remote add origin https://github.com/yourusername/anor-kids.git

# Push to GitHub
git push -u origin main
```

### **Subsequent Updates:**
```bash
# Check status
git status

# Add specific files or all changes
git add .

# Commit with message
git commit -m "Your commit message here"

# Push to remote
git push
```

---

## 🔒 Security Checklist

Before pushing to GitHub:

- [ ] **CRITICAL**: Verify `.env` is NOT in git:
  ```bash
  git ls-files | grep .env
  # Should return nothing!
  ```

- [ ] **Database**: Verify database files are NOT in git:
  ```bash
  git ls-files | grep "\.db"
  # Should return nothing!
  ```

- [ ] **Binaries**: Verify executables are NOT in git:
  ```bash
  git ls-files | grep "anor-kids"
  # Should only show source directories, not binaries!
  ```

- [ ] **Secrets**: Double-check no API tokens, passwords, or secrets in code

---

## 🚀 Production Deployment Checklist

### **1. Server Setup:**
- [ ] Install Go 1.21+ on production server
- [ ] Install supervisor/systemd for process management
- [ ] Configure firewall (open port 8080 or your chosen port)
- [ ] Setup SSL certificate (if using webhook mode)

### **2. Copy Required Files:**
```bash
# On production server, clone your repo
git clone https://github.com/yourusername/anor-kids.git
cd anor-kids

# Copy fonts
mkdir -p fonts
# Upload DejaVuSans.ttf and DejaVuSans-Bold.ttf to fonts/

# Create .env from template
cp .env.example .env
nano .env  # Edit with your production values
```

### **3. Configure Environment (.env):**
- [ ] Set real bot token from @BotFather
- [ ] Set actual admin phone numbers
- [ ] Set webhook URL (if using webhook mode)
- [ ] Set GIN_MODE=release for production

### **4. Build and Run:**
```bash
# Build for production
go build -o anor-kids-production ./cmd/bot

# Create required directories
mkdir -p temp_docs

# Run directly (for testing)
./anor-kids-production

# OR setup as a service (recommended)
```

### **5. Setup as Systemd Service:**

Create `/etc/systemd/system/anor-kids.service`:
```ini
[Unit]
Description=Anor Kids Complaint Bot
After=network.target

[Service]
Type=simple
User=your_user
WorkingDirectory=/path/to/anor-kids
ExecStart=/path/to/anor-kids/anor-kids-production
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable anor-kids
sudo systemctl start anor-kids
sudo systemctl status anor-kids
```

### **6. Verify Deployment:**
- [ ] Test health endpoint: `curl http://localhost:8080/health`
- [ ] Test bot responds to /start command
- [ ] Test user registration flow
- [ ] Test complaint submission with images
- [ ] Test admin panel access
- [ ] Monitor logs for errors

---

## 📊 Monitoring

### **View Logs:**
```bash
# If using systemd
sudo journalctl -u anor-kids -f

# If running directly
tail -f logs/bot.log
```

### **Check Database:**
```bash
sqlite3 anor_kids_bot.db "SELECT COUNT(*) FROM users;"
sqlite3 anor_kids_bot.db "SELECT COUNT(*) FROM complaints;"
```

---

## 🆘 Troubleshooting

### **Bot not responding:**
1. Check if process is running: `ps aux | grep anor-kids`
2. Check logs for errors
3. Verify BOT_TOKEN is correct
4. Check webhook is set correctly (if using webhook mode)

### **PDF generation failing:**
1. Check fonts directory exists and contains TTF files
2. Check temp_docs directory has write permissions
3. Check server logs for emoji-related errors

### **Database errors:**
1. Check database file permissions
2. Verify SQLite is installed
3. Check disk space

---

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Telegram Bot API](https://core.telegram.org/bots/api)
- [SQLite Documentation](https://www.sqlite.org/docs.html)

---

## ✅ Final Checklist Before Going Live

- [ ] All tests passed locally
- [ ] .gitignore is properly configured
- [ ] No secrets in git repository
- [ ] .env.example is in repo (for documentation)
- [ ] README.md is complete with setup instructions
- [ ] Fonts are in the repository
- [ ] Database migrations are in the repository
- [ ] Production .env is configured correctly
- [ ] Webhook URL is configured (if using webhook mode)
- [ ] Monitoring/logging is setup
- [ ] Backup strategy is in place

---

**Built with ❤️ for Anor Kids Kindergarten**
