package main

import (
    "database/sql"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strings"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/joho/godotenv"
)

var (
    dbUser        string
    dbPassword    string
    dbName        string
    telegramToken string
    chatID        string
)

func init() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    dbUser = os.Getenv("DB_USER")
    dbPassword = os.Getenv("DB_PASSWORD")
    dbName = os.Getenv("DB_NAME")
    telegramToken = os.Getenv("TELEGRAM_TOKEN")
    chatID = os.Getenv("CHAT_ID")
}

func main() {
    // Parse command-line flags
    domainPtr := flag.String("t", "", "Domain to monitor")
    flag.Parse()

    if *domainPtr == "" {
        log.Fatal("You must provide a domain using the -t flag")
    }

    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Get new subdomains using multiple tools
    subdomains := getSubdomains(*domainPtr)
    for _, subdomain := range subdomains {
        if !subdomainExists(db, subdomain) {
            saveSubdomain(db, *domainPtr, subdomain)
            sendTelegramAlert(subdomain)
        }
    }
}

func getSubdomains(domain string) []string {
    var subdomains []string

    // Sử dụng Subfinder
    subfinderSubdomains := runSubfinder(domain)
    subdomains = append(subdomains, subfinderSubdomains...)

    // Sử dụng Findomain
    findomainSubdomains := runFindomain(domain)
    subdomains = append(subdomains, findomainSubdomains...)

    // Sử dụng Crt.sh
    crtshSubdomains := runCrtsh(domain)
    subdomains = append(subdomains, crtshSubdomains...)

    return unique(subdomains)
}

func runSubfinder(domain string) []string {
    var subdomains []string

    cmd := exec.Command("subfinder", "-d", domain)
    output, err := cmd.Output()
    if err != nil {
        log.Fatal(err)
    }

    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if line != "" {
            subdomains = append(subdomains, line)
        }
    }

    return subdomains
}

func runFindomain(domain string) []string {
    var subdomains []string

    cmd := exec.Command("findomain", "-t", domain)
    output, err := cmd.Output()
    if err != nil {
        log.Fatal(err)
    }

    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if line != "" {
            subdomains = append(subdomains, line)
        }
    }

    return subdomains
}

func runCrtsh(domain string) []string {
    var subdomains []string

    // Gửi yêu cầu GET đến crt.sh
    url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", domain)
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    var result []struct {
        NameValue string `json:"name_value"`
    }

    if err := json.Unmarshal(body, &result); err != nil {
        log.Fatal(err)
    }

    seen := make(map[string]struct{})
    for _, item := range result {
        names := strings.Split(item.NameValue, "\n")
        for _, name := range names {
            name = strings.TrimSpace(name)
            if strings.HasSuffix(name, domain) && name != domain {
                if _, exists := seen[name]; !exists {
                    seen[name] = struct{}{}
                    subdomains = append(subdomains, name)
                }
            }
        }
    }

    return subdomains
}

func unique(slice []string) []string {
    uniqueMap := make(map[string]struct{})
    for _, item := range slice {
        uniqueMap[item] = struct{}{}
    }

    var uniqueSlice []string
    for key := range uniqueMap {
        uniqueSlice = append(uniqueSlice, key)
    }

    return uniqueSlice
}

func subdomainExists(db *sql.DB, subdomain string) bool {
    var exists bool
    err := db.QueryRow("SELECT COUNT(*) > 0 FROM subdomains WHERE subdomain = ?", subdomain).Scan(&exists)
    if err != nil {
        log.Fatal(err)
    }
    return exists
}

func saveSubdomain(db *sql.DB, domain, subdomain string) {
    _, err := db.Exec("INSERT INTO subdomains (domain, subdomain) VALUES (?, ?)", domain, subdomain)
    if err != nil {
        log.Fatal(err)
    }
}

func sendTelegramAlert(subdomain string) {
    bot, err := tgbotapi.NewBotAPI(telegramToken)
    if err != nil {
        log.Fatal(err)
    }

    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("New subdomain found: %s", subdomain))
    _, err = bot.Send(msg)
    if err != nil {
        log.Fatal(err)
    }
}
