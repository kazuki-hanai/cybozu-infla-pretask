# cybozu-infla-pretask

## 概要
- 3つの関数を実装し、計測を行った
- `processSingle()`
  - 1つのスレッドで単純にファイルの各行からsha256のチェックサムを計算し、標準出力を行う
- `processConcurrent1()`
  - `processSingle()` にgoroutineを組み合わせた
  - sha256のチェックサム計算と標準出力をgoroutineに切り離した
  - goroutineの数が増えすぎると標準出力を行う際にgoroutine同士のlock待ちが起きプログラムが強制終了する
- `processConcurrent2()`
  - `processConcurrent1()` に対して、以下の変更を行ったgoroutineの数とチェックサムの標準出力順序が正しくなるように実装したもの
    - gochannelを利用して、goroutine間での情報のやりとりを行う
    - goroutineを利用すると、チェックサムの正しい終了順序が保証されなくなる
      - そのため、バッファを利用して、正しい順序で標準出力が行われるようにした

## 計測結果
- 以下のような行を持つ1.8GBのファイルに対して`processSingle()` と `processConcurrent2()` を実行した
- また、実行時間の計測は `time go run ./main.go > /dev/null` のようなコマンドで行った
```
bPlNFGdSC2wd8f2QnFhk5A84JJjKWZdKH9H2FHFuvUs9Jz8UvBHv3Vc5awx39ivu
wsp2nChCIwVQztA2n95rXrtzhwuSAd6heDZ0tHBxFq6Pysq3N267L1vqkgnBsUje
9FqBZonjaaWDcXMm8biABkerSuHpnMmMDF2EsjYyTQWCfIuilZxV2FCniRwo7StO
fGOILa0u1wXnEw1GDGuvdSewj77Ax7Tlfj84Qyu6uRn8CTECWzT5s4ZJHd0TxrtM
...
```
- `processSingle()`
```
go run ./main.go > /dev/null  11.46s user 4.26s system 100% cpu 15.597 total
go run ./main.go > /dev/null  11.30s user 4.08s system 103% cpu 14.904 total
go run ./main.go > /dev/null  11.29s user 4.08s system 102% cpu 14.950 total
go run ./main.go > /dev/null  11.42s user 4.11s system 103% cpu 15.081 total
go run ./main.go > /dev/null  11.58s user 4.25s system 101% cpu 15.592 total

平均実行時間: 15.224秒
```
- `processConcurrent2()`
  - `WORKER_NUM=1` で計測した
    - チェックサムを計算するgoroutineの数
```
go run ./main.go > /dev/null  24.77s user 8.15s system 228% cpu 14.393 total
go run ./main.go > /dev/null  24.70s user 8.06s system 233% cpu 14.005 total
go run ./main.go > /dev/null  24.94s user 8.18s system 231% cpu 14.292 total
go run ./main.go > /dev/null  24.93s user 8.19s system 232% cpu 14.245 total
go run ./main.go > /dev/null  24.75s user 8.08s system 237% cpu 13.837 total

平均実行時間: 14.1544
```

## 考察
- ファイル読み込み、チェックサムを計算、標準出力をgoroutineとして切り離すことで、処理時間を短縮することができた
  - 標準出力を行う処理にプログラムの多くの処理が費やされているため、`processConcurrent2()`で実装したような切り離しを行うと処理速度が向上する
    - `processSingle()`で標準出力のみを除いて計測を行うと、平均5.6秒程度で処理が終わった
    - 標準出力のバッファ処理を導入することで、処理時間の短縮を図ることが可能だと考えられる
- `WORKER_NUM=1`の結果が最も高速だった
  - 実行したプログラムではCPUのコアを1つのみ利用していた
  - 複数のgoroutineは1つのコアを共有するため、処理待ち状態のgoroutineが発生する
  - 複数コアを利用した処理を実装することで、高速化が可能だと考えられる