# Building for multiple architectures

Build for single:
```
docker buildx build -t docker.io/nilsgustavstrabo/nilssimple:latest --platform linux/arm64 --push .
```

Build for multiple:
```
docker buildx build --platform linux/amd64,linux/arm64 -t docker.io/nilsgustavstrabo/nilssimple:latest --push .
```

Inspect an existing image:
```
docker buildx imagetools inspect docker.io/nilsgustavstrabo/nilssimple:latest
```