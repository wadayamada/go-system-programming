# 学んだこと
- シグナルはプロセス間通信とソフトウェア割り込みのために使われる
- プロセスを強制終了するSIGKILLとプロセスを一時停止するSIGSTOPはアプリケーションではハンドルできない
- SIGTERMはkillシステムコールがデフォルトで送信するシグナルでサーバーアプリケーションでハンドルした方が良い
- SIGHUPはコンソールアプリケーションが擬似端末から切断される時に呼ばれる。サーバーアプリケーションはこの目的で呼ばれることはないため、設定ファイルの再読み込みを外部から指定する目的で使うことがデファクトスタンダードらしい。
- SIGINTはコンソールアプリケーションでユーザーがCtrl+C
- グレイスフルリスタートはユーザー影響なしでサーバーをリスタートするやつ。Server::StarterとHttp.Server.Shutdownを組み合わせるとできる

# 面白かったこと
- シグナルは0除算エラーやメモリの範囲外アクセスなどがCPUで発生して、それをカーネルが受けて、シグナルを生成することもある
- Dockerではコンテナを終了させる時はまずSIGTERMを送るようになってる。いきなりSIGKILLで強制終了するのではなく、まずSIGTERMを送るのがお行儀が良い
- シグナルを受け取ると任意のスレッドで既存の処理を中断して、シグナル処理が強制的に走る。それは危険なので、Goではシグナル処理用のスレッドを用意している。goroutineが特定のスレッドで実行されるようにしてる。それ以外のスレッドはシグナルをブロックするようにしてる。

# わからなかったこと
- グレイスフルリスタートの実装楽しそう
- Go言語ランタイムにおけるシグナルの内部実装がむずい