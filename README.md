"# logger" 
# 使用方法
### 基本使用
```
err := logger.InitLogger(logger.Config{
 		Method:       "file",
 		LogPath:      "log/",//需要预先建立log文件夹
 		LogName:      "test",
 		LogLevel:     logger.LogLevelInfo,
 		LogSplitType: logger.LogSplitTypeSize,
 		LogSplitSize: 50 << 20,
 	})
 	if err != nil {
 		log.Fatal(err)
 	}
logger.Info("user server is running")
```