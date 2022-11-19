基于Hyperledger fabric 网络模拟DID文档创建、获取过程<br>
<br>
文件目录：
- chiancode  存放链玛文件
- config 配置文件
- contract-sdk  SDK文件夹
- test-network 网络

操作步骤如下：
1. 启动网络 ./startFabric.sh
2. rm -rf wallet
3. cd contract-sdk
4. go run main.go

操作成功如下所示
```
2022/11/19 16:44:02 =======application-golang start==============
2022/11/19 16:44:02 ============ Populating wallet ============
certPath -> ../test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem
read cert file successful~~~
start connect network
 [fabsdk/core] 2022/11/19 08:44:02 UTC - cryptosuite.GetDefault -> INFO No default cryptosuite found, using default SW implementation
SDK初始化成功，成功接入Fabric网络
创建成功
创建私有DID耗时: 2.097522678s
获取creator doc 成功
创建公共DID总耗时： 2.076873004s
创建公有DID耗时: 2.076873004s
获取DID耗时: 10.229696ms
```