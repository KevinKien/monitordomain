
## Install 
- Install other tool
```
go get github.com/joho/godotenv
go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install github.com/projectdiscovery/puredns/v2/cmd/puredns@latest
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

## Running
- Run with command
```
./monitordomain -t example.com
```
