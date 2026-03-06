https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/getting-started.html

```bash
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia \
    && helm repo update
```

```bash
helm install --wait gpu-operator \
    -n gpu-operator --create-namespace \
    nvidia/gpu-operator \
    --version=v25.10.1 \
    --values values.yaml
```

```bash
helm upgrade gpu-operator \
    -n gpu-operator \
    nvidia/gpu-operator \
    --version=v25.10.1 \
    --values values.yaml
```