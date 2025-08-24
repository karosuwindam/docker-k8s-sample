# インストールについて

Raspberry Pi Imagerでユーザや無線LANなどの初期設定が可能なので設定を実施する

ubuntu server 24.10で作成するものとする

## 初期設定

ブートの修正
```
sudo sed -i \ 
's/$/ cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory/g' \ 
/boot/firmware/cmdline.txt 

```


以下の項目について設定を実施する。

* IPアドレスの固定

ubuntu serverの場合は`sudo vi /etc/netplan/99-static-fix.yaml`

```
network:
  version: 2
  renderer: networkd
  ethernets:
    eth0:
      dhcp4: false
      dhcp6: true
      addresses: [192.168.0.40/24]
      routes:
        - to: default
          via: 192.168.0.1
      nameservers:
        addresses: [192.168.0.1]
  wifis:
    wlan0:
      dhcp4: false
      addresses: [192.168.0.50/24]
      routes:
        - to: default
          via: 192.168.0.1
      nameservers:
        addresses: [192.168.0.1]
      access-points:
        "aterm-7c9b4e-a":
          auth:
            key-management: "psk"
            password: "058d2948e39271dc610e825bbcf36f1b7fcbc3fc9edefbad42c130fec5dd1a22"
```

以下のコマンドで、変更する

```
sudo chmod 600 /etc/netplan/99-static-fix.yaml
```

## 設定

* swapの無効はデフォルトで設定されているので問題ない

/etc/hostsを使用する場合は再起動ごとに以下のファイルをベースに書き換えるため以下のファイルを書き換える

sudo vi /etc/cloud/templates/hosts.debian.tmpl

```
sudo sh -c "cat <<EOF >> /etc/cloud/templates/hosts.debian.tmpl
192.168.0.50    k8s-master-01
192.168.0.60    k8s-worker-01
192.168.0.61    k8s-worker-02
192.168.0.62    k8s-worker-03
EOF
"
```

sudo systemctl restart cloud-init-local.service

### iptablesの設定

以下コマンドを実行してネットワークアダプタを設定します。

```
sudo sh -c "cat <<EOF > /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF
"
```

再起動なしに有効する場合は、以下のコマンドを実行します。

```
sudo sysctl --system
```

### containerdのインストール

実際にインストールする場合は次の手順を実行します。

```bash
sudo apt install containerd runc -y
```

http通信のレジストリを登録するため以下のコマンドでcontainerdの設定ファイルを作成する。なおこのコマンドで作成されるものはversion 3のためversion 2用と間違えないようにする

```
sudo mkdir /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
```

kubernetesのネットワークを構成するためには、overlayとbr_netfilterのモジュールを読み込む必要があるため以下のコマンドで起動時に読み込む設定を実行する。

```bash
sudo sh -c "cat <<EOF > /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF
"
```

コンテナ内の動作状態を確認するためのcrictlコマンドの設定で`/etc/crictl.yaml`ファイルを作成して以下の通り登録を実施する

```
sudo sh -c "cat <<EOF > /etc/crictl.yaml
runtime-endpoint: unix:///var/run/containerd/containerd.sock
image-endpoint: unix:///var/run/containerd/containerd.sock
timeout: 10
debug: false
EOF
"
```

以上で、kubernetesでcontainerdを使用する設定は完了です。

kubernetesのインストール準備

```
sudo apt-get install -y apt-transport-https ca-certificates curl gpg
```

公開鍵の登録

```
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key | \
sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
```

リポジトリの登録

```
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.32/deb/ /' | \
sudo tee /etc/apt/sources.list.d/kubernetes.list
```

コマンドインストールとバージョン固定

```
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
```

## nfsやsamba設定

nfsの接続はデフォルトではソフトが含まれていないので、以下のコマンドでインストール

```
sudo apt install nfs-common
```

sambaに接続するには以下のコマンドでインストール

```
sudo apt install samba-common cifs-utils jq
```

sudo apt install nfs-common samba-common cifs-utils jq  -y 

## Ubuntu Server 25.04の場合

cgroupのmemoryが有効でないので以下の通り修正する

cat /proc/cmdline

sudo vi /boot/firmware/cmdline.txt

の末尾に以下の文章を追加

cgroup_memory=1 cgroup_enable=memory

### DNSサーバを作成する

Ubuntu Serverは、systemd-resolvedが動いているポートが他のDNSを動かすたためポート重なっているので、
dnsmasqを使用してDNSサーバを構築することができないので、以下のコマンドで、止める必要がある

```
sudo systemctl disable --now systemd-resolved
```

以下のコマンドで、dnsmasqをインストールするただ、先にsystemd-resolvedを止めるとインストールできなくなるので、以下のコマンドでインストール後止める

```
sudo apt install dnsmasq -y
```

インストール後、以下コマンドでDNSサーバを起動する

```
sudo systemctl enable --now dnsmasq
```

ただこのままだと/run/systemd/resolve/resolv.confがないためcniプラグインが動かない可能性があるので、以下のファイルを編集します。

/etc/dnsmasq.conf

を編集して、参照先を/etc/dnsmasq-upstreom.confに変更

/etc/dnsmasq-upstreom.confに上位のDNSサーバを記載する

/etc/resolv.confが参照リンクのため削除して、自分自身を参照するようにする

上記設定完了後resolvやdnsmasqを起動する

### Firewallについて

Ubuntu Searverはufwと呼ばれるソフトFirewallが動いているので、Kubernetesクラスタを構築するためには
ポートを開放させるか、ufwというソフトを止めてFirewallを止める必要があります。
なお、Firewallを止める場合は以下のコマンドを実行してください

```
sudo systemctl disable --now ufw.service
```

### 独自のコンテナレジストリを有効にする

Ubuntu Server の場合は(version 3)設定方法は以下のように登録すると証明がないコンテナレジストリであっても使用することが可能

sudo vi /etc/containerd/config.toml

```
    [plugins.'io.containerd.cri.v1.images'.registry]
      config_path = ''
      [plugins."io.containerd.cri.v1.images".registry.mirrors]
        [plugins."io.containerd.cri.v1.images".registry.mirrors."192.168.0.91:5000"]
          endpoint = ["http://192.168.0.91:5000"]
```

上記設定を有効にするにはcontainerdのサービスを再起動させる必要がある

sudo systemctl restart containerd