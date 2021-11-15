# h2 database read only

## h2 数据库 文件权限问题导致无法写入

> Causedby: org.h2.jdbc.JdbcSQLNonTransientException:The database is read only; SQL statement: UPDATE PUBLIC.DATABASECHANGELOGLOCK SET LOCKED = TRUE, LOCKEDBY =‘10.16.0.5 (10.16.0.5)’, LOCKGRANTED =‘2020-06-17 15:07:20.707’ WHERE ID =1 AND LOCKED = FALSE [90097-200] at org.h2.message.DbException.getJdbcSQLException(DbException.java:505) at org.h2.message.DbException.getJdbcSQLException(DbException.java:429) at org.h2.message.DbException.get(DbException.java:205)

### emdbedded模式下java访问h2数据库文件需要获取权限