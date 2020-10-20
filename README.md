# goRDStest

### AWS RDS+Golang+GORMなテストリポジトリ

**Amazon Linux2 + AWS RDS(mysql)でgormを使って簡単なCRUDなAPIを実装するテスト**<br>
ユーザー名、パスワードを扱うシンプルなAPIです。<br>

### まず以下のようにRDS側でデータベース、テーブルを事前準備してください

```
CREATE DATABASE test DEFAULT CHARACTER SET utf8;

create table test.member(id int primary key, name varchar(8),password varchar(100));

ALTER TABLE test.member MODIFY id INT AUTO_INCREMENT;
```

### SSL証明書を準備してください。

[mkcert](https://kakakakakku.hatenablog.com/entry/2018/07/27/120009)などを使って**SSL証明書**を作成してください。

### 環境変数を定義してください。

OS側の環境変数を定義してください。```export API_USER=xxxx```みたいにして定義します。<br>

|変数名|定義内容|
|:---|:---|
|API_USER|RDSのユーザー名|
|API_PASS|RDSのパスワード|
|API_ADDRESS|RDSのDNS名|
|API_SALT|パスワード暗号化用文字列|

以下みたいに定義するといいでしょう。<br>

```
export API_USER=admin
export API_PASS=xxxx
export API_ADDRESS=xxxx.xxxx.us-east-2.rds.amazonaws.com
export API_SALT=api12345
```

### Amazon Linux(じゃなくてUbuntuとかでも良いけど)のインスタンスを起動します。

リポジトリをcloneするとか、zipを置くとかでコードを配置してください。<br>

```
go build api.go
```

で、go.modのモジュールが組み込まれて動くと思う。<br>
**ESCキーで終了**します。<br>

### 起動オプションがあります。

起動時に引数から与えるオプションがあります。<br>

|オプション名|定義内容|
|:---|:---|
|cert|SSL証明書の公開鍵ファイル|
|key|SSL証明書の秘密鍵ファイル|
|port|APIを起動するポート番号|
|debug|デバッグモード|

以下みたいに使います。<br>

```
./api -cert=test.pem -key=test-key.pem -debug
```

### 使えるAPI

以下APIがあります。<br>

|API名|JSON名|機能|
|:---|:---|:---|
|raw|raw|生SQLを投げ込みます|
|find|search|ユーザー名検索します|
|create|name / password|アカウントを作成します|
|read|id|該当ID番号のユーザーを表示します|
|update|id / name / password|アカウント情報をアップデートします|
|delete|id|該当ID番号のユーザーを消します|
|auth|name / password|ユーザー名、パスワードで認証します|

以下みたいに使います。<br>

```
curl -k -H "Content-Type: application/json" -X POST -d '{"raw":" SELECT * FROM test.member;"}' https://localhost:8080/raw
curl -k -H "Content-Type: application/json" -X POST -d '{"search":"user3"}' https://localhost:8080/find
curl -k -H "Content-Type: application/json" -X POST -d '{"name":"user2", "password": "pass"}' https://localhost:8080/create
curl -k -H "Content-Type: application/json" -X POST -d '{"id":"8"}' https://localhost:8080/read
curl -k -H "Content-Type: application/json" -X POST -d '{"id": "8", "name":"user2", "password": "pass"}' https://localhost:8080/update
curl -k -H "Content-Type: application/json" -X POST -d '{"id":1}' https://localhost:8080/delete
curl -k -H "Content-Type: application/json" -X POST -d '{"name":"user2", "password": "pass"}' https://localhost:8080/auth
```

### テスト用コード

~~既存のDBべったりで書いてしまったのでそのうち修正しないと。。~~

### 開発するのに超お役立ちなFYI

[Visual Studio Code で編集中のテストコードを実行する (Golang編)](https://qiita.com/ykato/items/6b50d7d14be05128f74d)
[VS CodeのGo言語テストコード生成ツールを使ってみたらめちゃくちゃ便利だった話とか](https://kdnakt.hatenablog.com/entry/2019/01/03/080000)


