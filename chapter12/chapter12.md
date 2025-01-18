# 学んだこと
- プロセスはプログラムの実行単位
- プロセスは入出力をファイル単位で行い、管理するテーブルがある。それのインデックスがファイルディスクリプタ
- プロセスと外界とのやり取りはシステムコール経由。プロセスが単体でできるのは数値計算くらいで、データの入出力も、時刻の取得もシステムコール経由。
- シェルは実行パスからプログラムを特定して、引数をパースして、環境変数を渡して、プロセスを起動できる
- プロセスは誰かしらのユーザー権限で実行される。ユーザーIDなどを持つ
- プロセスのユーザーIDは親プロセスのものを引き継ぐが、実効ユーザーIDは個別に設定できる。実効ユーザーIDはそのプロセスのアクセス権限
- プロセスの入力はコマンドライン引数と環境変数で、出力は終了コード。終了コードは親プロセスに返される。0なら正常。waitのシステムコールで子プロセスの終了コードを待つ。
- OSはプロセスをプロセスディスクリプタで管理してる。プロセスごとにカレントディレクトリ、メモリ領域、ファイルディスクリプタなど。
- Goのプロセスからコマンドやプログラムを指定して、プロセスを作って実行できる。環境変数や実行時のディレクトリなどを設定できる。stdinPipeなどを使って子プロセスの標準出力・入力とやり取りできる
- ANSIエスケープシーケンスを使うことで、端末エミュレータでプロセスの出力に色をつけたりできる。

# 面白かったこと
- ファイルとソケットはファイルディスクリプタを扱うという点では共通化されており、read, write, closeなどのシステムコールは共通化されてるが、完全にファイルとして抽象化されてるわけではなく、それぞれでしか使えないシステムコールもある。
- MacOSではユーザーやグループのアクセス権限をApple Open Directoryで管理してる。ファイルのアクセス権限は別。
- OSがプロセスを実行した時点で、標準入力、標準出力、標準エラー出力のファイルが開かれてる
- 他のプロセスの標準出力を受け取って、加工して、ターミナルに表示する的な時、擬似端末と認識されず、ANSIエスケープシーケンス情報が省かれてしまうことがあるので、擬似端末と詐称できるライブラリなどがある
- Linuxカーネルではプロセスもスレッドもまとめてタスクとして扱っている。スレッドとプロセスは実装上の差はほぼない。スレッドは親のプロセスとメモリを共有できるプロセスとして扱ってる。
- Linuxにはメモリ共有もフラグで管理できるcloneシステムコールがあり、内部実装ではforkもスレッド生成もこれと同じロジックで処理される
- BusyBoxはコマンドのプログラムが全て1つの実行ファイルに入ってる。コマンドごとにリンクを用意して、どのリンク経由で実行ファイルが呼ばれたかによって、処理内容を変えている。

# パイプ
- パイプはプロセスの標準入出力を他のプロセスにすること
- パイプで繋げたコマンドのプロセスはプロセスグループ(別名ジョブ)になる
- 子プロセスを作ってプログラムを実行する時に、そのプロセスの標準入出力のパイプを受け取って、繋ぐことができる
  - プロセス間通信というよりは標準入出力をもらうって感じかな

# リダイレクト
- リダイレクトはプロセスの標準入出力を他のファイルにすること
- openシステムコールでファイルを開いて、dup2で標準入出力をそれと繋げる

# forkと並行処理
- forkによってプロセスをそのままコピーしたサブプロセスを生やすことができる。しかし、実行したスレッド以外はコピーされないので注意が必要。GCなどは別のgoroutin(OSスレッド)で実行されたりしてるので、Goで使うのは現実的じゃない。
- forkを実行しただけだとそのままコピーされるだけなので、execでプログラムやコマンドライン引数や環境変数を新しく読み込む。
- スクリプト言語だとグローバルインタプリタロックやジャイアントVMロックなどの制約により、マルチスレッドが現実的じゃないことがあり、この場合はforkを使ったマルチプロセスは結構使われる。
  - グローバルインタプリタロックはメモリ管理をスレッドセーフにするための機能で、1プロセスで1スレッドに限定する。ロックを取ってこれを実装してる。ただ、I/O待ちの場合はこのロックを解放する。I/Oが多いなら、マルチスレッドにする余地があるかも。
- forkの際は親子のどちらかでメモリの変更が入るまでは、メモリを共有するコピーオンライトによりメモリの節約をしていたが、PythonやRubyなどでGCなどにより即座にメモリの変更が入ってしまい、相性が悪かったりして、その対応が言語側に入ったりした。
- execで実行するコードを別のコードに置き換えることができる

# deamon
- deamonとはサーバープロセスなどのバックグラウンドプロセスを作るための機能。シェルから実行した場合は、ログアウトしたりシェルを閉じたらプログラムの実行が止まってしまうため、必要。
- deamon化はセッションやグループから切り離した後、forkしてinitシステムを親にするのと、標準入出力との切り離しを行う
  - ターミナルから実行されたコマンドやプログラムのプロセスはターミナルのセッションに属する。ターミナルを閉じると、同じセッションのプロセスにSIGHUPシグナルが送られて、全て終了させられる。バックグラウンドプロセスでも同様みたい。nohupなどで無視することはできる。
  - パイプで繋いだものは同じプロセスグループに属する
  - ターミナルでシェルから実行した場合、コマンドのプロセスの標準出力がターミナルになってるおり、ファイルディスクリプタ1はターミナルのデバイスである"/dev/tty"を指してるが、`/dev/null`を指すようになる
    - `/dev/null`とかは仮想的なファイルなので、プロセス同士で衝突とかはしない
    - deamon化とは別の話だが、`/dev/stdout`に書き込みをすると、標準出力されるらしい。`dev/null`だったら捨てられる
- Golangはforkの実行が現実的ではないため、systemdなどを使うことが多い

# やること
- [x] 外部コマンドの実行
- [x] パイプの実装
- [x] リダイレクトの実行
- [x] 終了コードを受け取る
- [x] プロセスを終了させる

# わからなかったこと
