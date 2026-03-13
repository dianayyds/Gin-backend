version - 20230615 by bo.zheng
# 暂无CICD,需本地build

## 本地debug配置：
    1. 修改目录名称：rap-backend -> rap_backend  ， (历史兼容问题) 
    2. 运行参数  rap_backend /Users/bz/airudder/rap_backend/docker/release/rap_backend.config
    3. 修改rap_backend.config中的配置
    4. 示例：
```
#LogFile: the log file store
LogFile = /Users/bz/airudder/rap_backend/log/rap_backend.log

#CfgRoot: the config root path
CfgRoot = /Users/bz/airudder/rap_backend/docker/rap_backend/release/etc_bak/
```
    4. 按需修改db.config、common.config、gauessdb.config、oss.config
    5. 运行rap_backend


## 生产发布流程
    1. 执行makeimge.sh
    2. 配置修改 - 直接去机器上修改!!!
    3. cmp发布
