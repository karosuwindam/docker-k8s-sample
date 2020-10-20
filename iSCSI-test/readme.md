#iSCSIの設定メモ

```
sudo apt install tgt

sudo mkdir /var/lib/iscsi_disks/
sudo dd if=/dev/zero of=/var/lib/iscsi_disks/disk01_10G.img count=0 bs=1 seek=10G

sudo vi /etc/tgt/conf.d/target01.conf

# 新規作成
# 複数デバイスを提供する場合は以下の<target>～</target>を増やして同じ要領で設定する
# 命名規則：[ iqn.年-月.ドメイン名の逆:任意の名前 ]
<target iqn.2020-10.local:k8s-worker-1.target01>
    # iSCSIターゲットとして提供するデバイス (複数設定する場合は [backing-store ***] の行を追加)
    backing-store /var/lib/iscsi_disks/disk01_10G.img
    # 接続を許可するiSCSIイニシエーターのIQN (複数設定する場合は [initiator-name *.*.*.*] の行を追加)
    initiator-name iqn.2020-10.local:bookserver2.initiator01
    # 接続を許可する際の認証情報 ( username, password は任意のものを設定)
    incominguser karosu windam1314
</target> 

sudo systemctl restart tgt

#状態確認
sudo tgtadm --mode target --op show
```

# クライアント設定
```
sudo apt -y install open-iscsi xfsprogs

sudo vi /etc/iscsi/initiatorname.iscsi

#InitiatorName=iqn.1993-08.org.debian:01:7d4298bc6bf
InitiatorName=iqn.2020-10.local:bookserver2.initiator01

sudo vi /etc/iscsi/iscsid.conf

#以下のコマンドがraspberrypiでは通らないので、本体再起動
#sudo systemctl restart iscsid open-iscsi
sudo /etc/init.d/open-iscsi restart
sudo /etc/init.d/iscsid restart

sudo iscsiadm -m discovery -t sendtargets -p 192.168.0.23

sudo iscsiadm -m node --login

sudo iscsiadm -m node --logout

sudo mkfs.xfs /dev/sdb 
```