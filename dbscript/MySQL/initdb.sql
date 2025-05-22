/* SQLEditor (MySQL (2))*/


DROP TABLE IF EXISTS s_departments;

DROP TABLE IF EXISTS s_depusers;

DROP TABLE IF EXISTS s_group;

DROP TABLE IF EXISTS s_groupuser;

DROP TABLE IF EXISTS s_items;

DROP TABLE IF EXISTS s_logs;

DROP TABLE IF EXISTS s_operators;

DROP TABLE IF EXISTS s_resources;

DROP TABLE IF EXISTS s_role;

DROP TABLE IF EXISTS s_rolegroup;

DROP TABLE IF EXISTS s_roleoperator;

DROP TABLE IF EXISTS s_users;

CREATE TABLE s_departments
(
id INTEGER AUTO_INCREMENT COMMENT '编号',
name VARCHAR(50) UNIQUE  COMMENT '组织名称',
createtime TIMESTAMP COMMENT '创建时间',
icon VARBINARY(100),
PRIMARY KEY (id)
) COMMENT='组织表';

CREATE TABLE s_depusers
(
userid VARCHAR(50) COMMENT '用户编号',
departmentid INTEGER COMMENT '组织编号',
PRIMARY KEY (userid,departmentid)
) COMMENT='组织成员表';

CREATE TABLE s_group
(
id INTEGER NOT NULL AUTO_INCREMENT,
name VARCHAR(20) COMMENT '组名',
createtime TIMESTAMP COMMENT '创建时间',
PRIMARY KEY (id)
) COMMENT='用户组';

CREATE TABLE s_groupuser
(
groupid INTEGER COMMENT '组编号',
userid VARCHAR(50) COMMENT '用户编号',
PRIMARY KEY (groupid,userid)
);

CREATE TABLE s_items
(
id INTEGER AUTO_INCREMENT,
name VARCHAR(100),
type INTEGER COMMENT '1、names',
PRIMARY KEY (id)
);

CREATE TABLE s_logs
(
id INTEGER AUTO_INCREMENT,
account VARCHAR(50),
ip VARCHAR(50),
logintime TIMESTAMP,
PRIMARY KEY (id)
);

CREATE TABLE s_operators
(
id INTEGER UNIQUE  COMMENT '功能编号',
name VARCHAR(20) COMMENT '功能名称',
PRIMARY KEY (id)
) COMMENT='操作功能表';

CREATE TABLE s_resources
(
resid INTEGER AUTO_INCREMENT COMMENT '资源编号',
userid VARCHAR(50) COMMENT '拥有者编号',
`read` BOOLEAN COMMENT '用户读权限',
`write` BOOLEAN COMMENT '用户写权限',
download BOOLEAN COMMENT '用户下载权限',
`delete` BOOLEAN COMMENT '删除权限',
PRIMARY KEY (resid,userid)
) COMMENT='资源权限表

每一个文档有多个这个对应';

CREATE TABLE s_role
(
id INTEGER NOT NULL AUTO_INCREMENT COMMENT '角色编号',
name VARCHAR(20) UNIQUE  COMMENT '角色名称',
createtime TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
PRIMARY KEY (id)
) COMMENT='角色';

CREATE TABLE s_rolegroup
(
roleid INTEGER COMMENT '角色编号',
groupid INTEGER COMMENT '组编号',
PRIMARY KEY (roleid,groupid)
) COMMENT='角色用户组';

CREATE TABLE s_roleoperator
(
roleid INTEGER COMMENT '角色编号',
operatorid INTEGER COMMENT '操作功能的编号',
acte BOOLEAN COMMENT '是否可以操作',
PRIMARY KEY (roleid,operatorid)
) COMMENT='角色功能操作表';

CREATE TABLE s_users
(
id VARCHAR(50),
email VARCHAR(50) UNIQUE ,
password VARCHAR(80),
name VARCHAR(50),
icon VARCHAR(100) COMMENT '头像的url',
state INTEGER COMMENT '状态：1，可用，0禁用,2自动注册',
createtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
modtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
sex INTEGER COMMENT '性别：1男，0女，2未知',
mobile VARCHAR(20),
aliasname VARCHAR(50),
PRIMARY KEY (id)
);
CREATE TABLE s_urlmappings
(
id INTEGER AUTO_INCREMENT,
operatorid INTEGER,
url VARCHAR(255),
PRIMARY KEY (id)
);
