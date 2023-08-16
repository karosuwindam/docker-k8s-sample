## 概要
koel

## 起動後以下のコマンドを実行
php artisan koel:init --no-assets
php artisan koel:init

email: admin@koel.dev
password: KoelIsCool

sambaを使用する場合は以下のコマンドでパスワードを登録する
```
kubectl -n koel create secret generic smbcreds --from-literal username=USERNAME --from-literal password="PASSWORD"
```

## 参考
https://koel.dev/

https://hub.docker.com/r/phanan/koel/