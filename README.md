[Goならわかるシステムプログラミング](https://www.lambdanote.com/products/go-2
)を読む

読んだ内容のメモや感想やスニペットを書く

## 掘り下げたいこと
この本を読んでて気になったこと
別の本などで改めて掘り下げたい
- tcpやその下のレイヤーの具体的な実装がちゃんとわかってない
- パイプライニングとHTTP/2がわからない
- HTTP/1.0, HTTP/1.1, HTTP/2, HTTP/3, QUICがちゃんとわかってない
  - 実装方法やライブラリの工夫などはあるけど、結局、HTTPはプロトコルでしかないから、サクッとプロトコルを調べられるようにしたい
- UDPを使ってZoomのようにカメラ映像や音声を送受信できるサービス作ってみたい
- IPパケットの送信がどう実現されてるか理解したい。ルーターとかスイッチとか
- MACアドレス
- CDN
- VPN
- LAN(ローカルネットワーク)、プライベートIPアドレス
- NAT
- datagram型のunixドメインソケットの場合、クライアント用のUnixドメインソケットファイルが必要な理由がわからない
- サーバー側でnet.ListenPacketしてからクライアント側でnet.Dialして、conn.WriteToしたものはクライアントで受け取れたが、新しくサーバー側でnet.Dialでconnを作って、conn.Writeしたものはクライアントで受け取れなかった
- インターフェースはファイルシステムで、内部ではGCP使ってるとか、HTTPリクエストしてるとかのやつやりたい
- おうちK8sと合わせてネットワーク周りを勉強
- goで作るネットワークのやつも合わせてやる
- ソケットとポートの本も合わせてやる
- イーサネット
- 分電盤、電気工事士
- シェルの実装楽しそう
- 統合テストの文脈ででてきた、ネットワークインターフェース？のdocker0とかちゃんと理解したい
- ネットワーク通信もCPUのハードウェア割り込みで処理してる。バッファとかはあるけど。
- コンテキストスイッチは自作OSしないとわからんわ
  - カーネルとユーザーアプリケーションの行き来がめっちゃ多そう
    - システムコール
    - スケジューリング
  - 割り込み
    - ネットワーク通信のたびに割り込みすると思わなかった