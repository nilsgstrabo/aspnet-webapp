# Building for multiple architectures

````
docker buildx build -t radixdev.azurecr.io/nilssimple:latest --platform linux/arm64 --push .

docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t radixdev.azurecr.io/nilssimple:latest --push .
```