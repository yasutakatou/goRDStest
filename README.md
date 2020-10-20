# goRDStest

### AWS RDS+Golang+GORMなテストリポジトリ

Amazon Linux2 + AWS RDS(mysql)でgormを使ってCRUDなAPIを実装するテスト<br>

### まず以下のようにRDS側でデータベース、テーブルを事前準備してください

```
CREATE DATABASE test DEFAULT CHARACTER SET utf8;

create table test.member(id int primary key, name varchar(8),password varchar(100));

ALTER TABLE test.member MODIFY id INT AUTO_INCREMENT;
```

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

*note) SALTはパスワードをRDSに保存時に暗号化する際に使う文字列です。8～16桁で指定できます。*

### Amazon Linux(じゃなくてUbuntuとかでも良いけど)のインスタンスを起動します。

リポジトリをcloneするとか、zipを置くとかでコードを配置してください。<br>

```
go build api.go
```

で、go.modのモジュールが組み込まれて動くと思う。<br>

### 
