#!/bin/bash
awslocal s3 mb s3://tavinikkiy-local

cd /home/localstack/data/
awslocal s3 cp LP-5.png s3://tavinikkiy-local/
awslocal s3 ls s3://tavinikkiy-local

# バケットのパブリックアクセスブロックを無効化
awslocal s3api put-public-access-block \
    --bucket tavinikkiy-local \
    --public-access-block-configuration "BlockPublicAcls=false,IgnorePublicAcls=false,BlockPublicPolicy=false,RestrictPublicBuckets=false"

# バケットポリシーを設定
awslocal s3api put-bucket-policy \
    --bucket tavinikkiy-local \
    --policy '{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::tavinikkiy-local/*"
        }
    ]
}'

# CORSの設定
awslocal s3api put-bucket-cors --bucket tavinikkiy-local --cors-configuration '{
    "CORSRules": [
        {
            "AllowedHeaders": ["*"],
            "AllowedMethods": ["GET", "PUT", "POST", "DELETE"],
            "AllowedOrigins": ["*"],
            "ExposeHeaders": []
        }
    ]
}'
