package db

//func SendToKafka(topic,key string,data []string){
//	host:=KC.Broker
//	address:=strings.Split(host,",")
//	if len(address) == 0{
//		log.Error("Kafka主机地址配置错误!")
//		return
//	}
//	//设置配置
//	config := sarama.NewConfig()
//	//等待服务器所有副本都保存成功后的响应
//	config.Producer.RequiredAcks = sarama.WaitForAll
//	//随机的分区类型
//	config.Producer.Partitioner = sarama.NewRandomPartitioner
//	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
//	config.Producer.Return.Successes = true
//	config.Producer.Return.Errors = true
//	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
//	//config.Version = sarama.V2_2_0_0
//	config.Version = sarama.V0_8_2_0
//
//	//使用配置,新建一个异步生产者
//	producer, e := sarama.NewAsyncProducer(address, config)
//	if e != nil {
//		log.Error("创建生产者失败: ", e)
//		return
//	}
//	defer producer.AsyncClose()
//
//	//发送的消息,主题,key
//	msg := &sarama.ProducerMessage{
//		Topic:     topic,
//		Key:       sarama.StringEncoder(key),
//		Timestamp: time.Now(),
//	}
//
//	for _, v := range data {
//		//fmt.Println("Message is: ", v)
//		//将字符串转化为字节数组
//		//myVal,_:=json.Marshal(v)
//		//msg.Value = sarama.StringEncoder(myVal)
//		msg.Value = sarama.StringEncoder(v)
//
//		//使用通道发送
//		producer.Input() <- msg
//
//		//循环判断哪个通道发送过来数据.
//		select {
//		//suc :=
//		case <-producer.Successes():
//			//fmt.Println("offset: ", suc.Offset, "partitions: ", suc.Partition)
//			//fmt.Println("Send to Kafka success!!!")
//		case fail := <-producer.Errors():
//			log.Error("err: ", fail.Err)
//
//		}
//	}
//}

// 接收kafka里的数据,将数据导入ES
//func Consumer(topic string, index string) {
//	for {
//
//		host := KC.Broker
//		address := strings.Split(host, ",")
//		if len(address) == 0 {
//			log.Error("Kafka主机地址配置错误!")
//			return
//		}
//		//设置配置
//		config := sarama.NewConfig()
//		//接收失败通知
//		config.Consumer.Return.Errors = true
//		//提交 offset 的时间间隔
//		//config.Consumer.Offsets.CommitInterval = 100 * time.Millisecond
//		//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
//		config.Version = sarama.V2_2_0_0
//		//config.Version = sarama.V0_8_2_0
//
//		//新建一个消费者
//		consumer, e := sarama.NewConsumer(address, config)
//		if e != nil {
//			log.Error("Error get consumer: ", e)
//			return
//		}
//		defer func() error {
//			if err := consumer.Close(); err != nil {
//				log.Error(err.Error())
//				return err
//			}
//			return nil
//		}()
//
//		//根据消费者获取指定的主题分区的消费者,Offset这里指定为获取最新的消息.
//		partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
//		if err != nil {
//			log.Error("Error get partition consumer", err)
//			return
//		}
//
//		// Trap SIGINT to trigger a shutdown.
//		signals := make(chan os.Signal, 1)
//		signal.Notify(signals, os.Interrupt)
//
//		//bulkRequest := ESClient.Bulk()
//		count := 0
//		//循环等待接受消息.
//	ConsumerLoop:
//		for {
//			select {
//			//接收消息通道和错误通道的内容.
//			case msg := <-partitionConsumer.Messages():
//				//fmt.Println(" timestrap: ", msg.Timestamp.Format("2006-01-02 15:04:05"), " value: ", string(msg.Value))
//				var netlog models.NetLog
//				id := string(msg.Key)
//				json.Unmarshal(msg.Value, &netlog)
//				//fmt.Println(netlog)
//				// 方法一
//				//res, err := ESClient.CreateIndex(EC.Index).
//				//	BodyJson(netlog).Do(context.Background())
//				//if err != nil {
//				//	log.Error("新建索引失败: ", err)
//				//	return
//				//}
//				//fmt.Println(res.Acknowledged)
//				// 方法二
//				put,err:=ESClient.Index().
//					Index(index).Type(index).
//					Id(id).BodyJson(netlog).
//					Do(context.Background())
//				if err!=nil{
//					log.Error("新建索引,写入ES失败: ", err)
//					return
//				}
//				count++
//				fmt.Printf("Indexed tweet %s to index %s, type %s\n and count=%d\n",
//					put.Id, put.Index, put.Type,count)
//
//				// 方法三
//				//doc := elastic.NewBulkIndexRequest().
//				//	Index(index).Type(index).
//				//	Id(id).Doc(netlog)
//				//bulkRequest = bulkRequest.Add(doc)
//				//count++
//				//fmt.Println(index," 计数: ", count)
//				//for k,v:=range mp{
//				//	doc:=elastic.NewBulkIndexRequest().
//				//		Index(EC.Index).Type(tp).
//				//		Id(k).Doc(v)
//				//	bulkRequest=bulkRequest.Add(doc)
//				//}
//				//res,err:=bulkRequest.Do(context.TODO())
//				//if err!=nil{
//				//	fmt.Println("批量插入失败: ",err)
//				//	return
//				//}
//				//failed:=res.Failed()
//				//it:=len(failed)
//				//fmt.Println("Error: ",res.Errors,it)
//
//				//fmt.Println("Count: ", count, " offset: ", msg.Offset, " Put into HBase success!")
//			case err := <-partitionConsumer.Errors():
//				log.Error("Get info from kafka failed", err.Err)
//				return
//			case <-signals:
//				break ConsumerLoop
//			}
//		}
//
//		//res, err := bulkRequest.Do(context.TODO())
//		//if err != nil {
//		//	fmt.Println("批量插入失败: ", err)
//		//	return
//		//}
//		//fmt.Println("result: ",res.Took)
//
//		ticker := time.NewTicker(10 * time.Minute)
//		<-ticker.C
//	}
//
//}
