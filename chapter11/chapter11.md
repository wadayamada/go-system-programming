# 学んだこと
- シェルはユーザーがOSを扱うためのインターフェース。ユーザーとシステムの接点
  - FinderもDesktopもシェル
- コマンドシェルはユーザーがテキスト画面で入力して、コマンドを起動するためのもの。Unix系OSだと、ユーティリティプログラムは色々あるが、シェルはそれを起動するだけなので、シェル本体の機能はシンプル。
- MacOSのTerminalは端末エミュレータ。擬似端末を通してシェルに繋ぐ。zshはシェル。
- 環境変数はシェルから子プロセスにまとまった情報を渡すために用意されている機能
- シェルがコマンドを実行するまで
  - 入力を受ける。テキスト補完や履歴などもある
  - テキストの分解。コマンドと引数。パイプやリダイレクトの記号も意識した処理が必要。
  - コマンドと引数の前処理。環境変数やコマンドの実行結果の埋め込み
  - 実行ファイルの探索。環境変数のPATHを使う
  - 引数のワイルドカードの展開
  - プロセスの起動とコマンド同士の連携。パイプとかリダイレクト

# 面白かったこと
- シェルが存在しないOSもある
  - シェルの役割はプログラムを選んで実行することで、実行するプログラムが決まってるなら不要になる
  - シェルがないとセキュリティ的な安全性が高い
- `./script.sh`, `bash script.sh`で実行する場合は別プロセスで呼ばれるが、source や .で実行すると、元のシェルのプロセスで実行される
- 外部コマンドはシェルの外部にある実行可能ファイルのことで、`ls`, `grep`, `find`など。内部コマンドはシェルにビルドインされてるもので、`cd`, `pwd`, `fg`, `which`など。
- プログラムから外部コマンドを実行する方法
  - シェルを経由する方法
  - 直接、子プロセスを作って実行する方法
    - シェルがやってくれる引数の分解などはプログラム側でやる必要がある
- POSIXはSUSと呼ばれるUNIXを名乗るために必要な規格の別名。シェルのコマンド、シェルの言語仕様、システムコールなどが定義されている。LinuxにはLSBという規格がある。
- シェルのコマンドを見てみると、ファイルの管理・加工・表示と実行中のプロセスへの命令などが定義されていることがわかる
- bash, zshはPOSIXのシェルのコマンド要件をほぼ満たしてるので、POSIX互換シェルと呼ばれてるが、WindowsのコマンドプロンプトやPowerShellは互換ではない
- CGIはプログラムに環境変数でリクエストヘッダー情報などを渡していた。マルチプロセスでプロセスごとにリクエストを処理する方式だったため可能
  - CGIは昔使われていた、Perl, PHPのプログラムをNginxから実行するためのインターフェース
  - Perl, PHPでCGIのインターフェースを満たしたCGIプログラムを実装する
  - 1リクエスト1プロセスで処理する
- dockerコマンドはデーモンに対する命令を行っているに過ぎず、環境変数は伝播しないので、引数で渡す必要がある
- .envファイルは環境変数に似てるが、OSの機能ではないので、アプリケーションごとに読み込む対応が必要。Ruby on Railsで発明された。
- リダイレクト
  - `>` : 標準出力をファイルに上書き
  - `>>`: 標準出力をファイルに追記
  - `<`: ファイルを標準出力
- WSLでWindoswsとLinuxはネットワークドライブとして互いのファイルシステムにアクセスしてるため、通常より低速になる。コード補完とかは大量のファイル検索が入るので影響ありそう。

# わからなかったこと
- .envファイルに書かれた環境変数を読み込むコマンドをGolangで実装するの楽しそう。パスも通す。
- PATHを見て、実行ファイルの探索も楽しそう
- whichの実装
