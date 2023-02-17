ver 0.10
サーバ側にcontext機能を追加
CO2センサーのログ出力を整理
ver 0.9
加速度センサMA8452Qについて対応
ver 0.85
ビルド用のgoのバージョンを1.18.2に変更
ver 0.84
温度センサーDhtとbme280で湿度と温度が-1を検出したら更新しないように変更
ver 0.83
CO2センサーをMH-19cが動くものに変更
ver 0.8
BME280を組み込む
ver 0.7
CO2センサーを組み込む
ver 0.6
dhtのセンサー機能を組み込む
healthチェックとjson出力を組み込む
ソースコード整理

ver 0.2

tsl2561について読み取り機能を組み込む
I2cが初回で検出できない場合はエージェント動かさないようにする。


参考
* BME280 \
 https://www.switch-science.com/catalog/2236/ \
 https://akizukidenshi.com/download/ds/bosch/ BST-BME280_DS001-10.pdf \
 [Raspberrypi で使用する](https://deviceplus.jp/hobby/raspberrypi_entry_039/)

* 光センサー TSL2561 \
 https://wiki.seeedstudio.com/Grove-Digital_Light_Sensor/
 http://www.ne.jp/asahi/shared/o-family/ElecRoom/AVRMCOM/TSL2561/TSL2561.html
 https://www.switch-science.com/catalog/1801/

実行参考
 https://www.denshi.club/pc/raspi/5raspberry-pi-zeroiot27-1-i2c-tls2561.html

* MMA8452Q搭載 三軸加速度センサモジュール \
 https://www.switch-science.com/catalog/1927/