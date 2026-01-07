GRANT ALL PRIVILEGES ON test_tavinikkiy.* TO 'user'@'%';
GRANT ALL PRIVILEGES ON develop_tavinikkiy.* TO 'user'@'%';
FLUSH PRIVILEGES;

CREATE DATABASE IF NOT EXISTS test_tavinikkiy;
CREATE DATABASE IF NOT EXISTS develop_tavinikkiy;
