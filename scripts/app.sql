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
CREATE TABLE IF NOT EXISTS program (
    id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    uuid VARCHAR(255) NOT NULL UNIQUE,
    exec_user VARCHAR(255) DEFAULT NULL,
    name VARCHAR(255) DEFAULT NULL,
    description TEXT,
    team_id VARCHAR(255) DEFAULT NULL,
    install_path VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
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


CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task_id VARCHAR(255) NOT NULL UNIQUE,
    task_name VARCHAR(255) NOT NULL,
    task_type VARCHAR(255) NOT NULL,
    task_status VARCHAR(255) NOT NULL,
    parent_task_id VARCHAR(255),
    content TEXT,
    description TEXT,
    creater VARCHAR(255) NOT NULL,
    team_id VARCHAR(255) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL
);

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
    timeout INTERVAL,
    parent_record_id VARCHAR(255)
);

