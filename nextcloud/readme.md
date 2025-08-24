削除エラー解消方法


maintenanceモードに入る

select * from oc_file_locks;

SELECT * FROM `oc_file_locks` WHERE `lock` = 1;
DELETE FROM oc_file_locks WHERE `lock` = 1;

https://piano2nd.smb.net/PukiWiki/index.php?nextcloud+file+unlock

## 423 Locked

NextCloudのコンテナ実行

```bash
su -s /bin/bash www-data
php occ maintenance:mode --on
```

mysql系のコンテナで実行

```bash
# ログイン
mysql
mysql -u $MYSQL_USER -p$MYSQL_PASSWORD -D $MYSQL_DATABASE
# 利用しているデータベースを選ぶ
select * from oc_file_locks;

SELECT * FROM `oc_file_locks` WHERE `lock` = 1;
DELETE FROM oc_file_locks WHERE `lock` = 1;

メンテナンスモードをOFF
```
php occ maintenance:mode --off
```