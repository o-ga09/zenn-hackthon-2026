package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	Cfg "github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/image"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ulid"
)

type CloudflareR2Storage struct {
	client     *s3.Client
	bucketName string
	endpoint   string
}

func NewCloudflareR2Storage(ctx context.Context, accountID, accessKeyID, accessKeySecret, bucketName string) (*CloudflareR2Storage, error) {
	var endpoint string
	ctx, err := Cfg.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}
	envCfg := Cfg.GetCtxEnv(ctx)

	if envCfg.Env == "local" || envCfg.Env == "test" {
		endpoint = "http://localstack:4566"
	} else {
		endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           endpoint,
			SigningRegion: "auto",
			Source:        aws.EndpointSourceCustom,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
		config.WithDefaultRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &CloudflareR2Storage{
		client:     client,
		bucketName: bucketName,
		endpoint:   endpoint,
	}, nil
}

func (s *CloudflareR2Storage) Upload(ctx context.Context, key string, base64Data string) (string, error) {
	var err error

	// Data URLプレフィックスを除去（例: "data:image/jpeg;base64," を除去）
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	// Base64デコード
	fileData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		fmt.Println(err)
		return "", errors.ErrFailedDecodeImage
	}

	// ファイル名を生成
	fileID, err := ulid.GenerateULID()
	if err != nil {
		return "", errors.ErrInvalidULID
	}

	// ファイルの種類を判定
	contentType := image.DetectContentType(fileData)
	if !image.IsValidImageType(contentType) {
		return "", errors.ErrFailedImageName
	}

	ext := image.GetExtensionFromContentType(contentType)
	path := fmt.Sprintf("%s%s%s", key, fileID, ext)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(path),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String(contentType),
	}

	_, err = s.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return path, nil
}

// UploadFile はmultipart/form-dataから受け取ったファイルデータを直接アップロードする
func (s *CloudflareR2Storage) UploadFile(ctx context.Context, key string, fileData []byte, contentType string) (string, error) {
	// ファイル名をkeyとして使用（すでにパス込みで渡される）
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String(contentType),
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return key, nil
}

func (s *CloudflareR2Storage) Delete(ctx context.Context, key string) error {
	var err error
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	_, err = s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *CloudflareR2Storage) Get(ctx context.Context, key string) (string, error) {
	var (
		err    error
		result *s3.GetObjectOutput
	)

	if key == "" {
		return "", errors.ErrNotFoundImage
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	result, err = s.client.GetObject(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "error StatusCode: 404") {
			return "", errors.ErrNotFoundImage
		}
		return "", fmt.Errorf("failed to get file: %w", err)
	}

	// ファイルの内容を読み取ってBase64エンコード
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, result.Body); err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	// Base64エンコード文字列に変換
	base64Data := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(base64Data, buf.Bytes())
	base64Str := string(base64Data)

	// データ部を追加する
	contentType := aws.ToString(result.ContentType)
	base64Str = fmt.Sprintf("data:%s;base64,%s", contentType, base64Str)

	// ファイルの内容を文字列として返す
	return base64Str, nil
}

func (s *CloudflareR2Storage) List(ctx context.Context, prefix string) (map[string]string, error) {
	var (
		err    error
		result *s3.ListObjectsV2Output
	)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(prefix),
	}

	result, err = s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	base64ImagesMap := make(map[string]string)
	for _, obj := range result.Contents {
		// 各オブジェクトの内容を取得
		getInput := &s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    obj.Key,
		}

		getResult, err := s.client.GetObject(ctx, getInput)
		if err != nil {
			// 個別のファイル取得でエラーが発生した場合はスキップ
			continue
		}

		// ファイルの内容を読み取ってBase64エンコード
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, getResult.Body); err != nil {
			getResult.Body.Close()
			continue
		}
		getResult.Body.Close()

		// Base64エンコード文字列に変換
		base64Data := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
		base64.StdEncoding.Encode(base64Data, buf.Bytes())
		base64Str := string(base64Data)

		// データ部を追加する
		contentType := aws.ToString(getResult.ContentType)
		if contentType == "" {
			// ContentTypeが取得できない場合は、ファイルの内容から判定
			contentType = image.DetectContentType(buf.Bytes())
		}
		base64Str = fmt.Sprintf("data:%s;base64,%s", contentType, base64Str)

		// ファイル名をキーとしてマップに追加
		fileName := aws.ToString(obj.Key)
		base64ImagesMap[fileName] = base64Str
	}

	return base64ImagesMap, nil
}

func ObjectURKFromKey(endpoint, key string) string {
	return fmt.Sprintf("%s/%s", endpoint, key)
}
