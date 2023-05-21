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

-- 创建 Program 表
CREATE TABLE IF NOT EXISTS program (
    id INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(255) NOT NULL,
    exec_user VARCHAR(255),
    name VARCHAR(255),
    description TEXT,
    team_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建 versions 表
CREATE TABLE IF NOT EXISTS versions (
    uuid VARCHAR(255) PRIMARY KEY,
    program_uuid VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    release_note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (program_uuid) REFERENCES program (uuid) ON DELETE CASCADE ON UPDATE CASCADE
);

-- 创建 packages 表
CREATE TABLE IF NOT EXISTS packages (
    id INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    version_uuid VARCHAR(255) NOT NULL,
    os VARCHAR(255),
    arch VARCHAR(255),
    storage_path VARCHAR(255),
    download_path VARCHAR(255),
    md5 VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (version_uuid) REFERENCES versions (uuid) ON DELETE CASCADE ON UPDATE CASCADE
);

-- 创建 program_actions 表
CREATE TABLE IF NOT EXISTS program_actions (
    id INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    action_uuid VARCHAR(255) NOT NULL,
    program_uuid VARCHAR(255) NOT NULL,
    action_type VARCHAR(255),
    content TEXT,
    status VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (program_uuid) REFERENCES programs (uuid) ON DELETE CASCADE ON UPDATE CASCADE
);

-- 创建 program_action_templates 表
CREATE TABLE IF NOT EXISTS program_action_templates (
    id INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    template_name VARCHAR(255),
    status VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建 template_actions 表
CREATE TABLE IF NOT EXISTS template_actions (
    id INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    template_action_uuid VARCHAR(255) NOT NULL,
    action_uuid VARCHAR(255) NOT NULL,
    sequence INT(11) NOT NULL,
    FOREIGN KEY (template_action_uuid) REFERENCES program_actions (action_uuid) ON DELETE CASCADE ON UPDATE CASCADE
);
