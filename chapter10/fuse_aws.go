package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/winfsp/cgofuse/fuse"
)

type AwsFileSystem struct {
	fuse.FileSystemBase
}

func fuse_sample() {

	s3Client := setup_aws_client()

	awsFs := &AwsFileSystem{}
	host := fuse.NewFileSystemHost(awsFs)
	host.Mount("/tmp/awsfs", []string{})

	// バケット内のオブジェクトをリスト
	bucketName := "go-system-programming"
	output, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		log.Fatalf("unable to list objects, %v", err)
	}

	// オブジェクト一覧を表示
	for _, object := range output.Contents {
		fmt.Println(*object.Key)
	}
}

func setup_aws_client() *s3.Client {
	// CSV ファイルを開く
	file, err := os.Open("./go-system-programming-user_accessKeys.csv")
	if err != nil {
		log.Fatalf("unable to open CSV file, %v", err)
	}
	defer file.Close()

	// CSV ファイルを読み取る
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("unable to read CSV file, %v", err)
	}

	// CSV の1行目を使ってアクセスキーとシークレットキーを取得
	accessKeyID := records[1][0]     // 2行目の Access Key ID
	secretAccessKey := records[1][1] // 2行目の Secret Access Key

	// 環境変数に設定
	os.Setenv("AWS_ACCESS_KEY_ID", accessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", secretAccessKey)
	os.Setenv("AWS_REGION", "ap-northeast-1")

	// AWS 設定をロード
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// S3 クライアントを作成
	s3Client := s3.NewFromConfig(cfg)
	return s3Client
}

// アンマウントするコマンド
// umount /private/tmp/awsfs

// ディレクトリの内容をリストアップするメソッド
// lsコマンドを実行すると呼ばれる
// 全てのディレクトリに、"dir" という名前のディレクトリと、"file" という名前のファイルがある
func (cf *AwsFileSystem) Readdir(path string, fill func(name string, stat *fuse.Stat_t, offset int64) bool, offset int64, fh uint64) (errc int) {
	fmt.Println("Readdir", path, offset, fh)
	fill("dir", nil, 0)
	fill("file", nil, 0)
	return 0
}

// ファイルやディレクトリの属性を返すメソッド
// ls -l コマンドを実行すると呼ばれる
// catを実行した時にも呼ばれ、ディレクトリの場合はIs a directoryと言われる
func (cf *AwsFileSystem) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	fmt.Println("Getattr", path, stat, fh)
	// 末尾が "file" で終わる場合はファイル、それ以外はディレクトリとする
	// ファイルのサイズは100バイトで、最終更新日時は2000年1月1日とする
	if len(path) < 4 {
		stat.Mode = fuse.S_IFDIR | 0555
	} else if path[len(path)-4:] == "file" {
		stat.Mode = fuse.S_IFREG | 0555
		stat.Size = 100
		// 2000年
		stat.Mtim.Sec = 946684800
	} else {
		stat.Mode = fuse.S_IFDIR | 0555
	}
	return 0
}

// ファイルを開くメソッド
// ファイルハンドルは一律で0を返す
// catコマンドを実行すると呼ばれる
func (fs *AwsFileSystem) Open(path string, flags int) (n int, fh uint64) {
	fmt.Println("open", path, flags)
	return 0, 0
}

// ファイルの中身を読み込むメソッド
// 中身として、パス、オフセット、ファイルハンドルの情報を返す
// catコマンドを実行すると呼ばれる
func (cf *AwsFileSystem) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	fmt.Println("Read", path, ofst, fh)
	var result bytes.Buffer
	fmt.Fprintf(&result, "path=%s, ofst=%d, fh=%d", path, ofst, fh)
	copy(buff, result.Bytes())
	return int(len(buff))
}
