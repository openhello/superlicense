# README
原因: 证书签名

## 生成rsa
```bash
openssl genrsa -out license.key 4096
openssl rsa -in license.key -pubout -outform PEM -out license.key.pub
```

## license组织形式
```bash
$ ls license
<license-id>-<category>.dat // core license
<license-parent>-<category>.dat // others license
```