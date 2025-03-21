#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查是否安装了Docker和Docker Compose
check_dependencies() {
  echo -e "${BLUE}检查依赖...${NC}"
  
  if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: 未找到Docker命令。请先安装Docker。${NC}"
    echo "安装指南: https://docs.docker.com/get-docker/"
    exit 1
  fi
  
  if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}错误: 未找到Docker Compose命令。请先安装Docker Compose。${NC}"
    echo "安装指南: https://docs.docker.com/compose/install/"
    exit 1
  fi
  
  echo -e "${GREEN}依赖检查通过。${NC}"
}

# 检查.env文件
check_env_file() {
  if [ ! -f .env ]; then
    echo -e "${YELLOW}未找到.env文件，将使用示例配置创建。${NC}"
    if [ -f .env.example ]; then
      cp .env.example .env
      echo -e "${GREEN}已创建.env文件，请根据需要修改配置。${NC}"
    else
      echo -e "${RED}错误: 未找到.env.example文件。${NC}"
      exit 1
    fi
  fi
}

# 创建必要的目录
create_directories() {
  echo -e "${BLUE}创建必要的目录...${NC}"
  
  # 确保nginx配置目录存在
  mkdir -p nginx/conf.d
  mkdir -p nginx/ssl
  mkdir -p nginx/www
  
  # 确保数据库初始化脚本目录存在
  mkdir -p init-scripts/mysql
  
  echo -e "${GREEN}目录创建完成。${NC}"
}

# 检查Nginx配置
check_nginx_config() {
  if [ ! -f nginx/conf.d/default.conf ]; then
    echo -e "${YELLOW}未找到Nginx配置文件，将使用默认配置。${NC}"
    if [ -f nginx/conf.d/default.conf.example ]; then
      cp nginx/conf.d/default.conf.example nginx/conf.d/default.conf
    else
      echo -e "${RED}错误: 未找到默认的Nginx配置示例文件。${NC}"
      exit 1
    fi
  fi
}

# 部署服务
deploy() {
  echo -e "${BLUE}开始部署服务...${NC}"
  
  # 拉取最新镜像
  echo -e "${YELLOW}拉取最新镜像...${NC}"
  docker-compose pull
  
  # 构建和启动服务
  echo -e "${YELLOW}构建和启动服务...${NC}"
  docker-compose up -d --build
  
  echo -e "${GREEN}部署完成！${NC}"
  echo -e "可以通过以下命令查看日志: docker-compose logs -f"
}

# 停止服务
stop() {
  echo -e "${BLUE}停止服务...${NC}"
  docker-compose down
  echo -e "${GREEN}服务已停止。${NC}"
}

# 重启服务
restart() {
  echo -e "${BLUE}重启服务...${NC}"
  docker-compose restart
  echo -e "${GREEN}服务已重启。${NC}"
}

# 显示服务状态
status() {
  echo -e "${BLUE}服务状态:${NC}"
  docker-compose ps
}

# 显示服务日志
logs() {
  echo -e "${BLUE}显示服务日志:${NC}"
  docker-compose logs -f $1
}

# 清理所有数据（危险操作）
clean() {
  echo -e "${RED}警告: 此操作将删除所有容器和卷数据！${NC}"
  read -p "是否继续? (y/n) " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}停止并移除所有容器...${NC}"
    docker-compose down -v
    echo -e "${YELLOW}清理未使用的镜像...${NC}"
    docker image prune -a -f
    echo -e "${GREEN}清理完成。${NC}"
  fi
}

# 备份数据库
backup_db() {
  local backup_dir="backups"
  local date_str=$(date +"%Y%m%d_%H%M%S")
  local backup_file="${backup_dir}/mysql_backup_${date_str}.sql"
  
  mkdir -p $backup_dir
  
  echo -e "${BLUE}开始备份数据库...${NC}"
  
  # 从.env文件加载数据库信息
  source .env
  
  docker-compose exec -T mysql mysqldump -u root -p"${MYSQL_ROOT_PASSWORD}" --all-databases > $backup_file
  
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}数据库备份成功: ${backup_file}${NC}"
  else
    echo -e "${RED}数据库备份失败${NC}"
    exit 1
  fi
}

# 恢复数据库
restore_db() {
  if [ -z "$1" ]; then
    echo -e "${RED}错误: 请指定要恢复的备份文件。${NC}"
    echo "用法: $0 restore_db backups/mysql_backup_file.sql"
    exit 1
  fi
  
  local backup_file=$1
  
  if [ ! -f $backup_file ]; then
    echo -e "${RED}错误: 备份文件不存在: ${backup_file}${NC}"
    exit 1
  fi
  
  echo -e "${BLUE}开始恢复数据库...${NC}"
  
  # 从.env文件加载数据库信息
  source .env
  
  docker-compose exec -T mysql mysql -u root -p"${MYSQL_ROOT_PASSWORD}" < $backup_file
  
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}数据库恢复成功${NC}"
  else
    echo -e "${RED}数据库恢复失败${NC}"
    exit 1
  fi
}

# 显示帮助信息
show_help() {
  echo -e "${BLUE}KY-Admin 部署管理脚本${NC}"
  echo ""
  echo "用法: $0 [命令]"
  echo ""
  echo "命令:"
  echo "  deploy     - 部署服务"
  echo "  stop       - 停止服务"
  echo "  restart    - 重启服务"
  echo "  status     - 显示服务状态"
  echo "  logs       - 显示服务日志"
  echo "  backup_db  - 备份数据库"
  echo "  restore_db - 恢复数据库"
  echo "  clean      - 清理所有数据（危险操作）"
  echo "  help       - 显示此帮助信息"
  echo ""
  echo "示例:"
  echo "  $0 deploy"
  echo "  $0 logs ky-admin"
  echo "  $0 backup_db"
  echo "  $0 restore_db backups/mysql_backup_20250322_123000.sql"
}

# 主函数
main() {
  # 如果没有参数，显示帮助信息
  if [ $# -eq 0 ]; then
    show_help
    exit 0
  fi
  
  # 解析命令
  case "$1" in
    deploy)
      check_dependencies
      check_env_file
      create_directories
      check_nginx_config
      deploy
      ;;
    stop)
      stop
      ;;
    restart)
      restart
      ;;
    status)
      status
      ;;
    logs)
      logs $2
      ;;
    backup_db)
      backup_db
      ;;
    restore_db)
      restore_db $2
      ;;
    clean)
      clean
      ;;
    help|*)
      show_help
      ;;
  esac
}

# 执行主函数
main "$@" 