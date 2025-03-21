#!/usr/bin/env node

/**
 * 该脚本用于将Swagger文档同步至YAPI平台
 * 
 * 使用方法:
 * 1. 确保已安装axios: 通过swagger-update.sh脚本安装
 * 2. 设置环境变量或直接修改脚本中的配置
 * 3. 运行: NODE_PATH=scripts/node/node_modules node scripts/sync-swagger-to-yapi.js
 */

const fs = require('fs');
const path = require('path');
let axios;

// 尝试导入axios，考虑不同的安装位置
try {
  axios = require('axios');
} catch (err) {
  console.error('错误: 无法加载axios模块。请先运行swagger-update.sh安装依赖。');
  process.exit(1);
}

// 获取脚本所在目录
const scriptDir = path.dirname(__filename);
const projectDir = path.resolve(scriptDir, '..');

// 配置信息
const config = {
  // YAPI服务器地址
  yapiBaseUrl: process.env.YAPI_URL || 'http://localhost:3000',
  
  // YAPI项目的token，可在项目设置中获取
  yapiToken: process.env.YAPI_TOKEN || '932165370c9bb310570efda1ce1606bd40220b158d31c1db6efa553f410fb2c7',
  
  // Swagger文档路径 - 优先使用环境变量，否则默认为ky-admin/swagger/swagger.json
  swaggerPath: process.env.SWAGGER_PATH || path.join(projectDir, 'ky-admin/swagger/swagger.json'),
  
  // YAPI导入接口
  importUrl: '/api/open/import_data',

  // 测试模式（不实际发送请求）
  testMode: process.env.TEST_MODE === 'true' || false
};

// 控制台颜色
const colors = {
  green: '\x1b[32m',
  red: '\x1b[31m',
  yellow: '\x1b[33m',
  reset: '\x1b[0m'
};

/**
 * 打印彩色日志
 */
function log(message, color = colors.reset) {
  console.log(`${color}${message}${colors.reset}`);
}

/**
 * 检查配置有效性
 */
function validateConfig() {
  if (!config.yapiToken && !config.testMode) {
    log('错误: 未设置YAPI项目token，请设置环境变量YAPI_TOKEN或直接修改脚本中的yapiToken', colors.red);
    process.exit(1);
  } else if (!config.yapiToken && config.testMode) {
    log('警告: 未设置YAPI项目token，但处于测试模式，将继续执行', colors.yellow);
  }
  
  if (!fs.existsSync(config.swaggerPath)) {
    log(`错误: Swagger文件不存在: ${config.swaggerPath}`, colors.red);
    process.exit(1);
  }
}

/**
 * 读取Swagger文档
 */
function readSwaggerFile() {
  try {
    const content = fs.readFileSync(config.swaggerPath, 'utf8');
    return JSON.parse(content);
  } catch (error) {
    log(`错误: 读取Swagger文件失败: ${error.message}`, colors.red);
    process.exit(1);
  }
}

/**
 * 同步Swagger到YAPI
 */
async function syncToYapi(swaggerData) {
  try {
    log('开始同步Swagger文档至YAPI...', colors.yellow);
    
    if (config.testMode) {
      log('测试模式: 不实际发送请求至YAPI服务器', colors.yellow);
      log('将要发送的数据:', colors.yellow);
      log(`- YAPI地址: ${config.yapiBaseUrl}${config.importUrl}`);
      log(`- Token: ${config.yapiToken.substring(0, 3)}...`);
      log(`- Swagger数据: 包含 ${Object.keys(swaggerData.paths || {}).length} 个接口`);
      
      // 打印接口列表
      log('\n接口列表:', colors.green);
      for (const path in swaggerData.paths) {
        for (const method in swaggerData.paths[path]) {
          const endpoint = swaggerData.paths[path][method];
          log(`${method.toUpperCase().padEnd(6)} ${path.padEnd(30)} ${endpoint.summary || ''}`);
        }
      }
      
      return true;
    }
    
    // 实际发送请求
    const importUrl = `${config.yapiBaseUrl}${config.importUrl}`;
    log(`发送请求至: ${importUrl}`, colors.yellow);
    
    try {
      const response = await axios.post(importUrl, {
        type: 'swagger',
        token: config.yapiToken,
        json: JSON.stringify(swaggerData),
        merge: 'merger',  // merger=合并，good=智能合并
        server_url: '' // 设置为空表示使用Swagger中定义的server_url
      });
      
      log(`服务器响应: ${JSON.stringify(response.data)}`, colors.yellow);
      
      if (response.data && response.data.errcode === 0) {
        log('Swagger文档成功同步至YAPI!', colors.green);
        // 只在data有效时尝试访问pid
        if (response.data.data && response.data.data.pid) {
          log(`YAPI项目地址: ${config.yapiBaseUrl}/project/${response.data.data.pid}/interface/api`, colors.green);
        } else {
          log(`YAPI项目地址: ${config.yapiBaseUrl}`, colors.green);
        }
        return true;
      } else {
        log(`同步失败: ${response.data.errmsg || '未知错误'}`, colors.red);
        return false;
      }
    } catch (error) {
      log(`请求错误: ${error.message}`, colors.red);
      
      // 检查是否是连接错误
      if (error.code === 'ECONNREFUSED' || error.code === 'ENOTFOUND') {
        log(`无法连接到YApi服务器: ${config.yapiBaseUrl}`, colors.red);
        log('请确认YApi服务器已启动，或检查yapiBaseUrl配置是否正确', colors.yellow);
      }
      
      // 检查是否是超时
      if (error.code === 'ETIMEDOUT') {
        log(`连接YApi服务器超时`, colors.red);
      }
      
      // 详细错误信息
      if (error.response) {
        log(`服务器响应状态码: ${error.response.status}`, colors.red);
        log(`服务器响应内容: ${JSON.stringify(error.response.data)}`, colors.red);
      } else if (error.request) {
        log('已发送请求但没有收到响应', colors.red);
      }
      
      return false;
    }
  } catch (error) {
    log(`程序异常: ${error.message}`, colors.red);
    return false;
  }
}

/**
 * 主函数
 */
async function main() {
  log('Swagger到YAPI同步工具', colors.green);
  log('------------------------');
  
  // 验证配置
  validateConfig();
  
  // 读取Swagger文件
  const swaggerData = readSwaggerFile();
  log(`已读取Swagger文件: ${config.swaggerPath}`, colors.green);
  
  // 同步到YAPI
  const success = await syncToYapi(swaggerData);
  
  if (success) {
    log('同步完成!', colors.green);
    process.exit(0);
  } else {
    log('同步失败!', colors.red);
    process.exit(1);
  }
}

// 执行主函数
main().catch(error => {
  log(`程序异常: ${error.message}`, colors.red);
  process.exit(1);
}); 