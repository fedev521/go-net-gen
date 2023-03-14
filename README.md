# Go Net Gen

## Setup

- `terraform init && terraform apply`
- download service account key from GCP
- `cd secrets && export GOOGLE_APPLICATION_CREDENTIALS="$PWD/gonetgen-sa-key.json"`

## Boilerplate

Features:

- Dockerfile that leverages cache to avoid unnecessary downloads, compilations
  and tests. It uses a distroless Debian image. Build with `DOCKER_BUILDKIT=1
  docker build -t "<imagename>:latest" .`
- logging: uses [logur](https://github.com/logur/logur) as a facade (and
  adapter) and [logrus](https://github.com/sirupsen/logrus) under the hood.
  Always uses JSON formatter.
- configuration: uses [viper](https://github.com/spf13/viper). Convention:
  mandatory configuration file named `config.toml`.
- flags: with [pflag](https://github.com/spf13/pflag)
- a `run()` function to facilitate tests.

## Resources

- [install d2](https://d2lang.com/tour/install)

D2 quick start:

```bash
echo 'x -> y' > helloworld.d2
d2 -w helloworld.d2 out.png
```

Useful:

```powershell
d2 \\wsl.localhost\Ubuntu-20.04\home\federico\projects\personal\go-net-gen\cmd\app\out.d2 out.png
```

## Next

- read information from file rather than from GCP (test)
- exclude projects or VMs
- don't print "Project" if project name is too long
- try another diagram with shared VPC relationships
- remove dead code
- support load balancers
- containerize
- CI/CD integration
- improve visually
- fix TODOs
- improve configuration management
- generate service account key in Terraform without storing it state: [blueprint](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/tree/v19.0.0/blueprints/cloud-operations/onprem-sa-key-management)
