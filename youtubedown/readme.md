# youbutu downloadツール

## 概要
Youtubeのダウンロード用のツール


## 適用

以下のファイル
```
git clone https://github.com/karosuwindam/youtubedown.git
```

```
kubectl apply -f youtubedown/ns.yaml
kubectl -n youtube-down create secret generic smbcreds --from-literal username=USERNAME --from-literal password="PASSWORD"
kubectl apply -f youtubedown/deployment.yml
```
