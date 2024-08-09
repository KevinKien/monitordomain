
## Install 
- Install other tool
```
go get github.com/joho/godotenv
go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install github.com/projectdiscovery/puredns/v2/cmd/puredns@latest
```

- Install findomain
```
wget https://github.com/Findomain/Findomain/releases/download/v7.0.0/findomain-linux
chmod +x findomain-linux
mv findomain-linux /usr/local/bin/findomain
```

- Install database
```
CREATE DATABASE subdomain_monitor;
USE subdomain_monitor;

CREATE TABLE subdomains (
    id INT AUTO_INCREMENT PRIMARY KEY,
    domain VARCHAR(255) NOT NULL,
    subdomain VARCHAR(255) NOT NULL,
    UNIQUE KEY (subdomain)
);
```

## Complie file
- With linux, macos
```
go build -o httpsmonitor
```

- With windows
```
GOOS=windows GOARCH=amd64 go build -o httpsmonitor.exe .
```

## Running
- Run with command
```
./monitordomain -t example.com
```
