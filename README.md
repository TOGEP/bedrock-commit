## bedrock-commits

## Requirement
Go1.20 or later

## Usage

- 開発がある程度すすむまで各自ビルド
```
$ go build
$ ./bedrock-commits
```


## Setup
### AWS credentials

本ツールの利用にあたり、AWS Bedrockが利用可能なAWSアカウントの認証情報が必要となります。  

認証情報の発行については、[AWSの公式ドキュメント](https://docs.aws.amazon.com/ja_jp/powershell/latest/userguide/pstools-appendix-sign-up.html)を参照してください。  

発行した認証情報はAWS設定ファイル（通常は~/.aws/configまたは~/.aws/credentials）に格納するか、プロジェクトルートに`.env`ファイルを作成し、以下のように環境変数として設定してください。

```env
AWS_ACCESS_KEY_ID=《アクセスキーを入力》
AWS_SECRET_ACCESS_KEY=《シークレットアクセスキーを入力》
AWS_REGION=《リージョンを入力》
```

### AWS Bedrockの有効化

[Model Access](https://ap-northeast-1.console.aws.amazon.com/bedrock/home?region=ap-northeast-1#/modelaccess)の画面に行き、希望のモデル利用を申請する
