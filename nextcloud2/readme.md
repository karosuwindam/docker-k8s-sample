削除エラー解消方法


maintenanceモードに入る

select * from oc_file_locks;

SELECT * FROM `oc_file_locks` WHERE `lock` = 1
DELETE FROM oc_file_locks WHERE `lock` = 1

https://piano2nd.smb.net/PukiWiki/index.php?nextcloud+file+unlock