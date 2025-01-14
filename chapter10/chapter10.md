# 学んだこと
- ファイル変更の監視のAPIはLinuxだとinotify系、BSD系OSはkqueue
- syscall.Flock()でファイルのロックができるが、これは強制力のないアドバイザリーロック。Windowsはこれが利用できず、代わりにLockFileEx()を使うが、これは強制ロック
- 他のプロセスの読み込みは許可する共有ロックと、全てを許可しない排他ロックがある
- syscall.Mmap()ではファイルをメモリに展開して、書き換えた内容をそのままファイルに反映もできる
  - ファイルを指定しないと、メモリだけ確保できる。mmapシステムコールはメモリ確保に使われている
  - コピーオンライトという複数のプロセスで同じファイルをmmapした時に、書き込みが起こるまではメモリを共有して、書き込みが起きたら、コピーする仕組み

## 同期とブロッキング
1スレッドでI/Oを処理する時という前提
- 同期・ブロッキング: I/Oをリクエストしたら、結果返ってくるまでスレッドの処理を止めて待つ(read, writeのシステムコール)
- 同期・ノンブロッキング: I/Oをリクエストしたら、結果を待たない。結果は再度I/Oリクエストして、こちらから問い合わせる。(read, writeをノンブロッキングのフラグを付与して呼ぶ)
- 非同期・ブロッキング: I/Oをリクエストしたら、結果を待たない。結果はチャネルで届くのを待つ。I/O多重化と呼ばれる。イベント駆動モデル。(select, kqueue)
- 非同期・ノンブロッキング: メインとは別のスレッドで処理を行い、結果もそっちのスレッドで受け取る。あまり使われてない。欠点を解消したio_uringがLinuxに実装された。(非同期I/0)
  - io_uring: 処理を依頼するキューと結果を受け取るキューがカーネルに用意される

## ファイルシステム
- FUSEはシステムコールをユーザーランドのプロセスに転送することができて、これのインターフェースを実装することで、S3などを使ってオリジナルのファイルシステムを作れる。オリジナルのファイルシステムをマウントできる。

# 面白かったこと
- Golangでは同期・ブロッキングのI/Oが基本だが、goroutine, channelを使えば、非同期やノンブロッキングを簡単に実装できる
- Golangは気軽にgoroutineを作れるので、ブロッキング・同期I/Oでも問題がないが、スレッドを作成するコストが高い言語だと、ノンブロッキングはパフォーマンス的に重要になってくる
- システムコールのselect属とchannelのselect構文は、準備ができたタスクの通知を受け取って、処理を行うという点で同じ
- ファイルシステムを作成できれば、Webサービスから得られる情報をJSONやCSVのファイルとしてみせるようなインターフェースも実現できる
  - インターフェースはファイルシステムだけど、裏側ではHTTPリクエストしてる的な、面白い

# わからなかったこと
- 勧告ロックで問題ないのかな？様々な開発者によって開発されたアプリケーションをLinux上で動かしたりすることは普通にありそうで、非協調的プロセスが混じったりしそうだけど。

# やること
- [] select属システムコールで同期・ノンブロッキングの実装
- [] オリジナルのファイルシステムのマウント