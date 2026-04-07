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

WIP tool that monitors and restarts unhealthy Docker containers.  
Inspired by https://github.com/willfarrell/docker-autoheal but rewritten in Go.

Autoheal monitors and restarts unhealthy containers labeled with `autoheal=true`
and belonging to one or more specified Docker Compose projects.

---

## Install

```bash
wget -qO- https://github.com/napicella/autoheal/raw/refs/heads/main/install | bash
```

Or run the `install` script, which by default builds a Docker image called `autoheal:latest`.

---

## Example usage (docker)

```bash
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  autoheal:latest \
  -project myproject -verbose
```

---

## Example usage (docker compose)

```yaml
version: "3.9"

services:
  web:
    image: nginx:latest
    container_name: web
    labels:
      autoheal: "true"
      autoheal.strategy: "project"  # optional
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
    command: ["-project", "my-docker-compose-project", "-verbose"]
```

---

## Configuration

### `-project`

Comma-separated list of Docker Compose projects to monitor.

Only containers belonging to these projects will be considered for auto-healing.

```bash
-project myproject
-project projectA,projectB,projectC
```

This matches the Docker label:

```
com.docker.compose.project
```

---

### `-verbose`

Enable verbose logging.

---

### `-restart-limit`

Maximum number of restarts before stopping a container.

Default: `10`

---

### `-stop-timeout`

Timeout (in seconds) before force-stopping a container.

Default: `10`

---

### `-interval`

Interval between health checks.

Default: `5s`

---

## Labels

Autoheal relies on container labels to decide what to monitor and how to restart.

---

### `autoheal=true`

Required for a container to be monitored.

```yaml
labels:
  autoheal: "true"
```

---

### `autoheal.strategy=project`

Optional.

When set to `project`, autoheal will restart the **entire Docker Compose project**
instead of just the unhealthy container.

```yaml
labels:
  autoheal: "true"
  autoheal.strategy: "project"
```

#### Behavior

- Default (no strategy or different value)  
  → Restart only the unhealthy container

- `autoheal.strategy=project`  
  → Restart the whole compose project (using the original compose file)

This uses the Docker label:

```
com.docker.compose.project.config_files
```

to locate the compose file.

---

## How it works

- Listens to Docker events for `health_status: unhealthy`
- Filters containers by:
  - `autoheal=true`
  - matching Compose project(s)
- Applies restart strategy:
  - container-level restart (default)
  - project-level restart (if configured)

---

## Notes

- Containers **must define a healthcheck** to be monitored
- Only containers with `autoheal=true` are considered
- Compose project filtering is required via `-project`