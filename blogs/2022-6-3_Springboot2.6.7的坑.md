> 有些依赖，不要选择最新的版本

在使用 JPA 的简单查询 findAllByNameContainsIgnoreCase 时，第一次正常，之后就会一直报错。

说实话，我是真没想到这是一个 hibernate 的 bug。

[https://github.com/spring-projects/spring-data-jpa/issues/2472](https://github.com/spring-projects/spring-data-jpa/issues/2472)

在 hibernate 版本 5.6.6 ~ 5.6.8 中用 startwith 以及 like 查询都会出现这个问题，该 bug 已在 hibernate 版本 5.6.9 中得到修复([HHH-15142](https://hibernate.atlassian.net/browse/HHH-15142) - Hibernate JIRA)

最后，将 Springboot 版本从 2.6.7 降到 2.6.3 后就不再报错。

```
java.lang.IllegalArgumentException: Parameter value [\] did not match expected type [java.lang.String (n/a)]
	at org.hibernate.query.spi.QueryParameterBindingValidator.validate(QueryParameterBindingValidator.java:54) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.query.spi.QueryParameterBindingValidator.validate(QueryParameterBindingValidator.java:27) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.query.internal.QueryParameterBindingImpl.validate(QueryParameterBindingImpl.java:90) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.query.internal.QueryParameterBindingImpl.setBindValue(QueryParameterBindingImpl.java:55) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.query.internal.AbstractProducedQuery.setParameter(AbstractProducedQuery.java:501) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.query.internal.AbstractProducedQuery.setParameter(AbstractProducedQuery.java:122) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.query.criteria.internal.compile.CriteriaCompiler$1$1.bind(CriteriaCompiler.java:141) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.query.criteria.internal.CriteriaQueryImpl$1.buildCompiledQuery(CriteriaQueryImpl.java:364) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.query.criteria.internal.compile.CriteriaCompiler.compile(CriteriaCompiler.java:171) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.internal.AbstractSharedSessionContract.createQuery(AbstractSharedSessionContract.java:774) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
	at org.hibernate.internal.AbstractSharedSessionContract.createQuery(AbstractSharedSessionContract.java:114) ~[hibernate-core-5.6.8.Final.jar:5.6.8.Final]
    ...
```