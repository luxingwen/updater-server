# updater-server



openapi-generator 生成python代码

[openapi-generator](https://github.com/OpenAPITools/openapi-generator#1---installation)

> openapi-generator generate -i docs/swagger.json -g python -o scripts/test --skip-validate-spec

mysql导出表结构数据

mysqldump -u root -p updater --no-data > scripts/updater.sql