# グラフを作るだけのマイクロサービス

gnuplotを利用してグラフを返す

|api||説明|
|--|--|--|
|/help|||
|/api/v1/data|||


gnuplot -p -e "set terminal png;set output \"test.png\";plot 'gf.data' using 0:2:xtic(1) with linespoints"