CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY, -- 主键
    vmuuid VARCHAR(255) UNIQUE COMMENT '虚拟机UUID', -- 唯一索引，表示虚拟机UUID
    sn VARCHAR(255) COMMENT '序列号', -- 序列号
    hostname VARCHAR(255) COMMENT '主机名', -- 主机名
    ip VARCHAR(255) COMMENT 'IP地址', -- IP地址
    proxy_id VARCHAR(255) COMMENT '代理ID', -- 代理ID
    status VARCHAR(255) COMMENT '状态', -- 状态
    os VARCHAR(255) COMMENT '操作系统', -- 操作系统
    arch VARCHAR(255) COMMENT '架构', -- 架构
    created TIMESTAMP COMMENT '创建时间', -- 创建时间
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间' -- 更新时间
);


CREATE INDEX IF NOT EXISTS idx_sn ON clients (sn);
CREATE INDEX IF NOT EXISTS idx_hostname ON clients (hostname);
CREATE INDEX IF NOT EXISTS idx_ip ON clients (ip);
CREATE INDEX IF NOT EXISTS idx_proxy_id ON clients (proxy_id); -- 为proxy_id字段创建索引



-- 创建 program 表
CREATE TABLE `program` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) NOT NULL,
  `exec_user` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `description` text,
  `team_id` varchar(255) DEFAULT NULL,
  `windows_install_path` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `linux_install_path` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB;

-- 创建 versions 表
CREATE TABLE IF NOT EXISTS versions (
    id INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(255) UNIQUE NOT NULL,
    program_uuid VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    release_note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (program_uuid) REFERENCES program (uuid) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- 创建 packages 表
CREATE TABLE IF NOT EXISTS packages (
    id INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(255) UNIQUE NOT NULL,
    version_uuid VARCHAR(255) NOT NULL,
    os VARCHAR(255),
    arch VARCHAR(255),
    storage_path VARCHAR(255),
    download_path VARCHAR(255),
    md5 VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (version_uuid) REFERENCES versions (uuid) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;


-- 创建 program_action 表
CREATE TABLE IF NOT EXISTS program_action (
    id INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(255) UNIQUE NOT NULL,
    program_uuid VARCHAR(255) NOT NULL,
    action_type ENUM('Download', 'Install', 'Start', 'Stop', 'Uninstall', 'Backup', 'Status', 'Version', 'Single', 'Composite'),
    name VARCHAR(255),
    content TEXT,
    status VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (program_uuid) REFERENCES program (uuid) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;




CREATE TABLE `tasks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `task_id` varchar(255) NOT NULL,
  `task_name` varchar(255) DEFAULT NULL,
  `task_type` varchar(255) DEFAULT NULL,
  `task_status` varchar(255) DEFAULT NULL,
  `parent_task_id` varchar(255) DEFAULT NULL,
  `content` longtext DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `creater` varchar(255) DEFAULT NULL,
  `team_id` varchar(255) DEFAULT NULL,
  `category` varchar(255) DEFAULT NULL,
  `next_task_id` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `ext` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `task_id` (`task_id`),
  INDEX `idx_parent_task_id` (`parent_task_id`)
) ENGINE=InnoDB;




CREATE TABLE IF NOT EXISTS task_execution_records (
    id SERIAL PRIMARY KEY,
    record_id VARCHAR(255) NOT NULL UNIQUE,
    task_id VARCHAR(255) NOT NULL,
    client_uuid VARCHAR(255) NOT NULL,
    task_type VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    stdout TEXT,
    stderr TEXT,
    message TEXT,
    script_exit_code INTEGER,
    code VARCHAR(255),
    content TEXT,
    timeout INT(11),
    parent_record_id VARCHAR(255),
    next_record_id VARCHAR(255),
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE=InnoDB;



CREATE TABLE users (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  uuid VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  phone VARCHAR(255) NOT NULL,
  role VARCHAR(255) NOT NULL,
  teamId VARCHAR(255) NOT NULL,
  avatar VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);



ALTER TABLE tasks
ADD COLUMN trace_id varchar(255) DEFAULT NULL,
ADD INDEX idx_trace_id (trace_id);


ALTER TABLE task_execution_records
ADD COLUMN trace_id varchar(70) DEFAULT NULL,
ADD INDEX idx_trace_id (trace_id);