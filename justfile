alias re := reload-edge

# reload NGINX
reload-edge:
    docker compose kill edge -s SIGHUP
