NextCloudで以下のようなエラーが表示される場合は、
「Error while trying to initialise the database: An exception occurred while executing a query: SQLSTATE[HY000]: General error: 4047 InnoDB refuses to write tables with ROW_FORMAT=COMPRESSED or KEY_BLOCK_SIZE」


以下の通り、オプションを追加してコマンドを実行するように
```
command: --transaction-isolation=READ-COMMITTED --binlog-format=ROW --innodb-file-per-table=1 --skip-innodb-read-only-compressed
```