## 第一章 前言

此SDK文档面向Hyperchain区块链平台的应用开发者，提供Hyperchain GoSDK的使用指南。f

**注意**：为支持Solidity方法和事件重载，1.2.15版本的GoSDK不再推荐使用ABI.Method[methodName]和ABI.Event[eventName]获取方法和事件对象，相应地，推荐使用方法和事件签名，并使用ABI.GetMethod(signature)和ABI.GetEvent(signature)获取相应的对象，签名的规则详见本文 *2.4 调用合约*。

## 第二章 接口使用流程示例

### 2.1 创建账户

账户分为国密和非国密类型，具体的账户创建可参考3.9文档。

创建SDK通用的accountJSON，再转换为GoSDK使用的账户Key。

例如：

```go
accountJson, err := account.NewAccountSm2("")
key, err := account.NewAccountSm2FromAccountJSON(accountJson,  "")
```

初始化RPC结构体来发送请求，具体初始化参考3.1文档。

```go
rpcAPI := rpc.NewRPC()
```

### 2.2 编译合约

创建合约文件，合约可分为Solidity、Java和HVM合约三种，其中Solidity合约支持在线编译，具体编译可查看3.3.1文档。

远程编译：

```go
rpcAPI := rpc.NewRPCWithPath("../../conf")
path := "../../conf/contract/Accumulator.sol"
contract, _ := common.ReadFileAsString(path)
cr, err := rpcAPI.CompileContract(contract)
fmt.Println("abi:", cr.Abi[0])
fmt.Println("bin:", cr.Bin[0])
fmt.Println("type:", cr.Types[0])
```

本地编译Solidity合约需要安装solc编译器，可以使用`npm install -g solc`安装 ，推荐使用本地编译。

Java合约编译则只能在本地编译，目前GoSDK暂不支持编译Java合约。

HVM合约的编写和编译请参考[http://hvm.internal.hyperchain.cn](http://hvm.internal.hyperchain.cn)

注意HVM合约编译之后需要用hvm-abi的maven插件处理HVM合约的jar包得到HVM合约的abi，具体参考*gosdk的HVM使用文档*。

### 2.3 部署合约

根据合约的构造函数是否有参数，实例化合约部署使用不同的方法，具体部署接口查看3.2.2文档。

```javascript
transaction := rpc.NewTransaction(gmKey.GetAddress()).Deploy(bin)
transaction.Sign(gmKey)
tx, err := rpcAPI.DeployContract(transaction)
fmt.Println(tx.ContractAddress)
```

其中VMType默认使用EVM，若为Java合约，则需要显示声明使用JVM，实例化transaction后部署。

```go
tx := rpc.NewTransaction(gmKey.GetAddress()).Deploy(payload).VMType(rpc.JVM)
tx.Sign(gmKey)
txReceipt, err := rpcAPI.DeployContract(tx)
contractAddress := txReceipt.ContractAddress
```

HVM合约需要显式声明使用HVM，实例化transaction的部署

```go
    transaction := NewTransaction(guomiKey.GetAddress()).Deploy(payload).VMType(HVM)
    transaction.Sign(guomiKey)
    receipt, stdErr := rpc.DeployContract(transaction)
```

### 2.4 调用合约

调用合约需要封装方法名和参数，Solidity使用`abi.Pack(name string, args ...interface{})`来编码调用的方法和参数，Java合约使用`java.EncodeJavaFunc(methodName string, params ...string)`来编码，编码后作为交易体paylod来进行调用,HVM合约则需使用`GenPayload(beanAbi *BeanAbi, params ...interface{})`来编码调用的abi和参数。

```go
packed, _ = ABI.Pack("getSum")
transaction2 := rpc.NewTransaction(guomiKey.GetAddress()).Invoke(receipt.ContractAddress, packed)
transaction2.Sign(guomiKey)
receipt2, _ := rpcAPI.InvokeContract(transaction2)
```

调用合约时，支持同名重载方法的调用，此时方法名称需要替换为方法签名：

```go
// 原方法调用
packed, _ = ABI.Pack("testOverload", args...)

// 现方法调用
packed, _ = ABI.Pack("testOverload(int24,bool[])", args...)
```

说明：

1. 若调用时仅仅使用方法名称，则默认调用无参同名方法，若未找到无参同名方法，则调用最后一个匹配到的方法，即最后声明的同名方法；

1. 方法签名的构造方式为：

```go
"MethodName(inputType1,inputType2,...inputTypeN)"
```

即在方法名称之后使用括号将传参类型统一包括，类型之间以逗号分隔。

1. 方法签名括号中规定了方法传参类型，传参类型之间必须使用“,”分隔，传参类型必须与abi中inputs里对应的type一致，且数量和顺序必须保持一致。

1. 事件签名构造同理：

```go
"EventName(inputType1,inputType2,...inputTypeN)"
```

### 2.5 返回值解析

得到返回值结果后，获得状态码可以判断是否调用成功，若调用成功，解析返回值可看到调用之后的结果。

```go
var p uint32
if sysErr := ABI.UnpackResult(&p, "getSum", receipt2.Ret); sysErr != nil {
    fmt.Println(sysErr)
    return
}
fmt.Println(p)
```

### 2.6 完整示例（Solidity合约）

2.6.1 Solidity合约

```javascript
contract Accumulator {
    event sayHello(int64 addr1, bytes8 indexed msg);
    int64 sum = 0;
    bytes hello = "hello world";
    function Accumulator(int64 sum1, bytes hello1) {
        sum = sum1;
        hello = hello1;
    }
    function getSum() returns(int64) {return sum;}
    function getHello() constant returns(bytes32) {
        sayHello(1, "test");
        return "hello";
    }
    function getMul(int64 count) returns(bytes, int64, address) {
        return (hello, count + sum, msg.sender);
    }
    function add(uint32 num1,uint32 num2) {sum = sum+num1+num2;}
}
```

2.6.2 GoSDK使用示例

```go
package main
import (
    "github.com/meshplus/gosdk/rpc"
    "github.com/meshplus/gosdk/account"
    "github.com/meshplus/gosdk/utils/common"
    "fmt"
    "strings"
    "github.com/meshplus/gosdk/abi"
)
func main() {
    logger := common.GetLogger("main")
    // new account
    accountJson, err := account.NewAccountSm2("123")
    if err != nil {
        logger.Error(err)
        return
    }
    fmt.Println("accountJson:", accountJson)
    key, err := account.NewAccountSm2FromAccountJSON(accountJson, "123")
    if err != nil {cc
        logger.Error(err)
        return
    }
    fmt.Println("account address:", key.GetAddress())
    rpcAPI := rpc.NewRPCWithPath("../conf")
    // compile
    code, _ := common.ReadFileAsString("../conf/contract/Accumulator3.sol")
    cr, stdErr := rpcAPI.CompileContract(code)
    if stdErr != nil {
        logger.Error(stdErr)
        return
    }
    fmt.Println("contract abi:", cr.Abi[0])
    // deploy
    tranDeploy := rpc.NewTransaction(key.GetAddress()).Deploy(cr.Bin[0]).DeployArgs(cr.Abi[0], int64(1), []byte("demo"))
    tranDeploy.Sign(key)
    if stdErr != nil {
        logger.Error(stdErr)
        return
    }
    txDeploy, stdErr := rpcAPI.DeployContract(tranDeploy)
    if stdErkr != nil {
        logger.Error(stdErr)
        return
    }
    fmt.Println("contract address:", txDeploy.ContractAddress)
    // invoket
    ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))
    packed, err := ABI.Pack("getMul", int64(1))
    tranInvoke := rpc.NewTransaction(key.GetAddress()).Invoke(txDeploy.ContractAddress, packed)
    tranInvoke.Sign(key)
    txInvoke, stdErr := rpcAPI.InvokeContract(tranInvoke)
    if stdErr != nil {
        logger.Error(stdErr)
        return
    }
    fmt.Println("invoke transaction hash:", txInvoke.TxHash)
    // decode
    var p0 []byte
    var p1 int64
    var p2 common.Address
    result := []interface{}{&p0, &p1, &p2}
    if err = ABI.UnpackResult(&result, "getMul", txInvoke.Ret); err != nil {
        logger.Error(err)
        return
    }
    fmt.Println("p0, p1, p2:", string(p0), p1, p2.Hex())
}
```

### 2.7 完整示例（Java合约）

本节展示了使用GoSDK操作java合约的完整例子。

2.7.1 java合约

```java
public class SimulateBank extends ContractTemplate {
    public SimulateBank() {}
    //String account, double num
    private ExecuteResult issue(List<String> args) {
        if(args.size() != 2) {
            logger.error("args num is invalid");
            return result(false, "args num is invalid");
        }
        logger.info("account: " + args.get(0));
        logger.info("num: " + args.get(1));
        boolean rs = ledger.put(args.get(0).getBytes(), args.get(1).getBytes());
        if(rs == false) {
            logger.error("issue func error");
            return result(false, "put data error");
        }
        return result(true);
    }
    //String accountA, String accountB, double num
    private ExecuteResult transfer(List<String> args) {
        try {
            String accountA = args.get(0);
            String accountB = args.get(1);
            double num = Double.valueOf(args.get(2));
            Result result = ledger.get(accountA.getBytes());
            if(!result.isEmpty()) {
                double balanceA = result.toDouble();
                result = ledger.get(accountB.getBytes());
                double balanceB ;
                if(!result.isEmpty()){
                    balanceB = result.toDouble();
                    if (balanceA >= num) {
                        ledger.put(accountA, balanceA - num);
                        ledger.put(accountB, balanceB + num);
                    }
                }
            }else {
                String msg = "get account " + accountA  + " balance error";
                logger.error(msg);
                return result(false, msg);
            }
        }catch (Exception e) {
            logger.error(e.getMessage());
            return result(false, e.getMessage());
        }
        return result(true);
    }
    private ExecuteResult getAccountBalance(List<String> args) {
        if(args.size() != 1) {
            logger.error("args num is invalid");
        }
        try {
            Result result = ledger.get(args.get(0).getBytes());
            if (!result.isEmpty()) {
                return result(true, result.toDouble());
            }else {
                String msg = "getAccountBalance error no data found for" + args.get(0);
                logger.error(msg);
                return result(false, msg);
            }
        }catch (Exception e) {
            e.printStackTrace();
            return result(false, e);
        }
    }
    public ExecuteResult testPostEvent(List<String> args) {
        logger.info(args);
        for (int i = 0; i < 10; i ++) {
            Event event = new Event("event" + i);
            event.addTopic("simulate_bank");
            event.addTopic("test");
            event.put("attr1", "value1");
            event.put("attr2", "value2");
            event.put("attr3", "value3");
            ledger.post(event);
        }
        return result(true);
    }
}
```

2.7.2 GoSDK代码

```go
package main
import (
    "github.com/meshplus/gosdk/rpc"
    "github.com/meshplus/gosdk/utils/java"
    "github.com/meshplus/gosdk/utils/common"
    "github.com/meshplus/gosdk/account"
    "fmt"
)
func main() {
    // 获取logger
    logger := common.GetLogger("main")
    // 随机创建账户，获取account json
    accountJSON, sysErr := account.NewAccount("TomKK")
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    logger.Debugf(accountJSON)
    // 根据account json得到*ecdsa.Key结构体
    key, sysErr := account.NewAccountFromAccountJSON(accountJSON, "TomKK")
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    // 设置配置文件根目录，构建rpc结构体
    hrpc := rpc.NewRPCWithPath("../../conf")
    // 部署合约
    // 读取java合约，因为构造函数无参，故没有填写参数
    bin, sysErr := java.ReadJavaContract("../../conf/contract/contract01")
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    // 构建部署交易结构体
    tx := rpc.NewTransaction(key.GetAddress()).
        Deploy(bin).
        VMType(rpc.JVM)
    // 交易签名
    tx.Sign(key)
    // 向hyperchain发起部署交易
    txReceipt, stdErr := hrpc.DeployContract(tx)
    if stdErr != nil {
        logger.Error(stdErr.String())
        return
    }
    // 得到合约地址
    contractAddr := txReceipt.ContractAddress
    // 调用java合约的 issue 方法，为自己发行1000token
    // 构建调用交易结构体
    tx = rpc.NewTransaction(key.GetAddress()).
        Invoke(contractAddr, java.EncodeJavaFunc("issue", key.GetAddress(), "1000")).
        VMType(rpc.JVM)
    // 交易签名
    tx.Sign(key)
    // 向hyperchain发起调用合约
    txReceipt, stdErr = hrpc.InvokeContract(tx)
    if stdErr != nil {
        logger.Error(stdErr.String())
        return
    }
    // 查看自己账户内token余额
    tx = rpc.NewTransaction(key.GetAddress()).
        Invoke(contractAddr, java.EncodeJavaFunc("getAccountBalance", key.GetAddress())).
        VMType(rpc.JVM)
    // 交易签名
    tx.Sign(key)
    // 向hyperchain发起调用合约
    txReceipt, stdErr = hrpc.InvokeContract(tx)
    if stdErr != nil {
        logger.Error(stdErr.String())
        return
    }
    // 解码合约返回值
    fmt.Printf("getAccountBalance返回值解码前: %s\n", txReceipt.Ret)
    fmt.Printf("getAccountBalance返回值解码后: %s\n", java.DecodeJavaResult(txReceipt.Ret))
    // 调用java合约testPostEvent方法
    // 构造调用交易结构体
    tx = rpc.NewTransaction(key.GetAddress()).
        Invoke(contractAddr, java.EncodeJavaFunc("testPostEvent")).
        VMType(rpc.JVM)
    // 签名
    tx.Sign(key)
    // 向hyperchain发起调用交易请求
    txReceipt, stdErr = hrpc.InvokeContract(tx)
    if stdErr != nil {
        logger.Error(stdErr.String())
        return
    }
    size := len(txReceipt.Log)
    for i := 0; i < size; i++ {
        fmt.Printf("java log 解码前: %s\n", txReceipt.Log[i].Data)
        // 解码java合约log
        decoded, sysErr := java.DecodeJavaLog(txReceipt.Log[i].Data)
        if sysErr != nil {
            logger.Error(sysErr)
            return
        }
        fmt.Printf("java log 解码后: %s\n", decoded)
    }
}
```

### 2.8 完整示例（HVM合约）

本节展示了使用GoSDK操作HVM合约的完整例子。

2.8.1 HVM合约

```java
package cn.hyperchain.contract.invoke;
import cn.hyperchain.contract.BaseInvoke;
import cn.hyperchain.contract.logic.ArraysTest;
import cn.hyperchain.contract.logic.bean.Bean1;
import cn.hyperchain.contract.logic.bean.Person;
import cn.hyperchain.core.Logger;
import java.util.*;
public class EasyInvoke implements BaseInvoke<Boolean, ArraysTest> {
    private boolean aBool;
    private char aChar;
    private byte aByte;
    private short aShort;
    private int anInt;
    private long aLong;
    private float aFloat;
    private double aDouble;
    private String aString;
    private Person person;
    private Bean1 bean1;
    private List<String> strList;
    private List<Person> personList;
    private Map<String, Person> personMap;
    private Map<String, Bean1> bean1Map;
    public Boolean invoke(ArraysTest arraysTest) {
        Logger logger = Logger.getLogger(EasyInvoke.class);
        logger.error(this.aBool);
        logger.error(this.aChar);
        logger.error(this.aByte);
        logger.error(this.aShort);
        logger.error(this.anInt);
        logger.error(this.aLong);
        logger.error(this.aFloat);
        logger.error(this.aDouble);
        logger.error(this.aString);
        logger.error(this.person);
        logger.error(this.bean1);
        logger.error(this.strList);
        logger.error(this.personList);
        logger.error("personMap: " + this.personMap);
        logger.error("bean1Map: " + this.bean1Map);
        return true;
    }
}
```

2.8.2 gosdk代码

```go
package main
import (
    "fmt"
    "github.com/meshplus/gosdk/account"
    "github.com/meshplus/gosdk/common"
    "github.com/meshplus/gosdk/hvm"
    "github.com/meshplus/gosdk/rpc"
    "github.com/meshplus/gosdk/utils/java"
)
func main() {
    logger := common.GetLogger("main")
    accountJson, sysErr := account.NewAccount("sys")
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    logger.Debugf(accountJson)
    key, sysErr :=account.NewAccountFromAccountJSON(accountJson,"sys")
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    hrpc := rpc.NewRPCWithPath("./conf")
    jarPath := "/Users/songyu/Desktop/hvm-abi-demo/target/hvmDemo-1.0.jar"
    payload, sysErr := hvm.ReadJar(jarPath)
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    transaction := rpc.NewTransaction(key.GetAddress()).Deploy(payload).VMType(rpc.HVM)
    transaction.Sign(key)
    receipt, sysErr := hrpc.DeployContract(transaction)
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    //read abi
    abiPath := "/Users/songyu/Desktop/hvm-abi-demo/target/hvm.abi"
    abiJson, sysErr := common.ReadFileAsString(abiPath)
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    abi, sysErr := hvm.GenAbi(abiJson)
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    // fmt.Println(abi)
    easyBean := "cn.hyperchain.contract.invoke.EasyInvoke"
    beanAbi, sysErr := abi.GetBeanAbi(easyBean)
    // fmt.Println(beanAbi)
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    person1 := []interface{}{"tom", "21"}
    person2 := []interface{}{"jack", "18"}
    bean1 := []interface{}{"hvm-bean1", person1}
    bean2 := []interface{}{"hvm-bean2", person2}
    // fmt.Println(len(beanAbi.Inputs))
    invokePayload, sysErr := hvm.GenPayload(beanAbi, "true", "c", "20", "100", "1000", "10000", "1.1", "1.11", "string", 
        //或者调用hvm.Convert(）方法（推荐！）
        //person1 := []string{{"tom", "21"}
        //Convert(person1)
        person1, 
        bean1,
        []interface{}{"strList1", "strList2"},
        []interface{}{person1, person2},
        []interface{}{[]interface{}{"person1", person1}, []interface{}{"person2", person2}},
        []interface{}{[]interface{}{"bean1", bean1}, []interface{}{"bean2", bean2}})
    //下面这种方式也可以，GenPayLoad接受的是一个个Json字符串
    //invokePayload, sysErr := hvm.GenPayload(beanAbi, "true", "c", "20", "100", "1000", "10000", "1.1", "1.11", "string", `["1f","2f","3f"]`,`[1,2,3]`, `{789:{456:12.2},234:{345:12.2}}`, `[[1,2],[2,4]]`, `{"name":"tom","age":21}`, `{"beanName":"hvm-bean1","person":{"name":"tom","age":21}}`, `["strList1","strList2"]`, `[{"name":"tom","age":21},{"name":"jack","age":18}]`, `{"person1":{"name":"tom","age":21},"person2":{"name":"jack","age":18}}`, `{"bean1":{"beanName":"hvm-bean1","person":{"name":"tom","age":21}},"bean2":{"beanName":"hvm-bean2","person":{"name":"jack","age":18}}}`)
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    invokeTx := rpc.NewTransaction(key.GetAddress()).Invoke(receipt.ContractAddress, invokePayload).VMType(rpc.HVM)
    invokeTx.Sign(key)
    invokeRe, sysErr := hrpc.InvokeContract(invokeTx)
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
    fmt.Printf("java log 解码前: %s\n", invokeRe.Ret)
    fmt.Printf("java log 解码后：%s\n", java.DecodeJavaResult(invokeRe.Ret))
    }
```

## 第三章 SDK文档

### 3.1 初始化

3.1.1 配置文件说明

配置文件hpc.toml

```toml
title = "GoSDK configuratoin file"
namespace = "global"
#发送重新连接请求间隔(/ms)
reConnectTime = 10000
[jsonRPC]
    nodes = ["172.16.5.3","172.16.5.3","172.16.5.3","172.16.5.3"]
    # JsonRpc connect port
    ports = ["8081", "8082", "8083", "8084"]
[webSocket]
    # webSocket connect port
    ports = ["11001", "11002", "11003", "11004"]
[polling]
    #重发次数
    resendTime = 10
    #第一次轮训时间间隔 unit /ms
    firstPollingInterval = 100
    #发送一次,第一次轮训的次数
    firstPollingTimes = 10
    #第二次轮训时间间隔 unit /ms
    secondPollingInterval = 1000
    #发送一次,第二次轮训的次数
    secondPollingTimes = 10
[privacy]
    #send Tcert during the request or not
    sendTcert = true
    #if sendTcert is true , you should add follow path.
    #the paths followed are relative to conf root path
    sdkcertPath = "certs/sdkcert.cert"
    sdkcertPrivPath = "certs/sdkcert.priv"
#    sdkcertPath = "certs/sdkcert_cfca.cert"
#    sdkcertPrivPath = "certs/sdkcert_cfca.priv"
    uniquePubPath = "certs/unique.pub"
    uniquePrivPath = "certs/unique.priv"
    cfca = false
[security]
    #Use Https
    https = true
    #If https is true, you shoule add follow properties
    #the paths followed are relative to conf root path
    tlsca = "certs/tls/tlsca.ca"
    tlspeerCert = "certs/tls/tls_peer.cert"
    tlspeerPriv = "certs/tls/tls_peer.priv"
[log]
    #设置日志输出门槛
    #"CRITICAL","ERROR","WARNING","NOTICE","INFO","DEBUG",
    log_level = "DEBUG"
    #存放日志文件夹
    log_dir = "../logs"
[transport]
    # MaxIdleConns controls the maximum number of idle (keep-alive)
    # connections across all hosts. Zero means no limit.
    maxIdleConns = 0
    # MaxIdleConnsPerHost, if non-zero, controls the maximum idle
    # (keep-alive) connections to keep per-host. If zero,
    # DefaultMaxIdleConnsPerHost is used.
    maxIdleConnsPerHost = 10
```

**说明**：

- `namespace`为初始化时传入的namespace；

- jspnRPC模块：`nodes`表示平台各节点的IP；`ports`也就是平台发送消息的端口。

- webSocket模块：`ports`为事件订阅连接的端口，他们分别对应着node中配置的IP。

- polling模块：`resendTime`参数表示重发次数；`pollingInterval`表示轮训去获取交易的时间间隔（分为第一次和第二次），单位为毫秒；`pollingTimes`表示发送一次交易的轮训次数（分为第一次和第二次）。

- privacy模块：`sendTcert`表示开启Tcert，请配合平台一起使用；`sdkcertPath`和`sdkcertPrivPath`表示请求tcert时需要用来签名的公私钥对；`uniquePubPath`和`uniquePrivPath`表示请求tcert时所使用的请求体公私钥对；`cfca`表示开启cfca证书。**需要注意证书路径均相对于**`**hpc.toml**`**文件**。

- security模块：`https`开启后则请求都会使用https来发送，`tlsca`表示客户端需要验证的服务器ca，`tlspeerCert`表示服务器需要验证的SDK的ca的公钥，`tlspeerPriv`则为对应的ca公钥的私钥。

- log模块：`log_level`表示输出的日志级别。`log_dir`表示存放日志的文件夹路径。

**注意**：

所有配置文件、证书都应放在一个文件夹下，默认配置文件夹名为`conf`，`hpc.toml`在配置文件夹下根目录，所有配置为文件路径都应相对于配置文件夹存放，例如我们以默认文件夹名`conf`为例，则上述配置文件中各文件存放应为：

```javascript
conf
├── certs
│   ├── sdkcert.cert
│   ├── sdkcert.priv
│   ├── unique.priv
│   └── unique.pub
│   ├── tls
│   │   ├── tls_peer.cert
│   │   ├── tls_peer.priv
│   │   └── tlsca.ca
└── hpc.toml
```

3.1.2 初始化RPC结构体(默认路径)

`func NewRPC() *RPC`

- 说明：默认路径为初始化RPC处的上层文件夹下的`conf`文件，即相对路径`../conf/`，若配置文件夹在别处，则需使用传路径的方式。

- 返回【返回值1】：RPC结构体。

- 实例

```go
rpc.NewRPC()
```

3.1.3 初始化RPC结构体(带路径)

`func NewRPCWithPath(confRootPath string) *RPC`

- 说明：传入`conf`配置文件夹路径来获取RPC结构体。

- 参数【confRootPath】：配置文件夹的根路径，具体文件内容可参考3.1.1。

- 返回【返回值1】：RPC结构体。

- 实例

```go
rpc.NewRPCWithPath("../conf")
```

3.1.4 构建绑定节点的RPC结构体

`func (r *RPC) BindNodes(nodeIndexes ...int) (*RPC, error)`

- 说明：生成一个和部分节点绑定的RPC结构体。

- 参数【nodeIndexes】：节点编号。如nodeIndex为0，1，即对应配置文件中jsonRPC.nodes中的0号和1号节点，返回的RPC对象将只会在0号节点和1号节点间负载均衡发送请求。

- 返回【返回值1】：一个新的RPC对象，与原有RPC对象互不影响。

- 实例

```go
hrpc := rpc.NewRPCWithPath("../conf")
proxy,err := hrpc.BindNodes(0)
```

3.1.5 创建默认RPC结构体

`func DefaultRPC(nodes ...*Node) *RPC`

- 说明：生成一个带有默认参数的RPC结构体。

- 参数【nodes】：节点列表。可以通过`NewNode`接口创建，见3.1.6。

- 返回【返回值1】：RPC对象指针。

- 实例

```go
rpc := DefaultRPC(NewNode("localhost", "8081", "11001"))
```

3.1.6 创建Node结构体

`func NewNode(url string, rpcPort string, wsPort string) (node *Node)`

- 说明：生成一个代表节点的结构体Node，用于`DefaultRPC`接口参数。

- 参数【url】：节点ip。【rpcPort】：节点监听的json-rpc端口。【wsPort】：节点监听的web socket端口。

- 返回【返回值1】：Node对象指针。

- 实例

```go
rpc := DefaultRPC(NewNode("localhost", "8081", "11001"))
```

3.1.7 各配置项的setters

`func (rpc *RPC) Namespace(ns string) *RPC`

`func (rpc *RPC) ResendTimes(resTime int64) *RPC`

`func (rpc *RPC) FirstPollInterval(fpi int64) *RPC`

`func (rpc *RPC) FirstPollTime(fpt int64) *RPC`

`func (rpc *RPC) SecondPollInterval(spi int64) *RPC`

`func (rpc *RPC) SecondPollTime(spt int64) *RPC`

`func (rpc *RPC) ReConnTime(rct int64) *RPC`

- 说明：各个配置项的简单setters

3.1.8 RPC结构体https相关配置

`func (rpc *RPC) Https(tlscaPath, tlspeerCertPath, tlspeerPrivPath string) *RPC`

- 说明：https配置。

- 参数【tlscaPath】：tlsca.ca证书的路径。【tlspeerCertPath】：tls_peer.cert证书的路径。【tlspeerPrivPath】：tls_peer.priv证书的路径。

- 返回【返回值1】：RPC对象的指针，用于链式调用。

- 实例

```go
rpc := DefaultRPC(NewNode("localhost", "8081", "11001")).Https("../conf/certs/tls/tlsca.ca", "../conf/certs/tls/tls_peer.cert", "../conf/certs/tls/tls_peer.priv")
```

3.1.9 RPC结构体tcert相关配置

`func (rpc *RPC) Tcert(cfca bool, sdkcertPath, sdkcertPrivPath, uniquePubPath, uniquePrivPath string) *RPC`

- 说明：tert配置。

- 参数【cfca】：是否启用cfca。【sdkcertPath】：sdkcert.cert证书的路径。【sdkcertPrivPath】：sdkcert.priv证书的路径。【uniquePubPath】：unique.pub证书的路径。【uniquePrivPath】：unique.priv证书的路径。

- 返回【返回值1】：RPC对象的指针，用于链式调用。

- 实例

```go
rpc := DefaultRPC(NewNode("localhost", "8081", "11001")).Tcert(true, "../conf/certs/sdkcert.cert", "../conf/certs/sdkcert.priv", "../conf/certs/unique.pub", "../conf/certs/unique.priv")
```

3.1.10 向配置中增加节点

`func (rpc *RPC) AddNode(url, rpcPort, wsPort string) *RPC`

- 说明：增加节点。

- 参数见3.1.6。

- 返回【返回值1】：RPC对象的指针，用于链式调用。

- 实例

```go
rpc := DefaultRPC().AddNode("localhost", "8081", "11001")
```

### 3.2 Transaction相关接口

交易是与hyperchain交互的重要形式，在GoSDK中，交易由**Transaction**结构体代表。GoSDK中交易体用户可配置参数共有如下（以下字段名首字母大写即为该字段的设值函数）：

【from】：交易发起方账户地址。**必填**。

【to】：交易接收方地址，合约调用时为合约地址。默认**0x0**。

【value】：转账金额。默认为**0**。

【payload】：交易荷载，在部署、调用、更新合约时使用。默认为**空字符串**。

【opcode】：操作码，在合约管理时使用，**1**代表*更新合约*，**2**代表*冻结合约*，**3**代表*解冻合约*。默认为**0**。

【extra】：存证信息，可以在交易中附带一些信息。默认为**空字符串**。

【simulate】：配置当前交易是否为模拟交易。若为true，则交易数据不会存入链上数据库。默认为**false**。

【vmType】：当前交易针对的虚拟机，目前支持配置EVM、JVM、HVM三种。默认为**EVM**。

3.2.1 实例化交易

GoSDK中构造交易体采用的是链式调用，应该由**NewTransaction**开始。

`func NewTransaction(from string) *Transaction`

- 说明：每一笔交易都应该有交易发起方，应该在链式调用的开始就定义。

- 参数【from】：交易的发起方地址。

- 返回【返回值1】：交易结构体指针，可以继续链式调用。

- 实例

```go
transaction := rpc.NewTransaction(ecdsaKey.GetAddress())
```

3.2.1.1 实例化普通交易

普通交易即为向某个账户的转账交易。

*需要的交易体字段有（字段含义以及默认值请看3.2）：*

**必填**：from、to、value、simulate（默认false）

**选填**：extra

`func (t *Transaction) Transfer(to string, value int64) *Transaction`

- 说明：普通（转账）交易的便捷构造方法。

- 参数【to】：接收方地址。【value】：转账金额。

- 返回【返回值1】：交易结构体指针，可以继续链式调用。

- 实例

```go
// 便捷构造
transaction := rpc.NewTransaction(ecdsaKey.GetAddress()).
    Transfer("0xbfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(1)).
    Extra("存证信息")
fmt.Println(transaction)
```

```go
// 自定义构造
transaction = rpc.NewTransaction(ecdsaKey.GetAddress()).
    To("0xbfa5bd992e3eb123c8b86ebe892099d4e9efb783").
    Value(int64(1)).
    Extra("存证信息").
    Simulate(true)
fmt.Println(transaction)
```

3.2.1.2 实例化合约部署交易

合约部署交易用来将一个solidity合约或者java合约或HVM合约部署到区块链上。部署交易结构体的payload字段是编码过的合约代码，如果合约构造函数有参还需要添加构造参数。

*需要的交易体字段有（字段含义以及默认值请看3.2）：*

**必填**：from、payload、vmType（默认EVM）、simulate（默认false）

**选填**：extra

`func (t *Transaction) Deploy(payload string) *Transaction`

- 说明：部署交易的构造方法。

- 参数【payload】：编码后的合约。如果是**solidity合约**该字段应该为合约编译后的bin，如果是**java合约**该字段应该是java.ReadJavaContract（见工具方法）返回的编码后的字符串，如果是**HVM合约**，该字段应该是hvm.ReadJar。

- 返回【返回值1】：交易结构体指针，可以继续链式调用。

- 实例

```go
// 构造函数无参solidity合约部署
cr, err := compileContract("../../conf/contract/Accumulator.sol")
if err != nil {
    t.Error(err)
    return
}
transaction := rpc.NewTransaction(ecdsaKey.GetAddress()).Deploy(cr.Bin[0])
fmt.Println(transaction)
```

```go
// 构造函数无参java合约部署
payload, _ := java.ReadJavaContract("../../conf/contract/contract01")
transaction = rpc.NewTransaction(ecdsaKey.GetAddress()).
    Deploy(payload)
fmt.Println(transaction)
```

// 构造函数无参HVM合约部署

```go
payload, sysErr := hvm.ReadJar(jarPath)
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
transaction := rpc.NewTransaction(key.GetAddress()).Deploy(payload).VMType(rpc.HVM)
```

如果**solidity合约**构造函数有参数，那么应该继续链式调用DeployArgs

`func (t *Transaction) DeployArgs(abiString string, args ...interface{}) *Transaction`

- 说明：用来在部署**solidity合约**时设置合约构造函数的构造参数。

- 参数【abiString】：合约编译结果的abi字符串。【args】：构造参数，可以任意个数，类型需要和合约构造函数声明的对应。

- 返回【返回值1】：交易结构体指针，可以继续链式调用。

- HVM合约构造函数带参和不带参没有区别，都是通过abi进行调用的。

- 实例

```go
// 构造函数带参solidity合约部署
cr, _ = compileContract("../conf/contract/Accumulator2.sol")
var arg [32]byte
copy(arg[:], "test")
transaction = rpc.NewTransaction(ecdsaKey.GetAddress()).
    Deploy(cr.Bin[0]).
    DeployArgs(cr.Abi[0], uint32(10), arg)
fmt.Println(transaction)
```

```go
// 构造函数带参java合约部署
payload, _ = java.ReadJavaContract("../../conf/contract/contract01",  "1")
transaction = rpc.NewTransaction(ecdsaKey.GetAddress()).
    Deploy(payload)
fmt.Println(transaction)
```

// 构造函数带参HVM合约部署

```go
payload, sysErr := hvm.ReadJar(jarPath)
    if sysErr != nil {
        logger.Error(sysErr)
        return
    }
transaction := rpc.NewTransaction(key.GetAddress()).Deploy(payload).VMType(rpc.HVM)
```

3.2.1.3 实例化合约调用交易

合约调用交易用来调用区块链上solidity合约或者java合约某个合约方法。

*需要的交易体字段有（字段含义以及默认值请看3.2）：*

**必填**：from、to、payload、vmType（默认EVM）、simulate（默认false）

**选填**：extra

`func (t *Transaction) Invoke(to string, payload []byte) *Transaction`

- 说明：用来构造合约调用交易的构造方法

- 参数【to】：合约地址。【payload】：若为**solidity合约**，那么为经过abi.Pack(methodName, param1...)（构造abi见3.11.7）编码后的调用字符串；若为**java合约**，那么为java.EncodeJavaFunc（见3.11.4）编码后的字符串，如果是**HVM合约**，那么为hvm.GenPayload编码后的字符串。

- 返回【返回值1】：交易结构体指针，可以继续链式调用。

- 实例

```go
// solidity合约调用
cr, _ := compileContract("../../conf/contract/Accumulator.sol")
contractAddress, err := deployContract(cr.Bin[0], cr.Abi[0])
ABI, sysErr := abi.JSON(strings.NewReader(cr.Abi[0]))
if err != nil {
    t.Error(sysErr)
    return
}
packed, sysErr := ABI.Pack("add", uint32(1), uint32(2))
if err != nil {
    t.Error(sysErr)
    return
}
transaction := rpc.NewTransaction(ecdsaKey.GetAddress()).
    Invoke(contractAddress, packed))
fmt.Println(transaction)
```

```go
// java合约调用
tx = rpc.NewTransaction(ecdsaKey.GetAddress()).
    Invoke(contractAddress, java.EncodeJavaFunc("issue", ecdsaKey.GetAddress(), "1000")).
    VMType(rpc.JVM)
fmt.Println(tx)
```

```go
// HVM合约调用
    easyBean := "cn.hyperchain.contract.invoke.EasyInvoke"
    beanAbi, sysErr := abi.GetBeanAbi(easyBean)
    invokePayload, sysErr := hvm.GenPayload(beanAbi, "true", "c", "20", "100", "1000", "10000", "1.1", "1.11", "string", person1, 
    bean1,
        []interface{}{"strList1", "strList2"},
        []interface{}{person1, person2},
        []interface{}{[]interface{}{"person1", person1}, []interface{}{"person2", person2}},
        []interface{}{[]interface{}{"bean1", bean1}, []interface{}{"bean2", bean2}})
    //下面这种方式也可以，GenPayLoad接受的是一个个字符
    //invokePayload, sysErr := hvm.GenPayload(beanAbi, "true", "c", "20", "100", "1000", "10000", "1.1", "1.11", "string", `["1f","2f","3f"]`,`[1,2,3]`, `{789:{456:12.2},234:{345:12.2}}`, `[[1,2],[2,4]]`, `{"name":"tom","age":21}`, `{"beanName":"hvm-bean1","person":{"name":"tom","age":21}}`, `["strList1","strList2"]`, `[{"name":"tom","age":21},{"name":"jack","age":18}]`, `{"person1":{"name":"tom","age":21},"person2":{"name":"jack","age":18}}`, `{"bean1":{"beanName":"hvm-bean1","person":{"name":"tom","age":21}},"bean2":{"beanName":"hvm-bean2","person":{"name":"jack","age":18}}}`)
undefined
invokeTx := rpc.NewTransaction(key.GetAddress()).Invoke(receipt.ContractAddress, invokePayload).VMType(rpc.HVM)
fmt.Println(invokeTx)
```

其中HVM合约支持的类型有：

| 类 型  | 写法1                                                        | 写法2(json格式）                                             |
| ------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| Bool   | "true"                                                       | "true"                                                       |
| Char   | "c"                                                          | "c"                                                          |
| Short  | "20"                                                         | "20"                                                         |
| Int    | "20"                                                         | "20"                                                         |
| Float  | "1.1"                                                        | "1.1"                                                        |
| Double | "1.11"                                                       | "1.11"                                                       |
| List   | []interface{}{"strList1", "strList2"}                        | `["strList1","strList2"]`                                    |
| Map    | []interface{}{[]interface{}{"bean1", bean1}, []interface{}{"bean2", bean2}} | ``{789:{456:12.2},234:{345:12.2}}`                           |
| Struct | bean1 := []interface{}{"hvm-bean1", person}  （其中person是基本类型） | `{"bean1":{"beanName":"hvm-bean1","person":{"name":"tom","age":21}},"bean2":{"beanName":"hvm-bean2","person":{"name":"jack","age":18}}}` |
| Array  | array1 = []interface{}{"strList1", "strList2"}               | `["strList1","strList2"]`                                    |

如上图所示，写法一是传入interface数组，写法二是传入json字符串，或者传入golang内在支持的类型然后用hvm.Convert()转成interface数组。

3.2.1.4 实例化维护合约交易

维护合约即为升级、冻结、解冻合约三种操作。

*需要的交易体字段有（字段含义以及默认值请看3.2）：*

**必填**：from、to、opCode、vmType（默认EVM）、simulate（默认false）

**选填**：extra、payload（若为升级合约则需要）

`func (t *Transaction) Maintain(op int64, to, payload string) *Transaction`

- 说明：实例化维护合约交易的构造方法。

- 参数【op】：操作码opCode。【to】：合约地址。【payload】：若为升级合约则需要填写新的payload（同部署合约的payload规则）。

- 返回【返回值1】：交易结构体指针，可以继续链式调用。

- 当为升级合约时，payload为新合约的jar string（通过ReadJar函数得到），Maintain函数会用新合约的jar包替换掉目前已经部署合约的jar包。

- 注意：当为升级合约时，新合约必须包含旧合约中有storefield注解的变量，且变量的类型和名称不能改变。

- 实例

```text
updateTx := rpc.NewTransaction(key.GetAddress()).Maintain(1 , receipt.ContractAddress, updatePayload).VMType(rpc.HVM)
```

```go
// 便捷调用
transactionUpdate := rpc.NewTransaction(ecdsaKey.GetAddress()).
    Maintain(1, contractAddress, compileUpdate.Bin[0])
```

```go
// 自定义调用
transactionUpdate := rpc.NewTransaction(ecdsaKey.GetAddress()).
        OpCode(1).
        To(contractAddress).
        Payload(compileUpdate.Bin[0])
```

3.2.2 交易签名

所以通过GoSDK发往Hyperchain的交易都需要进行签名。

`func (t *Transaction) Sign(key interface{})`

- 说明：用账户私钥对某个交易进行签名。目前只支持ECDSA和SM2两种签名算法。

- 参数【key】：账户结构体（账户相关请见3.9），注意该账户应该和交易体的from字段使用的同一个账户，如果是*ecdsa.Key类型，那么将使用ECDSA算法签名，如果是*gm.Key类型，那么将使用SM2算法签名。

- 实例

```go
transaction := rpc.NewTransaction(ecdsaKey.GetAddress())
transaction.Sign(ecdsaKey)
```

3.2.3 根据区块号查询范围内的交易

`func (r *RPC) GetTransactionsByBlkNum(start, end uint64) ([]TransactionInfo, StdError)`

- 说明：根据区块号获取范围内的交易

- 参数【start】：起始区块号。【end】：结束区块号。

- 返回【返回值1】：该区块号区间内的交易列表。交易结构体详细信息请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txs, err := hrpc.GetTransactionsByBlkNum(block.Number-1, block.Number)
```

3.2.4  获取所有非法交易

`func (r *RPC) GetDiscardTx() ([]TransactionInfo, StdError)`

- 说明：获取所有的非法交易

- 返回【返回值1】：非法交易列表。交易结构体详细信息请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txs, err := hrpc.GetDiscardTx()
if err != nil {
    t.Error(err)
    return
}
fmt.Println(len(txs))
fmt.Println(txs[len(txs)-1].Hash)
```

3.2.5 根据哈希获取交易

`func (r *RPC) GetTransactionByHash(txHash string) (*TransactionInfo, StdError)`

- 说明：根据交易哈希获取交易

- 参数【txHash】：交易哈希，应该以0x开头。

- 返回【返回值1】：交易信息结构体指针。交易结构体详细信息请看4.1.2。

- 实例

```go
transaction := rpc.NewTransaction(ecdsaKey.GetAddress()).Deploy(binContract)
transaction.Sign(ecdsaKey)
receipt, _ := hrpc.DeployContract(transaction)
tx, err := hrpc.GetTransactionByHash(receipt.TxHash)
```

3.2.6 根据交易哈希列表批量获取交易

`func (r *RPC) GetBatchTxByHash(hashes []string) ([]TransactionInfo, StdError)`

- 说明：根据交易哈希列表批量获取交易

- 参数【hashes】：交易哈希列表，应该以0x开头。

- 返回【返回值1】：交易信息列表；交易结构体详细信息请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txs, err := hrpc.GetBatchTxByHash(txhashes)
```

3.2.7 通过区块hash和交易序号返回交易信息

`func (r *RPC) GetTxByBlkHashAndIdx(blkHash string, index uint64) (*TransactionInfo, StdError)`

- 说明：通过区块哈希和交易在区块中的索引查询交易信息。

- 参数【blkHash】：区块哈希。【index】：交易在区块中的索引值。

- 返回【返回值1】：交易体信息。交易结构体详细信息请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
info, err := hrpc.GetTxByBlkHashAndIdx(block.Hash, 0)
```

3.2.8 通过区块号和交易序号返回交易信息

`func (r *RPC) GetTxByBlkNumAndIdx(blkNum, index uint64) (*TransactionInfo, StdError)`

- 说明：通过区块号和交易在区块中的索引查询交易信息。

- 参数【blkNum】：区块号。【index】：交易在区块中的索引值。

- 返回【返回值1】：交易体信息。交易结构体详细信息请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
info, err := hrpc.GetTxByBlkNumAndIdx(block.Number, 0)
```

3.2.9 通过区块号区间获取交易平均处理时间

`func (r *RPC) GetTxAvgTimeByBlockNumber(from, to uint64) (uint64, StdError)`

- 说明：根据区块号区间获取交易平均处理时间，单位ms

- 参数【from】：开始区块号。【to】：结束区块号。

- 返回【返回值1】：从[from, to]区间内的交易平均处理时间，单位（ms）。【返回值2】：错误类型请见4.1.5。

- 实例

```go
time, err := hrpc.GetTxAvgTimeByBlockNumber(block.Number - 2, block.Number)
```

3.2.10 批量获取交易回执

现有两个接口实现了该功能：

`func (r *RPC) GetBatchReceipt(hashes []string) ([]TxReceipt, StdError)`

- 说明：批量获取交易回执。

- 参数【hashes】：交易哈希数组。需要以0x开头。

- 返回【返回值1】：交易回执切片；交易相关结构体详细信息请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txs, err := hrpc.GetBatchReceipt(hashes)
```

3.2.11 根据区块hash获取区块上的交易数

`func (r *RPC) GetBlkTxCountByHash(blkHash string) (uint64, StdError)`

- 说明：根据区块hash获取区块上的交易数

- 参数【blkHash】：区块哈希。

- 返回【返回值1】：该哈希对应区块上的交易数量。【返回值2】：错误类型请见4.1.5。

- 实例

```go
count, err := hrpc.GetBlkTxCountByHash(block.Hash)
```

3.2.12 获取链上所有交易数量

`func (r *RPC) GetTxCount() (*TransactionsCount, StdError)`

- 说明：获取链上所有的交易数量

- 返回【返回值1】：包含交易数量和响应时间戳。交易结构体相关请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txCount, err := hrpc.GetTxCount()
```

3.2.13 查询区块区间内某合约的交易数量

`func (r *RPC) GetTxCountByContractAddr(from, to uint64, address string, txExtra bool) (*TransactionsCountByContract, StdError)`

- 说明：查询[from, to]区块号区间内涉及合约地址为address的交易的数量。

- 参数【from】：起始区块号。【to】：结束区块号。【address】：合约地址，以0x开头。【txExtra】：是否一定需要包含extra字段。若为true，则返回包含extra字段的交易，若为false，那么返回所有的匹配的交易无论是否有extra字段。

- 返回【返回值1】：包含总的交易数量、最后一条交易在区块中的索引号、最后一条交易所在的区块号。交易结构体相关请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
count, err := hrpc.GetTxCountByContractAddr(block.Number-1, block.Number, cAddress, false)
```

3.2.14 根据时间范围查询交易列表

`func (r *RPC) GetTxByTime(start, end uint64) ([]TransactionInfo, StdError)`

- 说明：查询时间范围内的交易信息。

- 参数【start】：起始时间时间戳。【end】：结束时间时间戳。

- 返回【返回值1】：交易列表。交易结构体信息见4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
infos, err := hrpc.GetTxByTime(1, uint64(time.Now().UnixNano()))
```

3.2.15 获取下一页的交易

`func (r *RPC) GetNextPageTxs(blkNumber, txIndex, minBlkNumber, maxBlkNumber, separated, pageSize uint64, containCurrent bool, contractAddr string) ([]TransactionInfo, StdError)`

- 说明：分页查询下一页某个合约相关的交易。

- 参数【blkNumber】：从该区块开始计数。【txIndex】：起始交易在blkNumber号区块的位置偏移量。 【minBlkNumber】：截止计数的最小区块号。 【maxBlkNumber】：截止计数的最大区块号。 即如果在遍历到maxBlkNumber这个区块时还是没有查够所需要的交易数目，那么将会结束计数直接返回不再继续向后遍历。【separated】：表示要跳过的交易条数（一般用于跳页查询）。 【pageSize】：表示要返回的交易条数。 【containCurrent】：true表示返回的结果中包括blkNumber区块中位置为txIndex的交易，如果该条交易不是合约地址为address合约的交易，则不算入。 【contractAddr】：查询的合约地址。

- 返回【返回值1】：交易列表。交易结构体信息见4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
// 查询从1号区块的2号交易到4号区块范围内的第6到第11个关于cAddress的交易列表
infos, err := rpc.GetNextPageTxs(1, 2, 3, 4, 5, 6, false, cAddress)
```

3.2.16 获取上一页的交易

`func (r *RPC) GetPrevPageTxs(blkNumber, txIndex, minBlkNumber, maxBlkNumber, separated, pageSize uint64, containCurrent bool, contractAddr string) ([]TransactionInfo, StdError)`

- 说明：分页查询上一页某个合约相关的交易。

- 参数【blkNumber】：从该区块开始计数。【txIndex】：起始交易在blkNumber号区块的位置偏移量。 【minBlkNumber】：截止计数的最小区块号。 即如果在遍历到minBlkNumber这个区块时还是没有查够所需要的交易数目，那么将会结束计数直接返回不再继续向前遍历。【maxBlkNumber】：截止计数的最大区块号。 【separated】：表示要跳过的交易条数（一般用于跳页查询）。 【pageSize】：表示要返回的交易条数。 【containCurrent】：true表示返回的结果中包括blkNumber区块中位置为txIndex的交易，如果该条交易不是合约地址为address合约的交易，则不算入。 【contractAddr】：查询的合约地址。

- 返回【返回值1】：交易列表。交易结构体信息见4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
// 查询从6号区块的5号交易到4号区块范围内的第3个关于cAddress的交易列表
infos, err := rpc.GetPrevPageTxs(6, 5, 4, 7, 2, 1, false, cAddress)
```

3.2.17 通过交易哈希获取交易回执

`func (r *RPC) GetTxReceipt(txHash string) (*TxReceipt, StdError)`

- 说明：通过交易哈希获取交易回执

- 参数【txHash】：交易哈希。

- 返回【返回值1】：交易回执。交易相关结构体详情请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
info,err := r.GetTxReceipt(txHash)
```

3.2.18 同步发送转账交易

`func (r *RPC) SendTx(transaction *Transaction) (*TxReceipt, StdError)`

- 说明：同步发送转账交易

- 参数【transation】：普通交易结构体（构造普通交易结构体见3.2.1.1）

- 返回【返回值1】：交易回执。交易相关结构体详情请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
receipt, err := hrpc.SendTx(transaction)
```

3.2.19 异步发送转账交易



- 取消支持
- 实例


3.2.20 查询区块区间交易数量（by method ID）

`func (rpc *RPC) GetTransactionsCountByMethodID(from, to uint64, address string, methodID string) (*TransactionsCountByContract, StdError)`

- 说明：查询区块区间交易数量（by method ID）

- 参数【from】：起始区块号【to】：终止区块号【address】：合约地址【methodID】：合约方法id

- 返回【返回值1】：交易数量。交易相关结构体详情请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txCount,err := r.GetTransactionsCountByMethodID(from, to, address, methodID)
```

3.2.21 查询区块交易数量（by block number）

`func (rpc *RPC) GetBlkTxCountByNumber(blkNum string) (uint64, StdError)`

- 说明：通过区块number获取区块上交易数

- 参数【blkNum】：区块号

- 返回【返回值1】：交易数量。交易相关结构体详情请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txCount,err := r.GetBlkTxCountByNumber(blkNum)
```

3.2.22 获取交易签名哈希

`func (rpc *RPC) GetSignHash(transaction *Transaction) (string, StdError)`

- 说明：获取交易签名哈希

- 参数【transaction】：交易信息

- 返回【返回值1】：签名哈希。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txCount,err := r.GetSignHash(transaction)
```

3.2.23 查询指定时间区间内的非法交易

`func (rpc *RPC) GetDiscardTransactionsByTime(start, end uint64) ([]TransactionInfo, StdError)`

- 说明：查询指定时间区间内的非法交易

- 参数【start】：起始时间戳(单位ns)【end】：结束时间戳(单位ns)

- 返回【返回值1】：非法交易信息。交易相关结构体详情请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txCount,err := r.GetDiscardTransactionsByTime(start, end)
```

3.2.24 查询指定时间区间内的交易数量

`func (rpc *RPC) GetTransactionsCountByTime(startTime, endTime uint64) (uint64, StdError)`

- 说明：查询指定时间区间内的交易数量

- 参数【startTime】：起始时间戳(单位ns)【endTime】：结束时间戳(单位ns)

- 返回【返回值1】：交易数量。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txCount,err := r.GetTransactionsCountByTime(startTime, endTime)
```

3.2.25 查询指定区块范围内的非法交易数量

`func (rpc *RPC) GetInvalidTransactionsByBlkNumWithLimit(start, end uint64, metadata *Metadata) (*PageResult, StdError)`

- 说明：查询指定区块范围内的非法交易数量

- 参数【start】：起始区块号【end】：结束区块号

- 返回【返回值1】：该区块号区间内的交易列表。交易结构体详细信息请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
pageResult, err := rpc.GetInvalidTransactionsByBlkNumWithLimit(start, end, metadata)
```

3.2.26 根据区块号查询区块内的非法交易列表

`func (rpc *RPC) GetInvalidTransactionsByBlkNum(blkNum uint64) ([]TransactionInfo, StdError)`

- 说明：根据区块号查询区块内的非法交易列表

- 参数【blkNum】：目标区块号。

- 返回【返回值1】：该区块号内的非法交易列表。交易结构体详细信息请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txInfos, err := rpc.GetInvalidTransactionsByBlkNum(block.Number)
```

3.2.27 根据区块哈希查询区块内的非法交易列表

`func (rpc *RPC) GetInvalidTransactionsByBlkHash(hash string) ([]TransactionInfo, StdError)`

- 说明：根据区块哈希查询区块内的非法交易列表

- 参数【hash】：目标区块哈希。

- 返回【返回值1】：该区块号内的非法交易列表。交易结构体详细信息请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
txInfos, err := rpc.GetInvalidTransactionsByBlkHash(block.Hash)
```

3.2.28 获取链上的非法交易数

`func (r *RPC) GetInvalidTxCount() (uint64, StdError)`

- 说明：获取链上的非法交易数

- 参数 无。

- 返回【返回值1】：链上的非法交易数量。【返回值2】：错误类型请见4.1.5。

- 实例

```go
count, err := rpc.GetInvalidTxCount()
```

### 3.3 Contract相关接口

3.3.1 编译源代码

3.3.1.1 编译Solidity合约

`func (r *RPC) CompileContract(code string) (*`CompileResult`, StdError)`

- 说明：编译Solidity合约可以通过远程调用的方式来编译合约，获取合约的abi和bin。

- 参数【code】：合约源代码。

- 返回【返回值1】：合约编译结果，包括abi、bin和合约名称(type)。【返回值2】：错误类型请见4.1.5。

- 注意事项：返回结果为数组，可以一次编译多个合约

- 实例

```go
rpcAPI := rpc.NewRPCWithPath("../../conf")
path := "../../conf/contract/Accumulator.sol"
contract, _ := common.ReadFileAsString(path)
cr, err := rpcAPI.CompileContract(contract)
if err != nil {
    fmt.Println("can not get compile return, ", err.String())
}
fmt.Println("abi:", cr.Abi[0])
fmt.Println("bin:", cr.Bin[0])
fmt.Println("type:", cr.Types[0])
```

3.3.1.2 编译Java合约

Java合约暂时不提供编译接口调用，可以通过`javac`在本地编译后(需要hyperjvm jar包)将`class`文件放到指定目录。

编译后需要在指定目录下创建一个contract.properties的配置文件，配置合约名和合约主类全名。

配置文件内容如下：

```properties
contract.name=AccountSum
main.class=cn.hyperchain.jcee.contract.examples.sb.src.SimulateBank
```

编译后的文件目录如下：

```shell
├── cn
│   └── hyperchain
│       └── jcee
│           └── contract
│               └── examples
│                   └── sb
│                       └── src
│                           └── SimulateBank.class
└── contract.properties
```

3.3.1.3 编译HVM合约

HVM合约暂时不提供编译接口调用，在本地编译后将jar文件放到指定目录。

编译后需要使用hvm-abi插件(在pom.xml中加入如下内容)，通过执行*mvn hvm-abi*命令来获取HVM合约的abi.

```xml
<plugin>
    <groupId>cn.hyperchain.hvm</groupId>
    <artifactId>hvm-maven-plugin</artifactId>
    <version>0.0.1</version>
    <configuration>
        <jarFile>${project.basedir}/target/hvmDemo-1.0.jar</jarFile>
        <invokeBeanPath>${project.basedir}/target/classes</invokeBeanPath>
        <invokeBeanPackages>
            <param>cn.hyperchain.contract.invoke.ArraysTestInvoke</param>
            <param>cn.hyperchain.contract.invoke.InvokeBean1</param>
            <param>cn.hyperchain.contract.invoke.EasyInvoke</param>
        </invokeBeanPackages>
        <outputFile>${project.basedir}/target/hvm.abi</outputFile>
    </configuration>
</plugin>
```

在hvm.abi中

- version代表hvm.abi的版本

- beanName代表实现*BaseInvoke*接口的类的名称

- inputs代表invoke方法的参数  其中

  - name 代表变量名称

  - type代表变量类型

  - 如果是结构体则structname代表结构体类的名称

- output代表invoke方法的返回值

- structs里面包含HVM合约中的结构体

示例hvm.abi如下

```json
  {
    "version": "v1",
    "beanName": "cn.hyperchain.contract.invoke.EasyInvoke",
    "inputs": [
      {
        "name": "aBool",
        "type": "Bool",
        "structName": "boolean"
      },
      {
        "name": "aChar",
        "type": "Char",
        "structName": "char"
      },
      ...
      ...
      ...
      {
        "name": "strList",
        "type": "List",
        "properties": [
          {
            "name": "java.lang.String",
            "type": "String",
            "structName": "java.lang.String"
          }
        ]
      },
      {
        "name": "personList",
        "type": "List",
        "properties": [
          {
            "name": "cn.hyperchain.contract.logic.bean.Person",
            "type": "Struct",
            "structName": "cn.hyperchain.contract.logic.bean.Person"
          }
        ]
      },
      {
        "name": "bean1Map",
        "type": "Map",
        "properties": [
          {
            "name": "java.lang.String",
            "type": "String",
            "structName": "java.lang.String"
          },
          {
            "name": "cn.hyperchain.contract.logic.bean.Bean1",
            "type": "Struct",
            "structName": "cn.hyperchain.contract.logic.bean.Bean1"
          }
        ]
      }
    ],
    "output": {
      "name": "java.lang.Boolean",
      "type": "Bool",
      "structName": "java.lang.Boolean"
    },
    "classBytes": "ca",
    "structs": [
      {
        "name": "cn.hyperchain.contract.logic.bean.Bean1",
        "type": "Struct",
        "properties": [
          {
            "name": "beanName",
            "type": "String",
            "structName": "java.lang.String"
          },
          {
            "name": "person",
            "type": "Struct",
            "structName": "cn.hyperchain.contract.logic.bean.Person"
          }
        ]
      },
      {
        "name": "cn.hyperchain.contract.logic.bean.Person",
        "type": "Struct",
        "properties": [
          {
            "name": "name",
            "type": "String",
            "structName": "java.lang.String"
          },
          {
            "name": "age",
            "type": "Int",
            "structName": "int"
          }
        ]
      }
    ]
  },
```

3.3.2 部署合约

3.3.2.1 部署合约(同步)

`func (r *RPC) DeployContract(transaction *Transaction) (*TxReceipt, StdError)`

- 说明：部署合约需要初始化部署Transaction。

- 参数【transaction】：部署transaction。

- 返回【返回值1】：交易回执，回执结构体请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
transaction := rpc.NewTransaction(gmKey.GetAddress()).Deploy(bin)
transaction.Sign(gmKey)
tx, err := rpcAPI.DeployContract(transaction)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(tx.ContractAddress)
```


3.3.3 调用合约

3.3.3.1 调用合约(同步)

`func (r *RPC) InvokeContract(transaction *Transaction) (*TxReceipt, StdError)`

- 说明：调用合约需要初始化调用Transaction。

- 参数【transaction】：调用transaction，transaction中调用合约有参数需要用abi来进行编码。

- 返回【返回值1】：交易回执。回执结构体请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

```go
ABI, _ := abi.JSON(strings.NewReader(abi)
packed, _ := ABI.Pack("add", uint32(1), uint32(2))
transaction := rpc.NewTransaction(key.GetAddress()).Invoke(txDeploy.ContractAddress, packed)
transaction.Sign(key)
txInvoke, _ := rpcAPI.InvokeContract(transaction)
fmt.Println(txInvoke.Ret)
```


**注意**：ABI升级需要将Solidity合约的具体类型和Go的具体类型对应起来，开发应用时需要注意，下面给出了Silidity合约的具体类型和Go类型对应的具体实例：

| Solidity类型                                      | Go类型                                        |
| ------------------------------------------------- | --------------------------------------------- |
| bytes                                             | byte切片                                      |
| bytes0...32, 步长为1                              | byte数组，且byte数组大小需要与bytesXX对应     |
| int8 int16 int32 int64 uint8 uint16 uint32 uint64 | Go类型与Solidity类型相同                      |
| int8...256, uint8...256, 步长为8                  | 不在上述整形类型中的其他整形均需要使用big.Int |
| string                                            | 同Go string                                   |
| address                                           | 同Go common.address                           |

合约：

```javascript
contract TypeCheck {
    event event1(bytes log1, int log2);
    function fun1(bytes data1, bytes32 data2, bytes8 data3) returns (bytes, bytes32, bytes8) {
        return (data1, data2, data3);
    }
    function fun2(int data1, int256 data2, int72 data3, int64 data4, int8 data5) returns (int, int256, int72, int64, int8) {
        return (data1, data2, data3, data4, data5);
    }
    function fun3(uint data1, uint256 data2, uint72 data3, uint64 data4, uint8 data5) returns (uint, uint256, uint72, uint64, uint8) {
        return (data1, data2, data3, data4, data5);
    }
    function fun4(int56 data1, int16 data2, int24 data3, uint56 data4, uint16 data5, uint24 data6) returns (int56, int16, int24, uint56, uint16, uint24) {
        return (data1, data2, data3, data4, data5, data6);
    }
    function fun5(string data1, address data2) returns (string, address) {
        return (data1, data2);
    }
}
```

调用实例：

```go
ABI, _ := abi.JSON(strings.NewReader(abiStr))
// invoke fun1
{
    var data32 [32]byte
    copy(data32[:], "data32")
    var data8 [8]byte
    copy(data8[:], "byte8")
    packed1, _ := ABI.Pack("fun1", []byte("data1"), data32, data8)
    invokeTx1 := NewTransaction(guomiKey.GetAddress()).Invoke(contractAddress, packed1)
    invokeTx1.Sign(guomiKey)
    invokeRe1, _ := rpc.InvokeContract(invokeTx1)
    var p0 []byte
    var p1 [32]byte
    var p2 [8]byte
    testV := []interface{}{&p0, &p1, &p2}
    if err := ABI.UnpackResult(&testV, "fun1", invokeRe1.Ret); err != nil {
        t.Error(err)
        return
    }
    fmt.Println(string(p0), string(p1[:]), string(p2[:]))
}
// invoke fun2
{
    bigInt1 := big.NewInt(-100001)
    bigInt2 := big.NewInt(-1000001)
    bigInt3 := big.NewInt(10000001)
    int1 := int64(-10001)
    int2 := int8(101)
    packed, _ := ABI.Pack("fun2", bigInt1, bigInt2, bigInt3, int1, int2)
    invokeTx := NewTransaction(guomiKey.GetAddress()).Invoke(contractAddress, packed)
    invokeTx.Sign(guomiKey)
    invokeRe, _ := rpc.InvokeContract(invokeTx)
    var p0 interface{}
    var p1 *big.Int
    var p2 *big.Int
    var p3 interface{}
    var p4 int8
    testV := []interface{}{&p0, &p1, &p2, &p3, &p4}
    if err := ABI.UnpackResult(&testV, "fun2", invokeRe.Ret); err != nil {
        t.Error(err)
        return
    }
    fmt.Println(p0, p1.Int64(), p2, p3, p4)
}
// invoke fun3
{
    bigInt1 := big.NewInt(100001)
    bigInt2 := big.NewInt(1000001)
    bigInt3 := big.NewInt(10000001)
    int1 := uint64(10001)
    int2 := uint8(101)
    packed, _ := ABI.Pack("fun3", bigInt1, bigInt2, bigInt3, int1, int2)
    invokeTx := NewTransaction(guomiKey.GetAddress()).Invoke(contractAddress, packed)
    invokeTx.Sign(guomiKey)
    invokeRe, _ := rpc.InvokeContract(invokeTx)
    var p0 interface{}
    var p1 *big.Int
    var p2 *big.Int
    var p3 interface{}
    var p4 uint8
    testV := []interface{}{&p0, &p1, &p2, &p3, &p4}
    if err := ABI.UnpackResult(&testV, "fun3", invokeRe.Ret); err != nil {
        t.Error(err)
        return
    }
    fmt.Println(p0, p1, p2, p3, p4)
}
// invoke fun4
{
    bigInt1 := big.NewInt(-100001)
    a16int := int16(-10001)
    bigInt3 := big.NewInt(10001)
    bigInt4 := big.NewInt(1111111)
    a16uint := uint16(10001)
    bigInt6 := big.NewInt(111111)
    packed, _ := ABI.Pack("fun4", bigInt1, a16int, bigInt3, bigInt4, a16uint, bigInt6)
    invokeTx := NewTransaction(guomiKey.GetAddress()).Invoke(contractAddress, packed)
    invokeTx.Sign(guomiKey)
    invokeRe, _ := rpc.InvokeContract(invokeTx)
    var p0 interface{}
    var p1 int16
    var p2 *big.Int
    var p3 interface{}
    var p4 uint16
    var p5 *big.Int
    testV := []interface{}{&p0, &p1, &p2, &p3, &p4, &p5}
    if err := ABI.UnpackResult(&testV, "fun4", invokeRe.Ret); err != nil {
        t.Error(err)
        return
    }
    fmt.Println(p0, p1, p2, p3, p4, p5)
}
// invoke fun5
{
    address := common.Address{}
    address.SetString("2312321312")
    packed, _ := ABI.Pack("fun5", "data1", address)
    invokeTx := NewTransaction(guomiKey.GetAddress()).Invoke(contractAddress, packed)
    invokeTx.Sign(guomiKey)
    invokeRe, _ := rpc.InvokeContract(invokeTx)
    var p0 string
    var p1 common.Address
    testV := []interface{}{&p0, &p1}
    if err := ABI.UnpackResult(&testV, "fun5", invokeRe.Ret); err != nil {
        t.Error(err)
        return
    }
    fmt.Println(p0, p1)
}
```

3.3.4 智能合约管理

3.3.4.1 同步管理

`func (r *RPC) MaintainContract(transaction *Transaction) (*TxReceipt, StdError)`

- 说明：管理合约需要初始化管理Transaction。

- 参数【transaction】：管理transaction，transaction中不同的opCode代表不同的操作，1代表升级合约，2代表冻结合约，3代表解冻合约。

- 返回【返回值1】：交易回执。回执结构体请看4.1.2。【返回值2】：错误类型请见4.1.5。

- 实例

升级合约

```go
transactionUpdate := rpc.NewTransaction(gmKey.GetAddress()).Maintain(1, originContractAddress, compileUpdate.Bin[0])
transactionUpdate.Sign(gmKey)
receiptUpdate, err := rpcAPI.MaintainContract(transactionUpdate)
if err != nil {
    t.Error(err)
    return
}
fmt.Println(receiptUpdate.ContractAddress)
```

冻结和解冻合约

```go
// freeze contract
transactionFreeze := rpc.NewTransaction(gmKey.GetAddress()).Maintain(2, contractAddress, "")
transactionFreeze.Sign(gmKey)
receiptFreeze, err := rpcAPI.MaintainContract(transactionFreeze)
fmt.Println(receiptFreeze.TxHash)
status, err := rpcAPI.GetContractStatus(contractAddress)
fmt.Println("contract status >>", status)
// unfreeze contract
transactionUnfreeze := rpc.NewTransaction(gmKey.GetAddress()).Maintain(3, contractAddress, "")
transactionUnfreeze.Sign(gmKey)
receiptUnFreeze, err := rpcAPI.MaintainContract(transactionUnfreeze)
fmt.Println(receiptUnFreeze.TxHash)
status, _ = rpcAPI.GetContractStatus(contractAddress)
fmt.Println("contract status >>", status)
```

3.3.5 获取合约状态

`func (r *RPC) GetContractStatus(contractAddress string) (string, StdError)`

- 说明：通过合约的部署地址来查看当前合约的状态。

- 参数【contractAddress】：合约的部署地址。

- 返回【返回值1】：合约状态，返回值为`normal`表示正常，`frozen`表示冻结，`non-contract`表示非合约即普通转账交易。【返回值2】：错误类型请见4.1.5。

- 实例

```go
contractAddress := "0x3ffca734f03458d83d8cff1dc25e49c46feea3bc"
statu, err := rpcAPI.GetContractStatus(contractAddress)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(statu)
```

3.3.6 获取已部署的合约列表

`func (r *RPC) GetDeployedList(address string) ([]string, StdError)`

- 说明：通过账户地址来查看该账户部署的合约列表。

- 参数【address】：账户地址。

- 返回【返回值1】：合约地址列表（切片）。【返回值2】：错误类型请见4.1.5。

- 实例

```go
contracts, err := rpcAPI.GetDeployedList(gmKey.GetAddress())
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(contracts)
```

3.3.7 获取合约字节编码

`func (rpc *RPC) GetCode(contractAddress string) (string, StdError)`

- 说明：获取合约字节编码

- 参数【contractAddress】：合约地址。

- 返回【返回值1】：十六进制字节码。【返回值2】：错误类型请见4.1.5。

- 实例

```go
code, err := rpcAPI.GetCode(contractAddress)
```

3.3.8 获取合约数量

`func (rpc *RPC) GetContractCountByAddr(accountAddress string) (uint64, StdError)`

- 说明：获取合约数量

- 参数【accountAddress】：账户地址。

- 返回【返回值1】：合约数量。【返回值2】：错误类型请见4.1.5。

- 实例

```go
count, err := rpcAPI.GetContractCountByAddr(accountAddress)
```

3.3.9 获取同态加密之后的账户余额以及转账金额

`func (rpc *RPC) EncryptoMessage(balance, amount uint64, invalidHmValue string) (*BalanceAndAmount, StdError)`

- 说明：获取同态加密之后的账户余额以及转账金额

- 参数【balance】：账户未转账之前所有的余额【amount】：要转的金额 【invalidHmValue】：非法同态值

- 返回【返回值1】：包含同态加密之后的账户余额以及同态加密之后的转账金额的结构体。【返回值2】：错误类型请见4.1.5。

- 实例

```go
balanceAndAmount, err := rpcAPI.encryptoMessage(balance, amount, invalidHmValue)
```

3.3.10 获取收款方对所有未验证同态交易的验证结果

`func (rpc *RPC) CheckHmValue(rawValue []uint64, encryValue []string, invalidHmValue string) (*ValidResult, StdError)`

- 说明：获取收款方对所有未验证同态交易的验证结果

- 参数【rawValue】：经过收款方椭圆曲线私钥解密之后的所有未验证转账金额【encryValue】：所有未验证交易中的同态加密过的转账金额【invalidHmValue】：收款方当前所有非法转账金额同态值之

- 返回【返回值1】：验证结果。【返回值2】：错误类型请见4.1.5。

- 实例

```go
validResult, err := rpcAPI.checkHmValue(rawValue, encryValue, invalidHmValue)
```

3.3.11 查询合约部署者

`func (rpc *RPC) GetCreator(contractAddress string) (string, StdError)`

- 说明：查询合约部署者

- 参数【contractAddress】：合约地址

- 返回【返回值1】：合约部署者的地址。【返回值2】：错误类型请见4.1.5。

- 实例

```go
accountAddr, err := rpcAPI.GetCreator(contractAddress)
```

3.3.12 查询合约部署时间

`func (rpc *RPC) GetCreateTime(contractAddress string) (string, StdError)`

- 说明：查询合约部署者

- 参数【contractAddress】：合约地址

- 返回【返回值1】：合约部署的日期时间。【返回值2】：错误类型请见4.1.5。

- 实例

```go
dataTime, err := rpcAPI.GetCreateTime(contractAddress)
```

### 3.4 Block相关接口

3.4.1 取得最新区块信息

`func (r *RPC) GetLatestBlock() (*Block, StdError)`

- 说明：用来取得最新的区块信息

- 返回【返回值1】：最新区块信息，包含区块的一些元信息以及其中的交易信息。Block结构体详情见3.4。【返回值2】：错误类型请见4.1.5。

- 实例：

```go
block, err := hrpc.GetLatestBlock()
```

3.4.2 取得指定区块列表

`func (r *RPC) GetBlocks(from, to uint64, isPlain bool) ([]*Block, StdError)`

- 说明：取得指定区块号范围内所有区块信息。

- 参数【from】：起始区块号，要求必须为非0正整数。【to】：结束区块号，要求必须为大于from的正整数。【isPlain】：如果该值为false，返回结果将包括区内的交易信息。若为true，则不包括。

- 返回【返回值1】：区块号为[from, to]区间内的区块数组。Block结构体详情见3.4。若未查询到返回值2将为非nil。【返回值2】：错误类型请见4.1.5。

- 实例

```go
latestBlock, err := hrpc.GetLatestBlock()
if err != nil {
    t.Error(err.String())
    return
}
blocks, err := hrpc.GetBlocks(latestBlock.Number-1, latestBlock.Number, true)
if err != nil {
    t.Error(err.String())
    return
}
fmt.Println(blocks)
```

3.4.3 根据区块哈希查询区块信息

`func (r *RPC) GetBlockByHash(blockHash string, isPlain bool) (*Block, StdError)`

- 说明：根据区块的哈希取得区块哈希。

- 参数【blockHash】：区块哈希，应该以0x开头。【isPlain】：如果该值为false，返回结果将包括区内的交易信息。若为true，则不包括。

- 返回【返回值1】：区块信息。Block结构体详情见3.4。若未查询到返回值2将为非nil。【返回值2】：错误类型请见4.1.5。

- 实例

```go
latestBlock, err := hrpc.GetLatestBlock()
if err != nil {
    t.Error(err.String())
    return
}
block, err := hrpc.GetBlockByHash(latestBlock.Hash, true)
if err != nil {
    t.Error(err.String())
    return
}
```

3.4.4 根据区块哈希列表批量查询区块信息

`func (r *RPC) GetBatchBlocksByHash(blockHashes []string, isPlain bool) ([]*Block, StdError)`

- 说明：根据区块哈希列表批量查询区块信息

- 参数【blockHashes】：区块哈希切片，哈希应该以0x开头。【isPlain】：如果该值为false，返回结果将包括区内的交易信息。若为true，则不包括。

- 返回【返回值1】：区块哈希切片对应的区块列表。Block结构体详情见3.4。若未查询到返回值2将为非nil。【返回值2】：略

- 实例

```go
latestBlock, err := hrpc.GetLatestBlock()
if err != nil {
    t.Error(err.String())
    return
}
blocks, err := hrpc.GetBatchBlocksByHash([]string{latestBlock.Hash}, true)
```

3.4.5 根据区块号查询区块信息

`func (r *RPC) GetBlockByNumber(blockNum interface{}, isPlain bool) (*Block, StdError)`

- 说明：根据区块号查询区块信息

- 参数【blockNum】：区块号，`latest`表示最新区块【isPlain】：如果该值为false，返回结果将包括区内的交易信息。若为true，则不包括。

- 返回【返回值1】：该区块号对应的区块信息。Block结构体详情见3.4。若未查询到返回值2将为非nil。【返回值2】：错误类型请见4.1.5。

- 实例

```go
latestBlock, err := hrpc.GetLatestBlock()
if err != nil {
    t.Error(err.String())
    return
}
block, err := hrpc.GetBlockByNumber(latestBlock.Number, true)
```

3.4.6 根据区块号列表查询区块信息

`func (r *RPC) GetBatchBlocksByNumber(blockNums []uint64, isPlain bool) ([]*Block, StdError)`

- 说明：根据区块号列表查询区块信息

- 参数【blockNums】：区块号切片。【isPlain】：如果该值为false，返回结果将包括区内的交易信息。若为true，则不包括。

- 返回【返回值1】：区块号切片对应的区块信息。Block结构体详情见3.4。若未查询到返回值2将为非nil。【返回值2】：错误类型请见4.1.5。

- 实例

```go
latestBlock, err := hrpc.GetLatestBlock()
if err != nil {
    t.Error(err.String())
    return
}
blocks, err := hrpc.GetBatchBlocksByNumber([]uint64{latestBlock.Number}, true)
```

3.4.7 获取区间区块生成平均速度

`func (r *RPC) GetAvgGenTimeByBlockNum(from, to uint64) (int64, StdError)`

- 说明：计算[from, to]区间内的区块生成速度。

- 参数【from】：起始区块号。【to】：终止区块号。

- 返回【返回值1】：平均区块生成时间，单位（**ms**）。

- 实例

```go
block, err := hrpc.GetLatestBlock()
if err != nil {
    t.Error(err.String())
    return
}
avgTime, err := hrpc.GetAvgGenTimeByBlockNum(block.Number-2, block.Number)
```

3.4.8 查询指定时间区间内的区块数量

`func (r *RPC) GetBlocksByTime(startTime, endTime uint64) (*BlockInterval, StdError)`

- 说明：查询指定时间区间内的区块数量。

- 参数【startTime】：起始时间戳，单位（**ns**）。【endTime】：结束时间戳，单位（**ns**）。

- 返回【返回值1】：包含区块总数、起始区块号、结束区块号信息。BlockInterval结构体信息请看4.1.1节。【返回值2】：错误类型请见4.1.5。

- 实例

```go
blockInterval, err := hrpc.GetBlocksByTime(1, 1778959217012956575)
if err != nil {
    t.Error(err.String())
    return
}
```

3.4.9 查询指定时间区间内的区块和交易生成速度

`func (r *RPC) QueryTPS(startTime, endTime uint64) (*TPSInfo, StdError)`

- 说明：查询指定时间区间内的区块和交易生成速度。

- 参数【startTime】：起始时间戳，单位（**ns**）。【endTime】：结束时间戳，单位（**ns**）。

- 返回【返回值1】：包含该时间区间内的区块总数，每秒生成区块数（区块/s），每秒交易生成数（交易/s）。TPSInfo结构体信息请看4.1.1节。【返回值2】：错误类型请见4.1.5。

- 实例

```go
tpsInfo, err := hrpc.QueryTPS(1, 1778959217012956575)
if err != nil {
    t.Error(err.String())
    return
}
```

3.4.10 查询创世区块号

`func (r *RPC) GetGenesisBlock() (string, StdError)`

- 说明：查询创世区块的区块号。

- 返回【返回值1】：创世区块号，16进制字符串。【返回值2】：错误类型请见4.1.5。

- 实例

```go
blkNum, err := rpc.GetGenesisBlock()
```

3.4.11 查询最新区块号

`func (r *RPC) GetChainHeight() (string, StdError)`

- 说明：查询最新的区块号。

- 返回【返回值1】：最新区块号，16进制字符串。【返回值2】：错误类型请见4.1.5。

- 实例

```go
blkNum, err := rpc.GetChainHeight()
```

### 3.5 Node相关接口

节点信息类型详见4.1.3

3.5.1 获取节点信息

`func (r *RPC) GetNodes() ([]NodeInfo, StdError)`

- 说明：获取部署的所有节点的信息。

- 返回【返回值1】：节点信息（切片）。Node相关结构体信息请看4.1.3。【返回值2】：错误类型请见4.1.5。

- 实例

```go
nodes, err := rpcAPI.GetNodes()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(len(nodes))
```

3.5.2 获取节点hash

3.5.2.1 获取随机节点的hash

`func (r *RPC) GetNodeHash() (string, StdError)`

- 说明：随机从配置的节点中获取一个节点的hash。

- 返回​    【返回值1】：节点hash。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
hash, err := rpcAPI.GetNodeHash()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(hash)
```

3.5.2.2 获取指定节点的hash

`func (r *RPC) GetNodeHashByID(id int) (string, StdError)`

- 说明：随机从配置的节点中获取一个节点的hash。

- 参数​    【id】：节点序号，按配置文件中的位置

- 返回​    【返回值1】：节点hash。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
hash, err := rpcAPI.GetNodeHashByID(1)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(hash)
```

3.5.3 删除VP节点

`func (r *RPC) DeleteNodeVP(hash string) (bool, StdError)`

- 说明：根据节点hash来删除VP节点。

- 参数​    【hash】：节点hash。

- 返回​    【返回值1】：是否删除成功。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
success, err := rpcAPI.DeleteNodeVP(hash)
```

3.5.4 删除NVP节点

`func (r *RPC) DeleteNodeNVP(hash string) (bool, StdError)`

- 说明：根据节点hash来删除NVP节点。

- 参数​    【hash】：节点hash。

- 返回​    【返回值1】：是否删除成功。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
success, err := rpcAPI.DeleteNodeNVP(hash)
```

3.5.5 获取节点的状态信息

`func (r *RPC) GetNodeStates() ([]NodeStateInfo, StdError)`

- 说明：获取当前节点所连接的所有对端节点的状态信息。

- 返回​    【返回值1】：节点状态信息。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
infos, err := rpcAPI.GetNodeStates()
```

### 3.6 数据归档相关接口

3.6.1 预约某区块高度归档

`func (r *RPC) Snapshot(blockHeight interface{}) (string, StdError)`

- 说明：预约在blockHeight高度进行归档操作。

- 参数【blockHeight】：要求**大于等于**目前区块链上的高度，`latest`表示当前区块，即立即归档

- 返回【返回值1】：快照标号，用来查询是否已经归档。

- 实例

```go
res, err := hrpc.Snapshot(1)
```

3.6.2 查询快照制作结果

`func (r *RPC) QuerySnapshotExist(filterID string) (bool, StdError)`

- 说明：查询快照制作结果

- 参数【filterID】：快照标号。

- 返回【返回值1】：标识快照是否存在。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.QuerySnapshotExist("0x5d86cce7e537cd0e0346468889801196")
```

3.6.3 快照检查

`func (r *RPC) CheckSnapshot(filterID string) (bool, StdError)`

- 说明：快照检查，检查快照内容是否正确。

- 参数【filterID】：快照标号。

- 返回【返回值1】：快照内容是否正确。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.CheckSnapshot("0x5d86cce7e537cd0e0346468889801196")
```

3.6.4 数据归档

`func (r *RPC) Archive(filterID string, sync bool) (bool, StdError)`

- 说明：数据归档

- 参数【filterID】：快照标号。【sync】：是否同步执行。

- 返回【返回值1】：数据归档是否成功。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.Archive("0x5d86cce7e537cd0e0346468889801196", false)
```

3.6.5 查询数据归档结果

`func (r *RPC) QueryArchiveExist(filterID string) (bool, StdError)`

- 说明：查询数据归档结果。

- 参数【filterID】：快照标号。

- 返回【返回值1】：快照是否已归档。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.QueryArchiveExist("0x5d86cce7e537cd0e0346468889801196")
```

3.6.6 删除快照

`func (r *RPC) DeleteSnapshot(filterID string) (bool, StdError)`

- 说明：删除生成的快照。

- 参数【filterID】：快照标号。

- 返回【返回值1】：是否删除成功。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.DeleteSnapshot("0x5d86cce7e537cd0e0346468889801196")
```

3.6.7 获取快照列表

`func (r *RPC) ListSnapshot() (Manifests, StdError)`

- 说明：获取所有已完成的快照信息。

- 返回【返回值1】：快照信息数组，详见4.1.7。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.ListSnapshot()
```

3.6.8 查询快照

`func (r *RPC) ReadSnapshot(filterID string) (Manifest, StdError)`

- 说明：查询指定快照的详细信息，如果快照已完成，则返回对应信息，否则返回错误码-32013的error。

- 参数【filterID】：快照标号。

- 返回【返回值1】：快照信息，详见4.1.7。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.ReadSnapshot("0x5d86cce7e537cd0e0346468889801196")
```

3.6.9 恢复归档数据

`func (r *RPC) Restore(filterID string, sync bool) (bool, StdError)`

- 说明：执行恢复归档数据的快照id必须是一个已经完成的快照并且做过数据归档，否则会返回错误-32013，表示快照不存在或者这是个异常操作。

- 参数【filterID】：快照标号。【sync】：是否同步执行

- 返回【返回值1】：恢复是否成功。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.Restore("0x5d86cce7e537cd0e0346468889801196"， true)
```

3.6.10 恢复所有归档数据

`func (r *RPC) RestoreAll(sync bool) (bool, StdError)`

- 说明：恢复所有的归档数据。

- 参数【sync】：是否同步执行

- 返回【返回值1】：恢复是否成功。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.RestoreAll(true)
```

3.6.11 查询所有待完成的快照请求

`func (r *RPC) Pending() ([]SnapshotEvent, StdError)`

- 说明：返回一个预快照列表。

- 返回【返回值1】：预快照的快照ID和发生快照的区块号，详见4.1.7。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.Pending()
```

3.6.12 查询最近一次归档的状态

`func (rpc *RPC) QueryLatestArchive() (*ArchiveResult, StdError)`

- 说明：返回一个归档结果信息。

- 返回【返回值1】：归档的进度及失败原因，详见4.1.7。【返回值2】：错误类型请见4.1.5。

- 实例

```go
res, err := hrpc.QueryLatestArchive()
```

### 3.7 返回值解析

3.7.1 Solidity合约返回值解析

`func (abi *ABI) UnpackResult(v interface{}, name, data string) (err error)`

- 说明：使用ABI来解码合约调用返回的返回结果。

- 参数​    【v】：解码后的对象。​    【name】：解码方法名。​    【data】：合约返回结果的源数据（TxReceipt.Ret）。

- 返回​    【返回值1】：解码错误

- 实例

单返回值

```go
// 调用的合约方法
function getSum() returns(uint32) {return sum;}
// 单返回值
var p uint32
if sysErr := ABI.UnpackResult(&p, "getSum", receipt.Ret); sysErr != nil {
    fmt.Println(sysErr)
    return
}
fmt.Println(p)
```

多返回值

```go
// 调用的合约方法
function getMul() returns(bytes, int64, address) {
    return ("hello", 12, msg.sender);
}
// 多返回值 需要将Solidity的类型同Go中对应
var p0 []byte
var p1 int64
var p2 common.Address
testV := []interface{}{&p0, &p1, &p2}
fmt.Println(reflect.TypeOf(testV))
if sysErr := ABI.UnpackResult(&testV, "getMul", receipt.Ret); sysErr != nil {
    fmt.Println(sysErr)
    return
}
fmt.Println(string(p0), p1, p2.Hex())
```

3.7.2 Java合约返回值解析

`func DecodeJavaResult(ret string) string`

- 说明：解码Java合约返回的结果。

- 参数​    【ret】：合约返回结果。

- 返回​    【返回值1】：解码后的结果

- 实例

```go
fmt.Println(java.DecodeJavaResult(txReceipt.Ret))
```

3.7.3 HVM合约返回值解析

`func DecodeJavaResult(ret string) string`

- 说明：解码HVM合约返回的结果。

- 参数​    【ret】：合约返回结果。

- 返回​    【返回值1】：解码后的结果

- 实例

```go
fmt.Println(java.DecodeJavaResult(txReceipt.Ret))
```

3.7.4 静态类型解析说明

`func ByteArrayToString(v interface{}) (string, error)`

- 说明：在解析Solidity合约返回值或Log值时，需要解析的类型为静态类型，例如bytes32，我们解码时需要对应的Go类型（[32]byte），若实际有效数据不够实际长度，会导致多余的0值产生乱码，需要调用`abi`包下的静态类型解码。

- 参数​    【v】：解码数据，支持字节数组。

- 返回​    【返回值1】：解码后字符串。​    【返回值2】：错误类型请见4.1.5。

- 注意：Solidity中`address`实际为固定20字节的数组。

- 实例

```go
var result [32]byte
if err := ABI.UnpackResult(&result, "getText", data); err != nil {
    fmt.Println(err)
    return
}
resultReal, _ := abi.ByteArrayToString(result)
fmt.Println(resultReal)
```

### 3.8 Log值解析

由平台发来的合约生成的日志是被编码过的，需要使用GoSDK提供的解码工具来解码还原。

3.8.1 solidity合约

`func (abi ABI) UnpackLog(v interface{}, name string, data string, topics []string) (err error)`

- 说明：根据**solidity合约**abi将log输入填充到v中。**不支持解码indexed修饰的动态类型（如bytes、string等）。**

- 参数【v】：结构体指针，映射关系可以通过变量名来设置，也可以通过abi tag来强制设置，要求字段类型需要和合约中声明的类型兼容。【name】：event名称。【data】：TxLog.Data字段。TxLog结构体请看4.1.2。【topics】：TxLog.Topics字段。在event中有indexed参数时使用。

- 返回【返回值1】：略。

- 实例

solidity合约：

```solidity
contract Accumulator {
    event sayHello(int64 addr1, bytes8 indexed msg);
    uint32 sum = 0;
    bytes32 hello = "hello world";
    function Accumulator(uint32 sum1, bytes32 hello1) {
        sum = sum1;
        hello = hello1;
    }
    function getHello() constant returns(bytes32) {
        sayHello(1, "test");
        saySum("sum", sum);
        return hello;
    }
}
```

GoSDK：

```go
test := struct {
    // 如果该字段名为Addr1，那么不需要abi tag也能进行映射
    Addr int64   `abi:"addr1"`
    Msg1 [8]byte `abi:"msg"`
}{}
sysErr := ABI.UnpackLog(&test, "sayHello", receipt1.Log[0].Data, receipt1.Log[0].Topics)
if sysErr != nil {
    t.Error(sysErr)
    return
}
// 解码
msg, sysErr := abi.ByteArrayToString(test.Msg1)
```

3.8.2 java合约

`func DecodeJavaLog(data string) (string, error)`

- 说明：用来解码**java合约**返回的log值

- 参数【data】：TxLog的Data字段。TxLog结构体请看4.1.2。

- 返回【返回值1】：解码后的json字符串【返回值2】：略。

- 实例

java合约：

```java
public ExecuteResult testPostEvent(List<String> var1) {
    this.logger.info(var1);
    for(int var2 = 0; var2 < 10; ++var2) {
        Event var3 = new Event("event" + var2);
        var3.addTopic("simulate_bank");
        var3.addTopic("test");
        var3.put("attr1", "value1");
        var3.put("attr2", "value2");
        var3.put("attr3", "value3");
        this.ledger.post(var3);
    }
    return this.result(true);
}
```

GoSDK：

```go
res, err := DecodeJavaLog(txReceipt.Log[0].Data)
```

3.8.2 HVM合约

HVM合约Log值解析过程与java合约相同

### 3.9 账户相关接口

目前根据加密方式不同一共有10种类型账户，非国密账户和国密账户各5种。

非国密账户有：

`ECKDF2 = "0x01"` , `ECDES = "0x02"`, `ECRAW = "0x03"`, `ECAES = "0x04"`, `EC3DES = "0x05"`

国密账户有：

`SMSM4 = "0x11"`, `SMDES = "0x12"`, `SMRAW = "0x13"`, `SMAES = "0x14"`, `SM3DES = "0x15"`

3.9.1 国密

3.9.1.1 创建账户JSON

`func NewAccountSm2(password string) (string, error)`

- 说明：创建sm2国密算法的accountJSON，可以和JavaSDK兼容的账户。

- 参数​    【passwrod】：是否加密私钥，不为空则加密，为空则不加密。

- 返回​    【返回1】：账户JSON串。​    【返回2】：错误类型请见4.1.5。

- 实例

不加密私钥

```go
accountJson, err := account.NewAccountSm2("")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(accountJson)
```

加密私钥

```go
accountJson, err := account.NewAccountSm2("123")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(accountJson)
```

`func NewAccountJson(acType, password string) (string, error)`

- 说明：可根据传入的账户加密类型创建不同国密算法的accountJSON，可以和JavaSDK兼容的账户。

- 参数   【acType】：账户加密类型。"0x11"为SMSM4；"0x12"为SMDES；"0x13"为SMRAW；"0x14"为SMAES。​    【passwrod】：是否加密私钥，不为空则加密，为空则不加密。

- 返回​    【返回1】：账户JSON串。​    【返回2】：错误类型请见4.1.5。

- 实例

不加密私钥

```go
accountJson, err := account.NewAccountJson(SMRAW,"")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(accountJson)
```

SM4加密私钥

```go
accountJson, err := account.NewAccountJson(SMSM4,"123")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(accountJson)
```

3.9.1.2 获取账户（by JSON）

`func NewAccountSm2FromAccountJSON(accountjson, password string) (*gm.Key, error)`

- 说明：从accountJson中转化为GoSDK可使用的国密key。

- 参数​    【accountjson】：SDK统一的账户JSON串。​    【password】：JSON串的密码，对有密码的私钥进行解码。

- 返回​    【返回值1】：GoSDK使用的国密key。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
key, err := account.NewAccountSm2FromAccountJSON(accountJson,  "") // 无密
key, err := account.NewAccountSm2FromAccountJSON(accountJson,  "123") // 有密
```

`func GenKeyFromAccountJson(accountJson, password string) (key interface{}, err error)`

- 说明：从accountJson中转化为GoSDK可使用的国密key。

- 参数​    【accountjson】：SDK统一的账户JSON串。​    【password】：JSON串的密码，对有密码的私钥进行解码。

- 返回​    【返回值1】：GoSDK使用的国密key。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
key, err := account.GenKeyFromAccountJson(accountJson,  "") // 无密
key, err := account.GenKeyFromAccountJson(accountJson,  "123") // 有密
```

3.9.1.3 获取账户（by privateKey）

`func NewAccountSm2FromPriv(priv string) (*gm.Key, error)`

- 说明：从私钥字符串获取GoSDK可使用的国密key。

- 参数​    【priv】：私钥字符串。

- 返回​    【返回值1】：GoSDK使用的国密key。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
key, err := account.NewAccountSm2FromPriv(priv)
```

3.9.1.4 获取账户地址

`func (key *Key) GetAddress() string`

- 说明：从国密key中获取账户的地址。

- 返回​    【返回值1】：账户地址

- 实例

```go
key, _ := account.NewAccountSm2FromPriv(priv)
address := key.GetAddress()
```

3.9.2 非国密

3.9.2.1 创建账户JSON

`func NewAccount(password string) (string, error)`

- 说明：创建ecdsa算法的accountJSON，可以和JavaSDK兼容的账户。

- 参数​    【passwrod】：是否加密私钥，不为空则加密，为空则不加密。

- 返回​    【返回1】：账户JSON串。​    【返回2】：错误类型请见4.1.5。

- 实例

```go
accountJson, err := account.NewAccount("") // 无密账户
accountJson, err = account.NewAccount("123") // 加密账户
```

`func NewAccountJson(acType, password string) (string, error)`

- 说明：可根据传入的账户加密类型创建不同非国密算法的accountJSON，可以和JavaSDK兼容的账户。

- 参数   【acType】：账户加密类型。"0x01"为ECKDF2（目前不支持）；"0x02"为ECDES；"0x03"为ECRAW；"0x04"为ECAES。​    【passwrod】：是否加密私钥，不为空则加密，为空则不加密。

- 返回​    【返回1】：账户JSON串。​    【返回2】：错误类型请见4.1.5。

- 实例

不加密私钥

```go
accountJson, err := account.NewAccountJson(ECRAW,"")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(accountJson)
```

AES加密私钥

```go
accountJson, err := account.NewAccountJson(ECAES,"123")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(accountJson)
```

3.9.2.2 获取账户（by JSON）

`func NewAccountFromAccountJSON(accountjson, password string) (*ecdsa.Key, error)`

- 说明：从accountJson中转化为GoSDK可使用的ecdsa key。

- 参数​    【accountjson】：SDK统一的账户JSON串。​    【password】：JSON串的密码，对有密码的私钥进行解码。

- 返回​    【返回值1】：GoSDK使用的ecdsa key。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
key, err := account.NewAccountFromAccountJSON(accountJson,  "") // 无密账户
key, err := account.NewAccountFromAccountJSON(accountJson,  "123") // 加密账户
```

`func GenKeyFromAccountJson(accountJson, password string) (key interface{}, err error)`

- 说明：从accountJson中转化为GoSDK可使用的非国密key。

- 参数​    【accountjson】：SDK统一的账户JSON串。​    【password】：JSON串的密码，对有密码的私钥进行解码。

- 返回​    【返回值1】：GoSDK使用的非国密key。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
key, err := account.GenKeyFromAccountJson(accountJson,  "") // 无密
key, err := account.GenKeyFromAccountJson(accountJson,  "123") // 有密
```

3.9.2.3 获取账户（by privateKey）

`func NewAccountFromPriv(priv string) *ecdsa.Key`

- 说明：从私钥字符串获取GoSDK可使用的ecdsa key。

- 参数​    【priv】：私钥字符串。

- 返回​    【返回值1】：GoSDK使用的ecdsa key。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
key := account.NewAccountFromPriv(priv)
```

3.9.2.4 获取账户地址

`func (key *Key) GetAddress() string`

- 说明：从ecdsa key中获取账户的地址。

- 返回​    【返回值1】：账户地址

- 实例

```go
key := account.NewAccountFromPriv(priv)
address := key.GetAddress()
```

3.9.3 获取账户余额

`func (r *RPC) GetBalance(account string) (string, StdError)`

- 说明：根据账户地址获取账户余额。

- 参数​    【account】：账户地址。

- 返回​    【返回值1】：余额。​    【返回值2】：错误类型请见4.1.5。

- 实例

```go
balance, err := rpcAPI.GetBalance(key.GetAddress())
```

### 3.10 WebSocket相关接口

WebSocket系列接口是用来实现事件订阅的相关功能。用户通过GoSDK向Hyperchain订阅自己感兴趣的事件，然后当事件发生时平台会向GoSDK主动推送。

3.10.1 获取WebSocket客户端

`func (r *RPC) GetWebSocketClient() *WebSocketClient`

- 说明：获取WebSokcet客户端，所有的web socket的接口都是基于web socket客户端。

- 返回【返回值1】：web socket客户端。

- 实例

```go
wsCli := wsRPC.GetWebSocketClient()
```

3.10.2 向某个节点订阅事件

`func (wscli *WebSocketClient) Subscribe(nodeIndex int, filter EventFilter, eventHandler WsEventHandler) (SubscriptionID, StdError)`

- 说明：用来向hyperchain中的某个节点订阅某个事件。

- 参数【nodeIndex】：从1开始。代表节点编号，与hpc.toml配置文件中jsonRPC下的nodes配置一致，1号表示向1号节点订阅，2号表示向2号节点订阅，以此类推。【filter】：事件的过滤条件，要求实现EventFilter接口。所有过滤器类型请看3.10.2.1~3.10.2.3。【eventHandler】：用户自定义回调函数，需要实现WsEventHandler接口，当事件发生时会自动触发对应回调函数：在订阅成功时会触发`OnSubscribe()`，在取消订阅成功时会触发`OnUnSubscribe()`，在接收到事件推送时会触发`OnMessage([]byte)`，在连接关闭时会触发`OnClose()`（若已经取消订阅，那么将不会再触发关闭连接回调）。

- 返回【返回值1】：订阅ID，可以用来取消订阅。【返回值2】：错误类型请见4.1.5。

- 实例

3.10.2.1 创建Block事件过滤器

`func NewBlockEventFilter() *BlockEventFilter`

- 说明：创建一个Block事件过滤器

`func (bf *BlockEventFilter) SetBlockInfo(b bool)`

- 说明：设置是否返回区块详细信息。

- 参数【b】：为true表示通知信息中包括最新区块详细信息，为false表示通 知信息中只返回最新区块哈希。

- 过滤器设置实例

```go
bf := NewBlockEventFilter().
    SetBlockInfo(true)
```

- 返回实例

```json
# 通知
<< {
    "jsonrpc":"2.0",
    "namespace":"global",
    "result":{
        "event":"block",
        "subscription":"0xd0f98f277f33b1125b2a5f5ac2cbc5b",
        "data":{
            "version":"1.3",
            "number":"0x1",
            "hash":"0xcc6d22af42e3cd0241f9f1a1a166764e642de2c6c2c86a5b4222
            68c2e91f16e6",
            "parentHash":"0x0000000000000000000000000000000000000000000000
            000000000000000000",
            "writeTime":1499157819652020682,
            "avgTime":"0x11",
            "txcounts":"0x1",
            "merkleRoot":"0xe6d8cd299eed0bae53a2439d3e5a1fb85a47edc50161fc
            53774266d946a23b15"
        }
    }
}
```

3.10.2.2 创建系统状态过滤器

**事件返回**：

- module ：状态信息所属模块。

- subType ：状态信息所属类型。

- errorCode ：状态信息的标志码。

- message ：具体状态信息。

- date ：状态信息抛出日期。

***

`func NewSystemStatusFilter() *SystemStatusFilter`

- 说明：创建一个系统状态事件过滤器

`func (ssf *SystemStatusFilter) AddModules(modules ...string) *SystemStatusFilter`

- 说明：设置要订阅哪些模块的状态信息。

- 参数【modules】：表示要订阅哪些模块的状态信息，若为空，则表示订阅所有模 块。比如：p2p、consensus、executor等。

- 返回【返回值1】：可以继续链式调用。

`func (ssf *SystemStatusFilter) AddModulesExclude(modulesExclude ...string) *SystemStatusFilter`

- 说明：设置要排除哪些模块的状态信息。

- 参数【modulesExclude】：表示要排除哪些模块的状态信息，若为空，则表示不排除。

- 返回【返回值1】：可以继续链式调用。

`func (ssf *SystemStatusFilter) AddSubtypes(subtypes ...string) *SystemStatusFilter`

- 说明：设置要订阅模块下面的哪一类状态信息。

- 参数【subtypes】：表示要订阅模块下面的哪一类状态信息，若为空，则表示订阅 所有类型。比如：viewchange等。

- 返回【返回值1】：可以继续链式调用。

`func (ssf *SystemStatusFilter) AddSubtypesExclude(subtypesExclude ...string) *SystemStatusFilter`

- 说明：设置要排除模块下面的哪一类状态信息。

- 参数【subtypesExclude】：表示要排除模块下面的哪一类状态信息，若为空， 则表示不排除。

- 返回【返回值1】：可以继续链式调用。

`func (ssf *SystemStatusFilter) AddErrorCode(errorCodes ...string) *SystemStatusFilter`

- 说明：设置要订阅指定的具体哪一条状态信息。

- 参数【errorCodes】：要订阅指定的具体哪一条状态信息，若为空，则表示订阅所有状态信息。

- 返回【返回值1】：可以继续链式调用。

`func (ssf *SystemStatusFilter) AddErrorCodeExclude(errorCodesExclude ...string) *SystemStatusFilter`

- 说明：设置要排除指定的具体哪一条状态信息。

- 参数【errorCodesExclude】：表示要排除指定的具体哪一条状态信息，若为空，则表示不排除。

- 返回【返回值1】：可以继续链式调用。

- 过滤器设置实例

```go
sysf := NewSystemStatusFilter().
        AddModules("p2p").
        AddSubtypes("viewchange")
```

- 返回实例

```json
{
"jsonrpc":"2.0",
"namespace":"global",
    "result":{
        "event":"systemStatus",
        "subscription":"0x81291fec5090f8c214d5e947ef6d3848",
        "data":{
            "module":"executor",
            "status":false,
            "subType":"viewchange",
            "errorCode":-1,
            "message":"required viewchange to 1",
            "date":"2018-03-29T16:56:38.420469472+08:00"
        }
    }
}
```

3.10.2.3 创建日志过滤器

可以用来监听平台执行合约过程中产生的event。

`func NewLogsFilter() *LogsFilter`

- 说明：用来创建日志过滤器。

`func (lf *LogsFilter) SetFromBlock(from string) *LogsFilter`

- 说明：用来设置起始区块号。

- 参数【from】：起始区块号，为空则默认为没有区块下限限制。

`func (lf *LogsFilter) SetToBlock(to uint64) *LogsFilter`

- 说明：用来设置终止区块号。

- 参数【to】：终止区块号，为空则默认为没有区块上限限制。

`func (lf *LogsFilter) AddAddress(addresses ...string) *LogsFilter`

- 说明：用来设置地址列表。

- 参数【addresses】：地址列表，若为空，则默认接收所有来自任意一个合约的log事件；若不为空，则只接收来自地址列表中的合约的log事件。

`func (lf *LogsFilter) SetTopic(pos int, topics ...common.Hash) *LogsFilter`

- 说明：用来设置根据topic筛选event的条件。pos之间是**且**的关系，topics内部为**或**的关系。即若调用SetTopic(0, A, B)和SetTopic(1, C, D)，那么代表订阅topic数组中0号位置为A或B且1号位置为C或D的事件。

- 参数【pos】：代表参数topics将用作过滤topics第pos个位置。例如SetTopic(0, topics)，代表当产生event的topics数组0号位置为参数topics指定的数组中之一时平台将会推送。【topics】：topic某个位置上的过滤条件。注意这里的topic应该设置的是event的ID，可以通过abi结构体得到。

- 实例

```go
// 订阅cAddress合约产生的event中0号topic为getHello的
logf := NewLogsFilter().
    AddAddress(cAddress).
    SetTopic(0, ABI.Events["getHello"].Id())
```

- 返回实例

```json
<< {
"jsonrpc":"2.0",
"namespace":"global",
    "result":{
        "event":"logs",
        "subscription":"0xdd061810765e621fd93291a0146b2e10",
            "data":[{
                "address":"0x313bbf563991dc4c1be9d98a058a26108adfcf81",
                "topics":["0x24abdb5865df5079dcc5ac590ff6f01d5c16edbc5fab4e195
                d9febd1114503da"],
                "data":"000000000000000000000000000000000000000000000000000000
                0000000064",
                "blockNumber":4,
                "blockHash":"0xee93a66e170f2b20689cc05df27e290613da411c42a7bdf
                a951481c08fdefb16",
                "txHash":"0xa676673a23f33a95a1a5960849ad780c5048dff76df961e9f7
                8329b201670ae2",
                "txIndex":0,
                "index":0
            }]
    }
}
```

3.10.3 取消某个订阅

`func (wscli *WebSocketClient) UnSubscribe(id SubscriptionID) StdError`

- 说明：使用订阅ID取消某个订阅。

- 参数【id】：订阅ID。

- 返回【返回值1】：若为nil，代表订阅成功。

- 实例

```go
wsCli.UnSubscribe(subID)
```

3.10.4 关闭与某个节点的WebSocket连接

`func (wscli *WebSocketClient) CloseConn(nodeIndex int) StdError`

- 说明：关闭与某个节点的web socket连接。该条连接上的所有订阅均不再收到消息。

- 参数【nodeIndex】：从1开始。代表节点编号，与hpc.toml配置文件中jsonRPC下的nodes配置一致，1号表示关闭与1号节点的连接，2号表示关闭与2号节点的连接，以此类推。

- 返回【返回值1】：若为nil，代表关闭成功。

- 实例

```go
wsCli.CloseConn(1)
```

3.10.5 获取某个节点上的所有订阅

`func (wscli *WebSocketClient) GetAllSubscription(nodeIndex int) ([]Subscription, StdError)`

- 说明：获取某个节点上的所有订阅。

- 参数【nodeIndex】：从1开始。代表节点编号，与hpc.toml配置文件中jsonRPC下的nodes配置一致，1号表示获得1号节点上所有订阅，2号表示获2号节点上所有订阅，以此类推。

- 返回【返回值1】：订阅数组（取消订阅的不再出现）。包含事件类型和订阅ID。web socket相关结构体请看4.1.4。【返回值2】：若nodeIndex下标越界将会报错。

- 实例

```go
subs, err := wsCli.GetAllSubscription(1)
```

### 3.11 工具方法

3.11.1 读入文件为string

`func ReadFileAsString(path string) (string, error)`

- 说明：将文件内容读取为字符串，可以用于读取合约内容。

- 参数​ 【path】：文件路径，相对路径或绝对路径。

- 返回​ 【返回值1】：文件内容字符串。 【返回值2】：略。

- 实例

```go
contract, err := common.ReadFileAsString(path)
```

3.11.2 编码Solidity合约调用方法

`func (abi ABI) Pack(name string, args ...interface{}) ([]byte, error)`

- 说明：将Solidity合约的调用方法和参数编码为payload。

- 参数​ 【name】：方法名。​ 【args】：调用参数。

- 返回​ 【返回值1】：编码后的字节切片。​ 【返回值2】：略。

- 实例

```go
data, err := abi.Pack("getSum", uint32(1), uint32(2))
```

3.11.3 读取Java合约

`func ReadJavaContract(path string, params ...string) (string, error)`

- 说明：从指定目录读取Java合约，并编码为交易体。

- 参数：​ 【path】：Java合约根路径。​ 【params】：初始化Java合约的参数。

- 返回​ 【返回值1】：编码后交易体。 【返回值2】：略。

- 实例

```go
payload, err := java.ReadJavaContract("../../conf/contract/contract01")
```

3.11.4 编码Java合约调用方法

`func EncodeJavaFunc(methodName string, params ...string) []byte`

- 说明：调用Java合约时，指定调用的方法名和参数来编码成payload。

- 参数​ 【methodName】：方法名。​ 【params】：调用参数。

- 返回​ 【返回值1】：编码后的字节切片。

- 实例

```go
data := java.EncodeJavaFunc("issue", "1", "2")
```

3.11.5 十六进制字符串转字节切片

`func FromHex(s string) []byte`

- 说明：将16进制的字符串转为[]byte。

- 参数​ 【s】：16进制字符串。

- 返回​ 【返回值1】：字节切片。

- 实例

```go
data := common.FromHex(origin)
```

3.11.6 字节切片转十六进制字符串

`func ToHex(b []byte) string`

- 说明：将[]byte转为16进制的字符串。

- 参数​ 【b】：字节切片。

- 返回​ 【返回值1】：16进制字符串。

- 实例

```go
hexString := common.ToHex(data)
```

ABI结构体主要应用于实例化合约部署、合约调用结构体的构造，还有合约返回值的解码，合约日志（event）的解码。

3.11.7 初始化ABI结构体

`func JSON(reader io.Reader) (ABI, error)`

- 说明：ABI结构体的构造函数。

- 参数【reader】：包含solidity合约abi的reader。

- 返回【返回值1】：ABI结构体。【返回值2】：略。

- 实例

```go
ABI, serr := abi.JSON(strings.NewReader(cr.Abi[0]))
```

3.11.8 编译Solidity合约

`func CompileSourcefile(source string) ([]string, []string, []string, error)`

- 说明：Go调用Solidity编译器编译合约。

- 参数：【source】：合约源码。

- 返回：【返回值1】：abis【返回值2】：bins【返回值3】：types【返回值4】：略。

- 实例

```go
abi, bin, type, _ := compile.CompileSourcefile(code)
```

3.11.9 通过方法名解码input

`func (abi *ABI) UnpackInput(v interface{}, methodName string, data []byte) (err error)`

- 说明：通过abi方法名解析input。

- 参数：【v】：解码结构体【methodName】：解码方法【data】：源数据

- 返回：【返回值1】：略。

- 实例

```go
var r1 uint32
var r2 uint32
testV := []interface{}{&r1, &r2}
err := ABI.UnpackInput(&testV, "add", data)
```

3.11.10 通过源数据解码input

`func (abi *ABI) UnpackInputWithoutMethod(v interface{}, data []byte) (err error)`

- 说明：通过abi方法名解析input，不需要方法名。

- 参数：【v】：解码结构体【data】：源数据

- 返回：【返回值1】：略。

- 实例

```go
var r1 uint32
var r2 uint32
testV := []interface{}{&r1, &r2}
err := ABI.UnpackInputWithoutMethod(&testV, data)
```

3.11.11

`func GenPayload(beanAbi *BeanAbi, params ...interface{}) ([]byte, error)`

- 说明：调用HVM合约时，指定调用的方法名和参数来编码成payload。

- 参数​ 【beanAbi】：beanAbi名称。​ 【params】：调用参数。

- 返回​ ：编码后的字节切片。

- 实例

```go
    invokePayload, sysErr := hvm.GenPayload(beanAbi, "true", "c", "20", "100", "1000", "10000", "1.1", "1.11", "string", person1, bean1,
        []interface{}{"strList1", "strList2"},
        []interface{}{person1, person2},
        []interface{}{[]interface{}{"person1", person1}, []interface{}{"person2", person2}},
        []interface{}{[]interface{}{"bean1", bean1}, []interface{}{"bean2", bean2}})
```

3.11.12

`func Validate(account string, proofPath *AccountProofPath) bool`

- 说明：得到账户的证明路径后，调用验证方法验证证明路径的正确性。

- 参数​ 【account】：账户地址。​ 【path】：证明路径。

- 返回​ ：验证结果。

- 实例

```go
   proofPath, err := rpcAPI.GetAccountProof(key.GetAddress())
   res := rpcAPI.Validate(key.GetAddress(),proofPath)
```

### 3.12 MQ相关接口

3.12.1 获取MQ请求客户端

`func (r *RPC) GetMqClient() *MqClient`

- 说明：管理发送MQ请求。

- 返回：【返回值1】：发送MQ请求客户端。

- 实例

```go
client := rpc.GetMqClient()
```

3.12.2 与broker建立连接

`func (mc *MqClient) InformNormal(id uint, brokerURL string) (bool, StdError)`

- 说明：申请与broker建立连接

- 参数：【id】：节点ID【brokerURL】：指定broker的URL，可为空，为空表示使用平台配置的默认URL

- 返回：【返回值1】：是否正常建立连接【返回值2】：略

- 实例

```go
success, err := client.InformNormal(1, "")
```

3.12.3 注册MQ channel

`func (mc *MqClient) Register(id uint, meta *RegisterMeta) (*QueueRegister, StdError)`

- 说明：向broker注册queue，并提供订阅的事件的参数

- 参数：【id】：节点ID【meta】：事件相关参数

- 返回：【返回值1】：成功订阅的queue名称和exchanger名称【返回值2】：略

- 实例

```go
var hash common.Hash
hash.SetString("123")
rm := NewRegisterMeta(guomiKey.GetAddress(), "node1queue1", MQBlock).SetTopics(1, hash)
rm.Sign(guomiKey)
regist, err := client.Register(1, rm)
```

3.12.4 注销MQ channel

`func (mc *MqClient) UnRegister(id uint, meta *UnRegisterMeta) (*QueueUnRegister, StdError)`

- 说明：通知broker删除指定的queue

- 参数：【id】：节点ID【meta】：事件相关参数

- 返回：【返回值1】：取消订阅是否成功，队列名称，队列中剩余的未消费消息个数【返回值2】：略

- 实例

```go
meta := NewUnRegisterMeta(guomiKey.GetAddress(), "node1queue1", "global_fa34664e_1541655230749576905")
meta.Sign(guomiKey)
unRegist, err := client.UnRegister(1, meta)
```

3.12.5 获取所有队列名称

`(mc *MqClient) GetAllQueueNames(id uint) ([]string, StdError)`

- 说明：获取所有队列名称

- 参数：【id】：节点ID

- 返回：【返回值1】：所有队列的名称数组【返回值2】：略

- 实例

```go
queues, err := client.GetAllQueueNames(1)
```

### 3.13 PROOF相关接口

3.13.1 获取指定账户的证明路径

`func (rpc *RPC) GetAccountProof(account string) (*AccountProofPath, StdError)`

- 说明：获取指定账户的证明路径。

- 返回：【返回值1】：指定账户的证明路径。

- 实例

```go
proofPath, err := rpcAPI.GetAccountProof(key.GetAddress())
```

## 4.1 GoSDK涉及类型说明

该小结用来介绍GoSDK接口涉及的一些结构体的介绍。

### 4.1.1 Block相关

```go
type Block struct {
   Version      string  // hyperchain版本
   Number       uint64  // 区块号
   Hash         string  // 区块哈希
   ParentHash   string  // 父区块哈希
   WriteTime    uint64    // 区块写入时间，unix时间戳
    AvgTime      int64    // 该区块内交易执行的平均时间，单位(ms)
   TxCounts     uint64    // 该区块的交易数量   
   MerkleRoot   string  // 默克尔树根哈希
   Transactions []TransactionInfo // 该区块内的交易列表
}
```

```go
type BlockInterval struct {
   SumOfBlocks uint64  // 区块总数
   StartBlock  uint64  // 起始区块号
   EndBlock    uint64  // 结束区块号
}
```

```go
type TPSInfo struct {
   StartTime     string 
   EndTime       string
   TotalBlockNum uint64  // 区块总数
   BlocksPerSec  float64 // 每秒区块生成数
   Tps           float64 // 每秒交易生成数
}
```

### 4.1.2 Transaction相关

```go
type TransactionInfo struct {
   Version     string  // hyperchain版本
   Hash        string  // 交易哈希
   BlockNumber uint64  // 该交易所在区块号
   BlockHash   string  // 该交易所在区块哈希
   TxIndex     uint64  // 该交易在区块中的偏移量(下标)
   From        string  // 交易发起方地址
   To          string  // 交易接收方地址
   Amount      uint64  // 转账金额
   Timestamp   uint64  // 交易发生时间戳
   Nonce       uint64  // 随机数
   ExecuteTime int64   // 该交易执行时间，单位(ms)
   Payload     string  // 交易体荷载
   Extra       string  // 存证信息
   Invalid     bool    // 是否为非法交易
   InvalidMsg  string  // 非法交易的描述信息
}
```

```go
type TxReceipt struct {
   TxHash          string  // 交易哈希
   ContractAddress string  // 合约地址
   Ret             string  // 合约调用返回值
   Log             []TxLog // 合约event
   VMType          string  // 交易对应的虚拟机类型
   Version         string  // hyperchain版本
}
```

```go
type FailedResult struct {
    Hash  string //未查询到相应交易或回执的hash
    Error string //错误原因
}
```

```go
type TxsInfoAndFailedResults struct {
    TxsInfo   []TransactionInfo //可查询到的交易
    FailedTxs []FailedResult    //失败列表
}
```

```go
type TxReceiptsAndFailedResults struct {
    TxReceipts     []TxReceipt    //可查询到的回执
    FailedReceipts []FailedResult //失败列表
}
```

```go
type TxLog struct {
   Address     string    // 合约地址
   Topics      []string  // 包含事件签名和indexed参数值
   Data        string    // 非indexed事件参数值
   BlockNumber uint64    // 区块号
   TxHash      string    // 交易哈希
   TxIndex     uint64     // 交易所在区块中的索引
   Index       uint64    // 该日志在本条交易产生的所有日志中的偏移量
}
```

```go
type TransactionsCount struct {
   Count     uint64  // 交易数目
   Timestamp uint64  // 查询触发时的时间戳
}
```

```go
type TransactionsCountByContract struct {
   Count        uint64  // 交易数目
   LastIndex    uint64  // 最后一条交易在区块中的索引号
   LastBlockNum uint64  // 最后一条交易所在的区块号
}
```

### 4.1.3 Node相关

```go
type NodeInfo struct {
    Status    uint       // 节点状态
    IP        string     // 节点IP
    Port      string     // 节点端口
    ID        uint       // 节点index
    Isprimary bool       // 是否主节点
    Delay     uint       // 表示该节点与本节点的延迟时间（单位ns），若为0，则为本节点
    IsVp      bool       // 是否VP节点
    Namespace string     // 节点namespace
    Hash      string     // 节点hash
    HostName  string     // 节点主机名
}
```

### 4.1.4 Web socket 相关

```go
type Subscription struct {
   Event          EventType      `json:"event"`  // 事件类型
   SubscriptionID SubscriptionID `json:"subId"`  // 订阅ID
}
```

### 4.1.5 error相关

GoSDK中和网络相关的RPC调用均返回StdError类型的错误。StdError为一个接口。

```go
type StdError interface {
   fmt.Stringer
   error
   Code() int
}
```

所以，该结构提供三种方法：

`func (re *RetError) Error() string`

- 说明：实现了error接口。

- 返回【返回值1】：返回错误描述。

`func (re *RetError) String() string`

- 说明：实现了fmt.Stringer接口。

- 返回【返回值1】：将会返回以下格式字符串"error code: %d, error reason: %s", re.Code(), re.Error()，即错误码加错误描述。

`func (re *RetError) Code() int {`

- 说明：返回错误码。

- 返回【返回值1】：错误码。错误码表如下。

| code   | 含义                                                       |
| ------ | ---------------------------------------------------------- |
| 0      | 请求成功                                                   |
| -32700 | 服务端接收到无效的json。该错误发送于服务器尝试解析json文本 |
| -32600 | 无效的请求（比如非法的JSON格式）                           |
| -32601 | 方法不存在或者无效                                         |
| -32602 | 无效的方法参数                                             |
| -32603 | JSON-RPC内部错误                                           |
| -32000 | Hyperchain内部错误或者空指针或者节点未安装solidity环境     |
| -32001 | 请求的数据不存在                                           |
| -32002 | 余额不足                                                   |
| -32003 | 签名非法                                                   |
| -32004 | 合约部署出错                                               |
| -32005 | 合约调用出错                                               |
| -32006 | 系统繁忙                                                   |
| -32007 | 交易重复                                                   |
| -32008 | 合约操作权限不够                                           |
| -32009 | (合约)账户不存在                                           |
| -32010 | namespace不存在                                            |
| -32011 | 账本上无区块产生，查询最新区块的时候可能抛出该错误         |
| -32012 | 订阅不存在                                                 |
| -32013 | 数据归档、快照相关错误                                     |
| -32096 | http请求处理超时                                           |
| -32097 | Hypercli用户令牌无效                                       |
| -32098 | 请求未带cert或者错误cert导致认证失败                       |
| -32099 | 请求tcert失败                                              |
| -9995  | 请求失败(通常是请求体过长)                                 |
| -9996  | 请求失败(通常是请求消息错误)                               |
| -9997  | 异步请求失败                                               |
| -9998  | 请求超时                                                   |
| -9999  | 获取平台响应失败                                           |

### 4.1.6 MQ相关

```go
type RegisterMeta struct {
    //queue related
    RoutingKeys []routingKey `json:"routingKeys,omitempty"`
    QueueName   string       `json:"queueName,omitempty"`
    //self info
    From      string `json:"from,omitempty"`
    Signature string `json:"signature,omitempty"`
    // block accounts
    IsVerbose bool `json:"isVerbose"`
    // vm log criteria
    FromBlock string           `json:"fromBlock,omitempty"`
    ToBlock   string           `json:"toBlock,omitempty"`
    Addresses []common.Address `json:"addresses,omitempty"`
    Topics    [4][]common.Hash `json:"topics,omitempty"`
}
```

```go
// UnRegisterMeta UnRegisterMeta
type UnRegisterMeta struct {
    From         string
    QueueName    string
    ExchangeName string
    Signature    string
}
```

```go
// QueueRegister MQ register result
type QueueRegister struct {
    QueueName     string
    ExchangerName string
}
```

```go
// QueueUnRegister MQ unRegister result
type QueueUnRegister struct {
    Count   uint
    Success bool
    Error   error
}
```

### 4.1.7 归档相关

```go
// Manifest represents all basic information of a snapshot.
type Manifest struct {
	Height         uint64 `json:"height"`
	Genesis        uint64 `json:"genesis"`
	BlockHash      string `json:"hash"`
	FilterID       string `json:"filterId"`
	MerkleRoot     string `json:"merkleRoot"`
	Namespace      string `json:"Namespace"`
	TxCount        uint64 `json:"txCount"`
	InvalidTxCount uint64 `json:"invalidTxCount,omitEmpty"`
	Status         uint   `json:"status"`
	DBVersion      string `json:"dbVersion"`
	// use for hyperchain
	Date string `json:"date"`
}
```

```go
// Manifests
type Manifests []Manifest
```

```go
// SnapshotEvent
type SnapshotEvent struct {
    FilterId    string `json:"filterId"`
    BlockNumber uint64 `json:"blockNumber"`
}
```

```go
// ArchiveResult used for return archive result, tell caller which step is processing
type ArchiveResult struct {
	FilterID string `json:"filterId"`
	Status   string `json:"status"`
	Reason   string `json:"reason"`
}
```

## 第五章 样例数据

**样例账户（DES）**

**账户地址** from :

`0x0b110ed15f21a3ec73b051b59864ed6dec687ad9`

from 加密存储的私钥字符串:

```json
{"address":"0b110ed15f21a3ec73b051b59864ed6dec687ad9","algo":"0x02","encrypted":"787cf33e169a914f6d6bd2f75f0a9c2fc6e746f2b9e81dfa92b29bb767759bb8c6f431ea0e0e22ac","version":"1.0"}
```

**账户地址** to :

`0xf4d69ac5dc63869de4dc6add25690ad404641bf5`

to 加密存储的私钥字符串:

```json
{"address":"f4d69ac5dc63869de4dc6add25690ad404641bf5","algo":"0x02","encrypted":"e86ba44da8b325e4c177c3ef4d2b7e68a23bc8e1249d7efcc03875c3087935f4c6f431ea0e0e22ac","version":"1.0"}
```

**Solidity 源代码**

```solidity
contract Accumulator{ uint32 sum = 0; function increment(){ sum = sum + 1;
} function getSum() returns(uint32){ return sum; } function add(uint32 num1,uint32 num2) { sum = sum+num1+num2; } }
```

**Solidity bin**

```javascript
0x6060604052600080546 3ffffffff19168155609 e908190601e90 396000f3606060405260 e060020a60003504633ad14af381146030578 063569c5f6d146056578 063d09de08a14607c57 5b6002565b34600257600 0805463ffffffff81166004350 16024350163 ffffffff199091161790 555b005b34600257600 05463ffffffff 166040805163 ffffffff90921682525190 81900360200 190f35b34600 25760546000 805463ffffffff 19811663ffffffff909116600101 17905556
```

**Solidity ABI**

```json
[{"constant":false,"inputs":[{"name":"num1","type":"uint32"},{"name":"num2","type":"uint32"}],"name":"add","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"getSum","outputs":[{"name":"","type":"uint32"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"increment" ,"outputs":[],"payable":false,"type":"function"}]
```

