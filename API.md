### 接口

#### 服务器信息

- **端点**: `/server/info`
- **请求**:
  ```bash
  curl http://127.0.0.1:8080/server/info
  ```
- **描述**: 返回服务器名称与版本信息
- **响应**

  ```json
  {
    "name": "YeGame Group",
    "version": "v0.1.3.0"
  }
  ```

#### 玩家列表

- **端点**: `/player`
- **请求**:

  ```bash
  curl http://127.0.0.1:8080/player
  curl http://127.0.0.1:8080/player\?update\=true
  ```

- **查询参数**:
  - `update`（可选）: 一个布尔值（`"true"`）表示是否在请求时从服务器更新玩家数据。默认为 false。
- **描述**: 获取所有玩家的昵称、steamid、playeruid 和上次在线时间与当前在线情况（最后五分钟内在线也算作在线）。
- **响应**:

  ```json
  [
    {
      "last_online": "2024-01-26 13:43:33",
      "name": "全国可飞",
      "online": true,
      "playeruid": "357689484",
      "steamid": "xxx"
    },
    {
      "last_online": "2024-01-26 13:43:33",
      "name": "梵音丶",
      "online": true,
      "playeruid": "2144044083",
      "steamid": "xxx"
    },
    {
      "last_online": "2024-01-26 13:43:33",
      "name": "DZ",
      "online": true,
      "playeruid": "850234947",
      "steamid": "xxx"
    },
    {
      "last_online": "2024-01-25 21:15:44",
      "name": "宅记",
      "online": false,
      "playeruid": "1302283639",
      "steamid": "xxx"
    },
    {
      "last_online": "2024-01-25 21:06:53",
      "name": "ikun",
      "online": false,
      "playeruid": "00000000",
      "steamid": "<null/err>"
    }
  ]
  ```

#### 踢出玩家

- **端点**: `/player/:steamid/kick`
- **请求**:
  ```bash
  curl -X POST http://127.0.0.1:8080/player/:steamid/kick
  ```
- **路径参数**:
  - `steamid`: 要踢出的玩家的 SteamID/PlayerUID。
- **描述**: 使用玩家的 SteamID/PlayerUID 将玩家从服务器踢出。
- **响应**:
  ```json
  { "message": "踢出成功" }
  ```
  ```json
  { "error": "Failed to Kick: {id}" }
  ```

#### 封禁玩家

- **端点**: `/player/:steamid/ban`
- **请求**:
  ```bash
  curl -X POST http://127.0.0.1:8080/player/:steamid/ban
  ```
- **路径参数**:
  - `steamid`: 要封禁的玩家的 SteamID/PlayerUID。
- **描述**: 使用玩家的 SteamID/PlayerUID 封禁玩家。
- **响应**:
  ```json
  { "message": "封禁成功" }
  ```
  ```json
  { "error": "Failed to Ban: {id}" }
  ```

#### 广播消息

- **端点**: `/broadcast`
- **请求**:
  ```bash
  curl -X POST http://127.0.0.1:8080/broadcast -d '{"message": "Hello World"}'
  ```
- **请求体**:
  - `message`: 要广播的消息，暂不支持中文！
- **描述**: 向服务器上的所有玩家广播消息。
- **响应**:
  ```json
  { "message": "广播成功" }
  ```
  ```json
  { "error": "..." }
  ```

#### 关闭服务器

- **端点**: `/server/shutdown`
- **请求**:
  ```bash
  curl -X POST http://127.0.0.1:8080/shutdown -d '{"seconds": "60","message": "Shutdown in 60 sec"}'
  ```
- **请求体**:
  - `seconds`: 服务器关闭之前的倒计时时间（默认值："60"）。
  - `message`: 关闭前显示的消息。
- **描述**: 安排一个带有自定义倒计时和消息的服务器关闭。
- **响应**:
  ```json
  { "message": "关闭服务器成功" }
  ```
  ```json
  { "error": "..." }
  ```
