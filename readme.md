* 為什麼使用 rabbitmq
  1. 異步處理： 當你的應用程式需要處理耗時的任務，如圖像處理、文件轉換等，可以將這些任務放入消息隊列，讓消費者進行處理，從而提高應用程式的反應能力。
  2. 分佈式系統： 如果你正在構建一個分佈式系統，不同的服務需要進行通信，你可以使用 RabbitMQ 來進行服務之間的消息傳遞，以實現解耦和彈性。
  3. 事件驅動架構： 當你希望應用程式能夠對特定事件作出反應時，如用戶註冊、訂閱更新等，你可以使用 RabbitMQ 來發送事件消息，通知相關應用程式進行相應操作。
  4. 訊息通知： 如果你需要向用戶發送通知消息，如郵件、短信、推送通知等，可以使用 RabbitMQ 來將消息交給消息服務進行處理，確保消息可靠地傳遞。
  5. 日誌處理： 將應用程式的日誌信息發送到 RabbitMQ，可以讓你將日誌集中處理，進行分析和監控。
  6. 負載平衡： 當你需要將請求分配給多個服務實例時，可以使用 RabbitMQ 來實現請求的負載平衡。
  7. 延遲處理： 如果你需要實現延遲處理，例如在一段時間後執行某些任務，可以使用 RabbitMQ 的延遲隊列功能。
  8. 微服務架構： 在微服務架構中，不同的微服務之間可能需要通信，RabbitMQ 可以幫助實現這種通信，從而保持微服務之間的解耦。
  9. 即時數據流： 當你需要處理即時的數據流，如實時分析、監控系統等，RabbitMQ 可以幫助你管理和分發數據。

* 核心概念
  * 製作人從不直接向佇列傳送任何訊息。 實際上，製片人往往甚至不知道訊息是否會傳遞到任何佇列。

* 基本名詞及解釋
  * Producer（生產者）： 生產者是發送消息到 RabbitMQ 交換器的應用程式或服務。它將消息發布到交換器，然後由 RabbitMQ 根據交換器的規則將消息路由到一個或多個隊列。
  * Consumer（消費者）： 消費者是從 RabbitMQ 隊列接收並處理消息的應用程式或服務。它訂閱一個隊列，等待並接收從隊列中傳遞的消息，然後進行相應的處理。
  * Message（消息）： 消息是應用程式之間交換的信息單元。它可以是任何形式的數據，如文字、JSON、二進制等。生產者創建消息並將其發送到交換器，然後 RabbitMQ 將消息路由到隊列。
  * Queue（隊列）： 隊列是 RabbitMQ 中用於存儲消息的地方。消息進入隊列後，等待消費者來接收和處理。每個隊列都有一個名稱，消費者可以訂閱特定的隊列。
  * Exchange（交換器）： 交換器是消息的分發中心。生產者將消息發送到交換器，然後交換器根據規則（如路由鍵）將消息路由到一個或多個隊列。交換器的類型（直接、主題、扇形等）決定了消息的路由方式。
  * Binding（綁定）： 綁定是交換器和隊列之間的連接。它確定了交換器如何將消息路由到特定的隊列。綁定包括交換器名稱、隊列名稱和一個可選的路由鍵。
  * Routing Key（路由鍵）： 路由鍵是一個字符串，用於將消息從交換器路由到特定的隊列。在發布消息時，生產者會指定一個路由鍵，交換器根據這個路由鍵將消息分發到相應的隊列。
  * Acknowledgment（確認）： 確認是消費者告知 RabbitMQ 已成功處理消息的方式。當消費者成功處理一條消息後，它可以發送確認給 RabbitMQ，然後 RabbitMQ 將從隊列中刪除該消息。
  * Dead-Letter Exchange（死信交換器）： 死信交換器是一個特殊的交換器，用於重新定向無法被正常處理的消息。如果消息在一個隊列中達到最大重新嘗試次數或因某些原因無法處理，它會被重新路由到死信交換器。
  * Virtual Host（虛擬主機）： 虛擬主機是 RabbitMQ 中的一個邏輯分割，用於隔離不同應用程式之間的資源。每個虛擬主機擁有自己的隊列、交換器和用戶許可權，從而實現了多租戶的隔離。

* docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management
* 1: 一送一接
  * go run receive.go
    * 啟動接收器
  * go run send.go
    * 發送一則訊息
* 2: 一送多接，但是都是一個接
  * go run new_task.go hello world
    * go run new_task.go First message.
    * go run new_task.go Second message..
    * go run new_task.go Third message...
    * go run new_task.go Fourth message....
    * go run new_task.go Fifth message.....
  * go run worker.go
    * go run worker.go
    * go run worker.go
  * 基本啟動多個接收的程式，會自動平均分配給每個接收的程式
  * 使用第一個教學，啟動多個 receive.go 也有相同效果
  * 學習設置持久化，在 worker 關閉後，因為沒有確認訊息處理，所以會重新在隊列中排序
  * 在訊息本身也需要做設定（new_task）
* 3. 一送多接，不過大家都接，廣播模式
  * go run 3_emit_log.go hello world
  * go run 3_receive_logs.go
* 4. 使用交換機來做路由，要完全相同才會處理
  * go run 4_receive_logs_direct.go info warn error
    * 從 info 開始可以選擇登入什麼樣的 exchange key 也就是路由
    * 這個會接收 info warn error 三種
  * go run 4_receive_logs_direct.go error
    * 這個只會接收 error 的
  * go run 4_emit_log_direct.go info "Run. Run. Or it will explode."
    * 當中的 info 就是 routing key 會說送到哪個接收器
  * go run 4_emit_log_direct.go error "Run. Run. Or it will explode."
* 5. 使用交換機來做路由，部分相同也可以
  * 5_receive_logs_topic
    * go run 5_receive_logs_topic.go "kern.*" "*.critical"
  * 5_emit_log_topic
    * go run 5_emit_log_topic.go kern.critical "A critical kernel error"