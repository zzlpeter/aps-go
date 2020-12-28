# aps-go

项目结构说明
```
conf                        配置文件
    local.toml
    test.toml
    product.toml
docs                        swagger所使用
grequests                   三方HTTP请求
libs                        库函数
    log                         日志模块
    mysql                       数据库模块
    redis                       Redis模块
    tomlc                       解析toml文件模块
    utils                       其他库函数
        datekits                    时间相关
        type_convert                类型转换相关
        utils                       其他
models                      数据库models
Dockerfile
main.go                     入口函数
README.md
```

启动流程
```
修改conf/conf.toml配置文件
执行SQL创建表
CREATE TABLE `task`.`task`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键自增',
  `task_key` varchar(63) NOT NULL COMMENT '任务唯一标识',
  `execute_func` varchar(255) NOT NULL COMMENT '任务执行的方法',
  `spec` varchar(255) CHARACTER NOT NULL COMMENT '调度时间',
  `params` varchar(255) CHARACTER NOT NULL DEFAULT '' COMMENT '执行方法的参数',
  `is_valid` int(11) NOT NULL COMMENT '是否有效',
  `status` varchar(31) NOT NULL COMMENT '执行状态(ready/doing)',
  `extra` varchar(255) NOT NULL DEFAULT '{}' COMMENT '额外信息',
  `create_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务描述',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `ux__task__task_key_status`(`task_key`, `status`) USING BTREE COMMENT 'task_key/status唯一索引'
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '任务表' ROW_FORMAT = Compact;

CREATE TABLE `task`.`task_execute`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键自增',
  `task_id` int(11) NOT NULL COMMENT '任务ID',
  `status` varchar(255) NOT NULL COMMENT '执行状态',
  `extra` text COMMENT '额外信息',
  `create_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `trace_id` varchar(255) DEFAULT '' COMMENT 'trace id',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx__execute_task__task_id_create_at`(`task_id`, `create_at`) USING BTREE COMMENT '任务ID、创建时间'
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '任务执行表' ROW_FORMAT = Compact;

启动生产者
go run main.go -action=producer

启动消费者
go run main.go -action=consumer
```