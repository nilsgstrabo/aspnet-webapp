# Building for multiple architectures

Build for single:
````
docker buildx build -t radixdev.azurecr.io/nilssimple:latest --platform linux/arm64 --push .
```

Build for multiple:
```
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t radixdev.azurecr.io/nilssimple:latest --push .
```

Inspect an existing image:
```
docker buildx imagetools inspect radixdev.azurecr.io/nilssimple:latest
```