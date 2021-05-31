# readme
superlicense

原理: 数字证书

## design
1. 使用yaml作为license format

    可读性强
1. license包含parent_id属性

    多个license可组合使用或分批导入， 尽量满足更多场景

## todo
- [ ] 使用更安全的RSA OEAP替代PKCS1.5