# 学んだこと
- Unix Domain ソケットはコンピュータ内部でしか使えない代わりに高速
- カーネル内部で完結する高速なネットワークインターフェースを作って接続する。TCP, UDPの両方に対応してる。
- サーバーがファイルを作って、そのファイルパスを指定してクライアントが接続する

# 面白かったこと
- Windowsには元々はなくて、名前付きパイプで同様の機能が提供されてたが、Windows10から使えるようになった
- Unixドメインソケットはカーネルのバッファに書き込んで、サーバーの方にコピーするくらいの負荷しかかからないため、たとえばあるベンチマークだとTCPの80倍くらい高速だった

# わからなかったこと
- UDPの場合、クライアント用のUnixドメインソケットファイルが必要な理由がわからない