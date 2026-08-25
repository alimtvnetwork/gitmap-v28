# Imperative Commands

```bash
kubectl create deployment nginx-deployment --image=nginx:latest --replicas=2 --dry-run=client -o yaml > nginx-deployment.yaml
kubectl apply -f nginx-deployment.yaml

kubectl create service loadbalancer nginx-service --tcp=80:80 --dry-run=client -o yaml > nginx-service.yaml
kubectl apply -f nginx-service.yaml




kubectl create deployment mysql-deployment --image=mysql:5.6 --dry-run=client -o yaml > mysql-deployment.yaml

kubectl apply -f mysql-deployment.yaml


kubectl create service clusterip mysql-service --tcp=3306:3306 --dry-run=client -o yaml > mysql-service.yaml
kubectl apply -f mysql-service.yaml



kubectl create deployment wordpress-deployment --image=wordpress:4.8-apache --dry-run=client -o yaml > wordpress-deployment.yaml
kubectl apply -f wordpress-deployment.yaml

kubectl set env deployment/wordpress-deployment WORDPRESS_DB_HOST=mysql-service:3306 WORDPRESS_DB_PASSWORD=password --dry-run=client -o yaml >> wordpress-deployment.yaml




kubectl create service loadbalancer wordpress-service --tcp=80:80 --dry-run=client -o yaml > wordpress-service.yaml
kubectl apply -f wordpress-service.yaml


kubectl get all --all-namespaces -o wide
kubectl get nodes
kubectl get pods
kubectl get pv
kubectl get pvc

kubectl get nodes --all-namespaces -o wide
kubectl get pods --all-namespaces -o wide
kubectl get pv --all-namespaces -o wide
kubectl get pvc --all-namespaces -o wide
kubectl get svc --all-namespaces -o wide
```