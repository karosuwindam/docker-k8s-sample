# Slackにデータをポストするだけのマイクロサービス

## はじめに
SlackのAPIを事前に作成しておく \
`Bot`の`file:write`や`chat:write`を有効にしておく

## API使用

|API|補足|説明|
|--|--|--|
|/api/v1/uploadfile||ファイルをアップロードする enctype="multipart/form-data"を設定|
|--|file|送信ファイル|
|--|filename|送信ファイル名指定|
|/api/v1/postmessage||メッセージを送信する|
|--|message|送信メッセージ指定|


## 環境変数

|環境変数|説明|
|--|--|
|SLACK_TOKEN|SLACK botのトークン|
|SLACK_CHANNEL|slackの送信チャンネル|
|WEB_PORT|起動webサーバのポート指定|
|WEB_IP|起動webサーバの受信IP指定|

## 参考URL

https://qiita.com/RuyPKG/items/5ac07ddc04432ee7641b