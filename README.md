# Autoheal

```
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣄⠙⢷⣦⡀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⡀⢷⡀⠻⣶⣤⡈⠁⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣶⡖⠀⣤⣤⣤⣀⠙⣷⡀⠛⠦⣄⣀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢀⣠⡶⠀⣿⣿⡇⢸⣿⣿⣿⣿⣦⣈⠙⠷⣦⣤⡤⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢀⣴⣿⣿⠇⢸⣿⣿⠀⣾⣿⣿⡿⠿⠿⢿⠀⢠⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣰⣿⣿⣿⣿⠀⡾⠿⠛⢀⣉⣠⣤⣤⣤⣤⣤⣤⣄⣀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢰⣿⣿⡿⠟⠋⣀⣤⣶⣿⣿⣿⣿⣿⣿⡿⠿⠿⠿⠿⠿⠆⠀⠀⠀⠀⠀
⠀⠀⠀⡿⠛⢉⣠⣶⣿⣿⣿⠿⠛⠋⣉⣁⣤⣤⣤⣤⣴⠀⣤⡄⠀⠀⠀⠀⠀⠀
⠀⠀⢀⣤⣾⣿⣿⠿⠋⠉⣤⣶⡆⢸⣿⣿⣿⣿⣿⣏⠁⠀⣻⡇⠀⠀⠀⠀⠀⠀
⠀⠀⣾⣿⡿⠋⣁⣴⣿⠀⣿⣿⡇⠸⣿⣿⣿⣿⣿⠟⢠⣿⡿⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠹⠋⢠⣾⣿⣿⣿⡆⢸⣿⣿⠀⣿⣿⠟⠿⠋⣠⣿⡿⠁⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠈⠻⣿⣿⡿⠇⠈⣿⣿⡆⠘⠛⣁⡀⡺⠟⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠉⠓⠲⠆⢸⣿⣿⡀⠛⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
          AUTOHEAL
Restart unhealthy containers
```

WIP tool that monitors and restarts unhealthy docker containers.  
Inspired by [willfarrell/docker-autoheal](https://github.com/willfarrell/docker-autoheal) but rewritten in golang.

Autoheal only monitor and restarts containers labeled with "autoheal : true" that are part of the docker compose 
project passed as parameter.

## Install
```
wget -qO- https://github.com/napicella/autoheal/install | bash
```

Or run the `install` script, which by defaults builds a docker image called `autoheal:latest`.

## Example usage (docker)
```bash
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  autoheal:latest \
  -project myproject -verbose
```

## Example usage (docker compose)

```yaml
version: "3.9"

services:
  web:
    image: nginx:latest
    container_name: web
    labels:
      autoheal: "true"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 10s
      timeout: 5s
      retries: 3

  autoheal:
    image: autoheal:latest
    container_name: autoheal
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      COMPOSE_PROJECT_NAME: myproject
    command: ["-project", "my-docker-compose-project", "-verbose"]
```
