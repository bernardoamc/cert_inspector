# cert_inspector

Small Go utility that wraps [gungnir](https://github.com/g0ldencybersec/gungnir) to scan certificate transparency (CT) logs.

It is designed to:
- Automatically restart `gungnir` if it stops
- Create daily logs in `logs/` directory
- Separate logs for standard output and errors
- Provide all the infrastructure required to run it locally or through Docker

---

## 🚀 Getting Started

```bash
make build             # builds ./build/cert_inspector
make run               # builds & runs locally (adds logs/ and targets.txt, needs gungnir)
make docker-build      # builds the Docker image
make docker-run-local  # runs the Docker container with files shared by host
make docker-run-volume # runs the Docker container with volumes
make copy-logs         # creates a logs_backup on host
make clean             # removes ./build
make tidy              # runs go mod tidy
```

If you have to run some of these commands with `sudo`, you might need to do something like:

```bash
sudo PWD=$PWD make <command>
```

---

### Debugging the container

```bash
docker run --rm -it cert_inspector:latest sh
```

### Targets

Let's say you want to monitor subdomains from every Bug Bounty program.

1. Download the list provided by Project Discovery
2. Aggregate domains using `jq`
2. Create a targets.txt file with domains

```bash
curl -O https://raw.githubusercontent.com/projectdiscovery/public-bugbounty-programs/main/chaos-bugbounty-list.json
jq -r '.programs[].domains[]' chaos-bugbounty-list.json > targets.txt
```
