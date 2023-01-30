# Go Net Gen

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
