unzip myapp_mysql.zip && cd myapp
go mod tidy
go run ./cmd/server



cd /var/www/api-sportbook
go mod tidy
go build -o api-sportbook ./cmd/server

sudo nano /etc/systemd/system/api-sportbook.service

[Unit]
Description=API Sportbook Go App
After=network.target

[Service]
User=www-data
WorkingDirectory=/var/www/api-sportbook
ExecStart=/var/www/api-sportbook/api-sportbook
EnvironmentFile=/var/www/api-sportbook/.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target


sudo systemctl daemon-reload
sudo systemctl enable api-sportbook
sudo systemctl start api-sportbook
sudo systemctl status api-sportbook


sudo nano /etc/nginx/sites-available/api-sportbook

server {
    listen 80;
    server_name api-sportbook.2m-sy.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

sudo ln -s /etc/nginx/sites-available/api-sportbook /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx