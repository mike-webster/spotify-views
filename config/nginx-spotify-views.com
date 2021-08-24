server {
        server_name spotify-views.com www.spotify-views.com;
        root /home/spotify-views/releases-client/live;
        index index.html;
        # Other config you desire (TLS, logging, etc)...
        location / {
                try_files $uri /index.html;
        }

    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/spotify-views.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/spotify-views.com/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot


}

server {
        server_name api.spotify-views.com;
        location / {
                proxy_pass http://127.0.0.1:3001/;
        }

    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/spotify-views.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/spotify-views.com/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

}

# Virtual Host configuration for example.com
#
# You can move that to a different file under sites-available/ and symlink that
# to sites-enabled/ to enable it.
#
#server {
#       listen 80;
#       listen [::]:80;
#
#       server_name example.com;
#
#       root /var/www/example.com;
#       index index.html;
#
#       location / {
#               try_files $uri $uri/ =404;
#       }
#}

server {
    if ($host = www.spotify-views.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


    if ($host = spotify-views.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


        server_name spotify-views.com www.spotify-views.com;
        listen 80 default_server;
    return 404; # managed by Certbot




}

server {
    if ($host = api.spotify-views.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


        server_name api.spotify-views.com;
        listen 80;
    return 404; # managed by Certbot


}