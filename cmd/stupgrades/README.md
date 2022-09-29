# Deployment

This is currently deployed manually as a container instance.

```
% docker build -t docker.io/syncthing/stupgrades:latest -f Dockerfile.stupgrades .
% docker push docker.io/syncthing/stupgrades:latest
% az container create --subscription d241969e-6dfe-4a70-b3cf-81dfdcb1a5b3 \
    -g Websites --name upgrades --image docker.io/syncthing/stupgrades:latest \
    --cpu 1 --memory 1 --ports 8080 --dns-name-label upgrades
```
