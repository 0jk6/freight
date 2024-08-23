# Freight

Allows you to run code on Kubernetes using k8s jobs, just like leetcode, codeforces and other online judges


Build docker images for the backend, service and the ui and push them to the dockerhub

- `docker build -t YOUR-DOCKERHUB-USERNAME/freight-backend:0.0.1 -f Dockerfile.backend .`
- `docker build -t YOUR-DOCKERHUB-USERNAME/freight-service:0.0.1 -f Dockerfile.service .`

- `docker push YOUR-DOCKERHUB-USERNAME/freight-backend:0.0.1`
- `docker push YOUR-DOCKERHUB-USERNAME/freight-service:0.0.1`

Make sure that you edit the images in the install.yml file to match the above pushed images.

- `kubectl apply -f install.yml`