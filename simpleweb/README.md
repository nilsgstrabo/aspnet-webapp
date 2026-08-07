# Building for multiple architectures

Build for single:
```
docker buildx build \
	--build-arg VERSION=$(git describe --tags --exact-match HEAD 2>/dev/null || git rev-parse --short HEAD) \
	-t docker.io/nilsgustavstrabo/nilssimple:latest \
	--platform linux/arm64 \
	--push .
```

Build for multiple:
```
docker buildx build \
	--build-arg VERSION=$(git describe --tags --exact-match HEAD 2>/dev/null || git rev-parse --short HEAD) \
	--platform linux/amd64,linux/arm64 \
	-t docker.io/nilsgustavstrabo/nilssimple:latest \
	--push .
```

Inspect an existing image:
```
docker buildx imagetools inspect docker.io/nilsgustavstrabo/nilssimple:latest
```