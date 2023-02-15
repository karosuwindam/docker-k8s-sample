# Rook
## 概要
分散ストレージを実行するなお、jobの動作が64bitに固定できないので、テストに失敗した
すべて、64bitにすれば動作には問題ない


## データについて
wgetを手に入れる
```
wget https://raw.githubusercontent.com/rook/rook/master/deploy/examples/crds.yaml
wget https://raw.githubusercontent.com/rook/rook/master/deploy/examples/common.yaml
wget https://raw.githubusercontent.com/rook/rook/master/deploy/examples/operator.yaml
wget https://raw.githubusercontent.com/rook/rook/master/deploy/examples/cluster.yaml
wget https://raw.githubusercontent.com/rook/rook/master/deploy/examples/cluster-on-pvc.yaml
wget https://raw.githubusercontent.com/rook/rook/master/deploy/examples/cluster-test.yaml
```

以下のファイルを編集
* operator.yaml
affinityを追加して、arm64とamd64のみ実行するようにする。また、コンテナのイメージを`v1.10.0`を指定する
configファイルで、`CSI_CEPHFS_PLUGIN_NODE_AFFINITY`を有効に



## 実行手順
以下のコマンドでOperatorのコントローラを導入する
```
kubectl apply -f crds.yaml -f common.yaml -f operator.yaml
```
以下のコマンドでカスタムリソースを適用してクラスタを作成
```
kubectl apply -f cluster.yaml
```

# 参考
* [HomePage](https://rook.io/)
* [QuickStart](https://rook.io/docs/rook/v1.10/Getting-Started/quickstart/#prerequisites)
* [Github](https://github.com/rook/rook)
* [rook/ceph image](https://hub.docker.com/r/rook/ceph)
