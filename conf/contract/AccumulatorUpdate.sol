contract AccumulatorUpdate {

    event sayHello(address addr1, bytes32 msg);
    event saySum(bytes32 msg, uint32 sum);

	uint32 sum = 0;
	uint32 sum1 = 0;
	bytes32 hello = "hello world";

    function increment() {sum = sum + 1;}

    function getSum() returns(uint32) {return sum;}

    function getHello() constant returns(bytes32) {
        sayHello(msg.sender, hello);
        saySum("sum", sum);
        return hello;
    }

    function add(uint32 num1,uint32 num2) returns(uint32) {sum = sum+num1+num2; return sum;}

    function addUpdate(uint32 num1, uint32 num2) returns(uint32) {
        sum = sum1 + num1 + num2;
        return sum;
    }
}