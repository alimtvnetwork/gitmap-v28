#!/bin/bash

# Exit on error
set -e
echo "Md Alim Ul Karim"
echo "https://alimkarim.com | https://riseup-pro.com"
echo
echo "========== before joining worker node run crio services start and enable ==========="
echo "sudo systemctl enable crio.service && sudo systemctl start crio.service"
echo "=================== (run from sudo su - user only) initializing kubeadam init =============================================="
sudo systemctl enable crio.service && \
sudo systemctl start crio.service && \
kubeadm init && \
mkdir -p $HOME/.kube && \
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
chown $(id -u):$(id -g) "$HOME/.kube/config" && \
echo "========== Remember to join from Worker as `sudo su` root user ===========" && \
echo "========== If required run in the master node ===========" && \
echo "  $ sudo kubeadm token create --print-join-command" && \
echo "========== If required run network reset to start fresh ===========" && \
echo "  $ kubeadm reset --force" && \
echo "========== after initializing - in the worker nodes ===========" && \
echo "  $ sudo hostnamectl set-hostname \"new_hostname\"" && \
echo "========== after joining worker node - run kubectl apply on master node ===========" && \
echo "  $ kubectl apply -f https://reweave.azurewebsites.net/k8s/v1.31/net.yaml" && \
echo "  $ kubectl -n kube-system get pods"
