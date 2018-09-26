CREATE TABLE `cron_service` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `task_name` varchar(30) NOT NULL COMMENT '任务名',
    `time_format` varchar(30) NOT NULL COMMENT '时间格式',
    `service_url` varchar(100) NOT NULL COMMENT '路由',
    `service_method` varchar(50) NOT NULL COMMENT '被调用方方法名',
    `method` varchar(10) NOT NULL COMMENT 'http方法名',
    `service_header` text COMMENT 'header头',
    `form` text COMMENT 'form',
    `service_body` text COMMENT 'body',
    `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '当前定时任务是否可用：0-不可用; 1-可用',
    `lock_key` varchar(40) unique NOT NULL COMMENT 'redis lock key',
    `task_desc` text COMMENT '任务描述',
    `create_time` datetime NOT NULL COMMENT '创建时间',
    `update_time` datetime COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name_format_url` (`task_name`, `time_format`, `service_url`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='存储task的定时任务相关信息';