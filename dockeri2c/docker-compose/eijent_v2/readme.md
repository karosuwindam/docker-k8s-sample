## 機能
以下のコマンドでセンサーチェック対象のリセット送信
```
curl localhost:9140/reset -X POST
```

## 履歴
* ver 0.14.5
  * goを1.22に変更
  * コムポートで読み込みタイムアウト設定されていなかったので修正
  * コンテナを64bitへ変更
* ver 0.14.4
  * 温度センサーの最大値と最小値を超える場合はエラー表示
* ver 0.14.3
  * 読み取り動作がおかしい場合の処理
* ver 0.14.2
  * 読み取り動作がおかしい場合の処理
* ver 0.14.1
  * CO2センサーのヘルスチェック
* ver 0.14
  * bme280の読み取り失敗したときに値を維持する。
* ver 0.13
  * センサーの効率化
  * CO2センサーの異常時の排除設定
* ver 0.12
  * 待ち時間が発生するので修正の必要あり
  * DHTのセンサーによる待ち時間がでかいので修正する必要あり
* ver 0.11
  * コードを作り直し
  * リセットボタンの追加
* ver 0.10
  * サーバ側にcontext機能を追加
  * CO2センサーのログ出力を整理
* ver 0.9
  *  加速度センサMA8452Qについて対応
* ver 0.85
  * ビルド用のgoのバージョンを1.18.2に変更
* ver 0.84
  * 温度センサーDhtとbme280で湿度と温度が-1を検出したら更新しないように変更
* ver 0.83
  * CO2センサーをMH-19cが動くものに変更
* ver 0.8
  * BME280を組み込む
* ver 0.7
  * CO2センサーを組み込む
* ver 0.6
  * dhtのセンサー機能を組み込む
  * healthチェックとjson出力を組み込む
  * ソースコード整理
* ver 0.2
  * tsl2561について読み取り機能を組み込む
  * I2cが初回で検出できない場合はエージェント動かさないようにする。


## 参考
* BME280 \
 https://www.switch-science.com/catalog/2236/ \
 https://akizukidenshi.com/download/ds/bosch/ BST-BME280_DS001-10.pdf \
 [Raspberrypi で使用する](https://deviceplus.jp/hobby/raspberrypi_entry_039/)

* 光センサー TSL2561 \
 https://wiki.seeedstudio.com/Grove-Digital_Light_Sensor/
 http://www.ne.jp/asahi/shared/o-family/ElecRoom/AVRMCOM/TSL2561/TSL2561.html
 https://www.switch-science.com/catalog/1801/

## 実行参考
 https://www.denshi.club/pc/raspi/5raspberry-pi-zeroiot27-1-i2c-tls2561.html

* MMA8452Q搭載 三軸加速度センサモジュール \
 https://www.switch-science.com/catalog/1927/