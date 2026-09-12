package cmdvhost

const wordpressVHostTemplate = `server {
    listen {{.Port}};
    listen [::]:{{.Port}};
    server_name {{.Domain}}{{if .Aliases}} {{.Aliases}}{{end}};
    root {{.DocumentRoot}};
    index index.php index.html index.htm;
    access_log /var/log/nginx/{{.Domain}}_access.log;
    error_log /var/log/nginx/{{.Domain}}_error.log warn;
    client_max_body_size {{.MaxBodySize}};
    fastcgi_read_timeout {{.FastCGITimeout}};
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    location / {
        try_files $uri $uri/ /index.php?$args;
    }

    location ~ \.php$ {
        try_files $uri =404;
        fastcgi_split_path_info ^(.+\.php)(/.+)$;
        fastcgi_pass {{.FastCGIPass}};
        fastcgi_index index.php;
        include fastcgi_params;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        fastcgi_param PATH_INFO $fastcgi_path_info;
        fastcgi_buffer_size 32k;
        fastcgi_buffers 16 16k;
        fastcgi_busy_buffers_size 64k;
        fastcgi_temp_file_write_size 64k;
        fastcgi_intercept_errors on;
    }

    location = /wp-config.php { deny all; access_log off; log_not_found off; }
    location = /xmlrpc.php { deny all; access_log off; log_not_found off; }
    location ~* /wp-content/uploads/.*\.php$ { deny all; access_log off; log_not_found off; }
    location ~ /\.(?!well-known).* { deny all; access_log off; log_not_found off; }
    location ~* \.(ogg|ogv|svg|svgz|eot|otf|woff|woff2|mp4|ttf|css|rss|atom|js|jpg|jpeg|gif|png|ico|zip|tgz|gz|rar|bz2|doc|xls|exe|ppt|tar|mid|midi|wav|bmp|rtf)$ {
        expires max;
        log_not_found off;
        access_log off;
        add_header Cache-Control "public, no-transform";
    }
}

`

const laravelVHostTemplate = `server {
    listen {{.Port}};
    listen [::]:{{.Port}};
    server_name {{.Domain}}{{if .Aliases}} {{.Aliases}}{{end}};
    root {{.DocumentRoot}}/public;
    index index.php index.html;
    charset utf-8;
    access_log /var/log/nginx/{{.Domain}}_access.log;
    error_log /var/log/nginx/{{.Domain}}_error.log warn;
    client_max_body_size {{.MaxBodySize}};
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }
    error_page 404 /index.php;
    location ~ \.php$ {
        fastcgi_pass {{.FastCGIPass}};
        fastcgi_split_path_info ^(.+\.php)(/.+)$;
        fastcgi_index index.php;
        include fastcgi_params;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        fastcgi_param DOCUMENT_ROOT $realpath_root;
        fastcgi_buffer_size 32k;
        fastcgi_buffers 16 16k;
        fastcgi_busy_buffers_size 64k;
    }

    location ~ /\.(?!well-known).* { deny all; access_log off; log_not_found off; }
    location ~ /\.env { deny all; return 404; }
}

`

const phpVHostTemplate = `server {
    listen {{.Port}};
    listen [::]:{{.Port}};
    server_name {{.Domain}}{{if .Aliases}} {{.Aliases}}{{end}};
    root {{.DocumentRoot}};
    index index.php index.html index.htm;
    access_log /var/log/nginx/{{.Domain}}_access.log;
    error_log /var/log/nginx/{{.Domain}}_error.log warn;
    client_max_body_size {{.MaxBodySize}};
    fastcgi_read_timeout {{.FastCGITimeout}};
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    location / {
        try_files $uri $uri/ /index.php?$args;
    }

    location ~ \.php$ {
        try_files $uri =404;
        fastcgi_split_path_info ^(.+\.php)(/.+)$;
        fastcgi_pass {{.FastCGIPass}};
        fastcgi_index index.php;
        include fastcgi_params;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        fastcgi_param PATH_INFO $fastcgi_path_info;
        fastcgi_buffer_size 32k;
        fastcgi_buffers 16 16k;
    }

    location ~ /\.(?!well-known).* { deny all; access_log off; log_not_found off; }
    location ~* \.(css|js|jpg|jpeg|gif|png|ico|svg|woff|woff2|ttf|eot)$ {
        expires 30d;
        access_log off;
        add_header Cache-Control "public, no-transform";
    }
}

`

const staticVHostTemplate = `server {
    listen {{.Port}};
    listen [::]:{{.Port}};
    server_name {{.Domain}}{{if .Aliases}} {{.Aliases}}{{end}};
    root {{.DocumentRoot}};
    index index.html index.htm;
    access_log /var/log/nginx/{{.Domain}}_access.log;
    error_log /var/log/nginx/{{.Domain}}_error.log warn;
    client_max_body_size {{.MaxBodySize}};
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    location / {
        try_files $uri $uri/ =404;
    }

    location ~ /\.(?!well-known).* { deny all; access_log off; log_not_found off; }
    location ~* \.(css|js|jpg|jpeg|gif|png|ico|svg|woff|woff2|ttf|eot)$ {
        expires 30d;
        access_log off;
        add_header Cache-Control "public, no-transform";
    }
}

`
