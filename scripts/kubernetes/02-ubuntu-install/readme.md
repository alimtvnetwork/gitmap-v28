# Ubuntu Scripts

## `01-set-ip.sh`

`Failed to restart kube-apiserver.service: Unit kube-apiserver.service not found.`

```bash
 chmod +x ./01-set-ip.sh && \
 sudo ./01-set-ip.sh <ip> <route ip optinal / 192.168.0.1>
 ```

```bash
 chmod +x ./01-set-ip.sh && \
 sudo ./01-set-ip.sh 192.168.0.66
 ```

```bash
 chmod +x ./01-set-ip.sh && \
 sudo ./01-set-ip.sh 192.168.0.67
 ```

```bash
chmod +x ./01-set-ip.sh && \
sudo ./01-set-ip.sh 192.168.0.68
```

## `05-remove-used-ports.sh`

`Failed to restart kube-apiserver.service: Unit kube-apiserver.service not found.`

```bash
 chmod +x ./05-remove-used-ports.sh && \
 sudo ./05-remove-used-ports.sh 
 ```

## To run `04-kill-user-processes.sh`

```bash
chmod +x 04-kill-user-processes.sh && \
./04-kill-user-processes.sh <username>
```

```bash
chmod +x ./04-kill-user-processes.sh && \
./04-kill-user-processes.sh kube
```

## To run `05-omy-zsh-only.sh`

```bash
chmod +x ./05-omy-zsh-only.sh && \
./05-omy-zsh-only.sh <Theme> <is_apped_zsh:yes>
./05-omy-zsh-only.sh robbyrussell
./05-omy-zsh-only.sh amuse
./05-omy-zsh-only.sh dst
```


## To run `01-zsh-theme-change-v2.sh`

```bash
chmod +x ./01-zsh-theme-change-v2.sh && \
./01-zsh-theme-change-v2.sh <Theme> <is_apped_zsh:yes>
./01-zsh-theme-change-v2.sh -h
./01-zsh-theme-change-v2.sh robbyrussell
./01-zsh-theme-change-v2.sh amuse
./01-zsh-theme-change-v2.sh dst
```

## To run `09-create-root-user-v2.sh`

```bash
chmod +x 09-create-root-user-v2.sh && \
./09-create-root-user-v2.sh <new-username> <password> <homedir>
```

```bash
chmod +x ./09-create-root-user-v2.sh && \
./09-create-root-user-v2.sh masterK awe1233401212
```

```bash
chmod +x ./09-create-root-user-v2.sh && \
./09-create-root-user-v2.sh worker1 awe1233401212
```

```bash
chmod +x ./09-create-root-user-v2.sh && \
./09-create-root-user-v2.sh worker2 awe1233401212
```

```bash
chmod +x ./09-create-root-user-v2.sh && \
./09-create-root-user-v2.sh masterK awe1233401212 amuse && \
./09-create-root-user-v2.sh worker1 awe1233401212 half-life && \
./09-create-root-user-v2.sh worker2 awe1233401212 josh && \
./09-create-root-user-v2.sh worker3 awe1233401212 dirty && \
./09-create-root-user-v2.sh worker4 awe1233401212 jnrowe
```

```bash
chmod +x ./09-create-root-user-v2.sh && \
sudo ./09-create-root-user-v2.sh masterK awe1233401212 amuse && \
sudo ./09-create-root-user-v2.sh worker1 awe1233401212 half-life && \
sudo ./09-create-root-user-v2.sh worker2 awe1233401212 josh && \
sudo ./09-create-root-user-v2.sh worker3 awe1233401212 dirty && \
sudo ./09-create-root-user-v2.sh worker4 awe1233401212 jnrowe
```

## Script 16 - Root user v2 (with omyzsh)

```bash
chmod +x ./09-create-root-user-v2.sh && \
sudo ./09-create-root-user-v2.sh masterK awe1233401212 amuse && \
sudo ./09-create-root-user-v2.sh worker1 awe1233401212 half-life && \
sudo ./09-create-root-user-v2.sh worker2 awe1233401212 josh && \
sudo ./09-create-root-user-v2.sh worker3 awe1233401212 dirty && \
sudo ./09-create-root-user-v2.sh worker4 awe1233401212 jnrowe && \
sudo ./09-create-root-user-v2.sh test1 awe1233401212 dst

```
## List available Users

```bash
cat /etc/passwd
cut -d: -f1 /etc/passwd
getent group sudo
```

## To run `07-auto-purge.sh`

```bash
chmod +x ./07-auto-purge.sh && \
sudo ./07-auto-purge.sh
```

### `masterK` hostname


```bash
hostnamectl set-hostname masterK
```

### `worker1` hostname

```bash
hostnamectl set-hostname worker1
```

### `worker2` hostname

```bash
chmod +x ./09-kube-install-csa-kuril-v3.sh && \
sudo ./09-kube-install-csa-kuril-v3.sh worker2
```

```bash
hostnamectl set-hostname worker2
```


## Once install is done do `kubeadam init` from only the master node

### ⚠️ ! From Master Node Only run the below script

```bash
chmod +x ./10-kube-init.sh && \
sudo ./10-kube-init.sh
```

## Change ZSH Theme

```bash
chmod +x ./01-zsh-theme-change-v2.sh && \
sudo ./01-zsh-theme-change-v2.sh {theme-name} {true / false -- append zshrc}
```

```bash
chmod +x ./01-zsh-theme-change-v2.sh && \
sudo ./01-zsh-theme-change-v2.sh half-life
```

## 17 Repo Permission

```bash
sudo chmod +x 09-repo-permissions.sh && \
./09-repo-permissions.sh
```
## 18 Git Pull

```bash
sudo chmod +x 10-git-pull.sh && \
./10-git-pull.sh
```
