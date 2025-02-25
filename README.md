[Goならわかるシステムプログラミング](https://www.lambdanote.com/products/go-2)を読む

読んだ内容のメモや感想やをチャプターごとに書いた

# まとめ

## 実行ファイルを作る
goプログラムをコンパイルしてリンクすることで、実行ファイルを作ることができる。	
コンパイルはgoプログラムを機械語のコードに変換する
リンクは機械語のコードやランタイムを結合して実行ファイルを作る

リンカー
* 実行ファイルのヘッダーに最初に実行されるエントリーポイントのアドレスを書く。ランタイムの_start_がエントリーポイントになる。
* 各セクションをメモリのどこに置くかも書く

ランタイム
* OSへのシステムコールを呼ぶ
* goroutine, channel, スライスなどのデータ構造

実行ファイルはセクション単位でデータを扱う。実行コードの領域、静的に確保された初期化済みのメモリ領域、変数が置かれる領域

## コマンドシェルからプログラムを実行する

zshなどのコマンドシェルはTerminalなどの端末エミュレータから接続することができる。

シェルで外部コマンドやプログラムを実行できる。実行すると、PATHから実行ファイルを探して実行する。

Golangプログラムが実行されると、実行ファイルに従ってメモリにセクションを配置し、エントリーポイントから実行を開始する。最初はgoroutine, シグナルハンドラ, GCなどの初期化が実行される

インタプリタでプログラムを実行する方式もある
* CPUはインタプリタを実行し、インタプリタがプログラムを読み込んでアプリケーションを実行する
* Java, Rubyだと実行頻度が高い部分をネイティブコードに変換することで高速化するJITコンパイラという工夫がある。初回だけ時間がかかる。

コマンドやプログラムはシェルの子プロセスとして実行され、標準出力はシェルにパイプされる。標準入出力をプロセスに繋ぐのがパイプで、ファイルに繋ぐのがリダイレクト。

プロセスに対する引数は環境変数として渡される。CGIで1プロセス1リクエストで捌いていた時代はHTTPリクエストのヘッダーなどが環境変数で渡されていた。

ターミナルから実行されたコマンドはターミナルのプロセスのセッショングループになり、ターミナルを閉じるとSIGHUPシグナルが送られて終了させる。またパイプしてるプロセスは同じプロセスグループに属する。

終了させたくなかったら、セッションとグループから切り離して、forkしてinitシステムを親プロセスにして、標準入出力をターミナルから切り離してdeamon化する必要がある。forkは実行したスレッド以外はコピーされないため、golangでは現実的ではなく、deamon化するならsystemdなどを使うのが良い。

また、execでプロセスが実行するコードを他のコードに置き換えることができる。

## システムコール

OSはプロセス管理、プロセス間通信、ファイルシステム、メモリ管理、ネットワーク、タイマーなどを担っている

アプリケーションプロセスでストレージやメモリの入出力、ネットワーク通信、タイマー、プロセス間通信などを行いたい場合、システムコールを呼ぶ必要がある。プロセスが単体でできるのは数値計算くらい。標準入出力、ソケット、乱数生成などはファイルのRead/Writeとして統一したインターフェースで利用できる。ファイルのIDがファイルディスクリプタ。

Linuxでシステムコールを呼ぶと最終的に命令セットのSYSCALL命令が呼ばれる。SYSCALL命令が呼ばれると、CPUはレジスタにアドレス登録済みのLinuxカーネルのシステムコール関数を呼ぶ。レジスタ経由でシステムコール番号も渡す。WindowsとMacOSはシステムコールに標準ライブラリ(libc)を使ってる。

## ファイルシステム

論理ボリュームマネージャーが複数のストレージを1つの論理ボリュームとして抽象化する

ファイルシステムは論理ボリュームなどを使って、ファイルを管理する

VFSはマウントされたファイルシステムを管理して、統一したインターフェースを提供する。システムコールを捌くのもVFS。

ファイルの書き込みはメモリバッファに書き込めた時点で完了として返す。読込もメモリバッファが最新だったら、そこから読み込む

## メモリ管理

仮想メモリから物理メモリへのマッピングを管理してるのがページテーブル。ページテーブルのCPUキャッシュがTLB

共通ライブラリはプロセスごとにロードするのではなく、プロセスで共有されてる

Mmapを対象ファイルを設定せずに実行するとメモリの確保だけ行える。物理メモリが仮想メモリに割り当てられる。

Goで変数のデータをスタックに置くかヒープに置くかは自動で決まる。関数の中でしか使わないならスタックに置かれる。ポインタを他の関数に渡すなどする場合はヒープに置かれる。スタックの方が高速。

ヒープに変数を割り当てる場合にはGolangではTCMallocが呼ばれる。スレッドローカルキャッシュから取得し、なければ共有ヒープから取得し、それもなければ、Mmapする

オブジェクトキャッシュによって、malloc済みのオブジェクトを使いまわせる。mallocをせずに変数の割り当てができる。

GoのGCはstop the worldの時間が無視できるくらい短くなるように最適化されている

## ネットワーク

アプリケーション層
* HTTPリクエストは複数のTCPセグメントに分けて送っても問題ない
    * ヘッダーは行ごとに分けて送れる
    * ボディは事前にContentLengthを指定するか、チャンクにすることで、分けて送れる
* TLS
    * サーバーから渡された公開鍵を使ってクライアントは暗号化して送信する
    * 証明書がサーバーから渡される。認証局の署名があるので、別でもらった認証局の公開鍵で署名を検証する

トランスポート層
* TCP
    * 好きなタイミングでお互いにwriteもreadもできる。リクエスト・レスポンスという概念を作ってるのはHTTP
    * 輻輳制御、順序制御、再送制御、ウィンドウ制御とフロー制御
    * 3way handshake
* UDP
    * コネクションレスなので、いきなり送りつけることができる
    * マルチキャストができる
    * DNS, NTP, WebRTCなどに使われている
    * HTTP/3のためのQUICに使われる
* Unix Domain Socket
    * ローカルでしか使えない代わりに高速
    * カーネル内部で完結する高速なネットワークインターフェースを使う
    * カーネルのバッファに書き込んで、アプリケーションプロセスにコピーするくらいの負荷しかかからない

インターネット層
* IPアドレスとIPパケットでデータの送受信

ネットワークインターフェース層
* Ethernet, Wifiなどを使って、Macアドレスとフレームでデータの送受信

## シグナル

プロセス間通信には、ソケット、パイプ、共有メモリ以外に、シグナルやメッセージキューなどがある

シグナルはCPU割り込みハンドリングの結果受信にも使われる。CPUの0除算エラー、メモリの範囲外アクセス、killシステムコール、ctrl+cによって、CPU割り込みが発生し、CPU割り込みハンドラが呼ばれ、カーネルが対象プロセスにシグナルを送る。

システムコールを叩いて、シグナルを送ることもできる。 シグナルを受け取ると、プロセスが一時停止され、シグナルハンドラーが実行される。シグナルハンドラーはスレッドで実行中の処理を強制的に中断して処理が始まるため、ロックしている最中にシグナルが入ると、デッドロックの危険性がある。そのため、マルチスレッド環境ではシグナル処理用スレッドを用意するのが良い。

プロセスを強制終了するSIGKILL、一時停止するSIGSTOPはアプリケーション側でハンドリングできない。 SIGTERMはkillシステムコールがデフォルトで呼ぶやつで、アプリケーションでハンドリングした方が良い。 SIGINTはctrl+cで呼ばれる

## タイマー

タイマー、カウンターはタスクのCPUへの割り当て、I/Oスケジューリング、タイムアウト、ファイルのタイムスタンプなどで使われる

ハードウェアタイマーで適宜修正するリアルタイム時刻、調整をせず巻き戻しが派生しないモノトニック時刻がある

また、普段人間が使ってるのがウォールクロック時間で、CPUの実行時間がCPU時間

電源を切っても消えないのがリアルタイムクロックというハードウェアのクロックで、それを元にOSで保持してるのがシステムクロック

タイムスタンプカウンタはクロック周波数をカウントしたもので分解能が高く、タイマーデバイスは定期的に割り込みを入れて、プロセスのタイムスライスを減らしたりするのに使ってる
Linuxでは現在時刻取得をvDSOを使って高速化してる。Linuxはカーネルのメモリに書き込んで、ユーザープロセスはそこから読み込むやり方。

## 並列処理

CPUにおける処理時間が大きい時は並列、I/O待ちの場合は並行で処理する 並列化してもコア数以上はスケールしないし、CPUが元からパンパンだったら意味なし
### マルチプロセス

Python, Rubyだとスレッドセーフにするために、グローバルインタプリタロックで1プロセス1スレッドの制約がある

並列処理をしたいならマルチプロセスにする必要があるが、コンテキストスイッチのコストが高い

プロセスプールでプロセス作成コストは抑えることができる

### マルチスレッド

プロセスほどではないがOSスレッドも起動コストがかかる

コンテキストスイッチコストはプロセスと変わらない？

スレッドプールなどの工夫はできる

### goroutine

goroutineはOSスレッドに比べて初期スタックメモリのサイズが小さいので、起動処理が軽い

プログラムカウンタ、スタックポインタ、データレジスタの3つのレジスタを退避するだけでタスクを切り替えられるようになっていて、カーネルに処理を渡す必要がないため、コンテキストスイッチコストが低い

goroutineで並列化して、channelで直列化できる

### スレッドセーフ

同時に実行されるとrace conditionが起きる可能性があるコードをクリティカルセクションと呼ぶ

sync.Mutexでロックを取るのが基本だが、sync/atomicの不可分操作でロックフリーでスレッドセーフに処理が行える

メモリのアクセス待ちを吸収するためにCPUにはメモリへの書き込み順序を変えるアウトオブオーダーという仕組みがある。シンプルな処理なら基本的に問題ないが、Mmapを使ったメモリマップドI/Oでハードウェアの制御をしている場合、実行順序が重要。sync/atomicを使うとメモリバリアという機能で実行順序も保証してくれるらしい。

## 同期・ブロッキング

* 同期・ブロッキング
    * I/Oをgoroutineの処理を止めて待つ
* 同期・ノンブロッキング
    * I/Oを待たない。後ほど問い合わせる。
* 非同期・ブロッキング
    * I/Oをリクエスト時は待たない。後ろの方でselectで待つ。ノンブロッキングで複数I/Oリクエストして、selectでまとめて待つのがよくあるやつ。I/O多重化
* 非同期・ノンブロッキング
    * ノンブロッキングでI/Oリクエストして、別のgoroutineで結果を待つ

## 仮想化・コンテナ

ホストOSの上にハイパーバイザを入れて、ゲストOSを入れることで仮想化できる

CPUを完全にエミュレーションする方式は他のCPU用のソフトウェアも動かせるが、パフォーマンスは落ちる。同じアーキテクチャに限定するネイティブ方式の方がパフォーマンスは良い。

素朴に実装するとゲストOSが特権命令を実行するたびにハイパーバイザの処理が入ってしまう。最近のCPUには仮想化支援機能というのが用意されており、VT-xではゲストOSが特権命令を実行すると、ハイパーバイザ用OSモードでCPUを実行でき、エミュレーションが不要になる。

Proxmoxは仮想化技術が統合されたOSで、Linuxカーネルに埋め込まれたKVMというハイパーバイザを使っている。またCPUの仮想化支援機能も使ってる。

ゲストOSが仮想化を意識して、ハイパーバイザのシステムコールを呼ぶ、準仮想化というのもある

コンテナはCPU、メモリ、ネットワーク、デバイスへの使用量を制限するコントロールグループと、プロセス、ネットワーク、ファイルシステム、ユーザーなどの名前空間を分離できる名前空間というlinuxの機能で実現できる Linux前提なので、MacOSではRancher Desktopなどを使って仮想化環境が必要。Windowsはコンテナをサポートしてる。