# readme
superlicense

原理: 数字证书

## license组织形式
```bash
$ ls license
<license-id>-<category>.dat // core license
<license-parent>-<category>.dat // others license
```

## design
1. 使用yaml作为license format

    可读性强
1. license包含parent_id属性

    多个license可组合使用或分批次导入, 尽量满足更多场景

## example
```bash
cd cmd/superlicense
./build.sh
./superlicense generate --i license.yaml
./superlicense parse --i license.dat
```

## todo
- [ ] 使用更安全的RSA OEAP替代PKCS1.5