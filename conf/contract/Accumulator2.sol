contract Accumulator {

    event sayHello(int64 addr1, bytes8 indexed msg);
    event saySum(bytes32 msg, uint32 sum);

	uint32 sum = 0;
	bytes32 hello = "hello world";

    function Accumulator(uint32 sum1, bytes32 hello1) {
        sum = sum1;
        hello = hello1;
    }

    function increment() {sum = sum + 1;}

    function getSum() returns(uint32) {return sum;}

    function getHello() constant returns(bytes32) {
        sayHello(1, "test");
        saySum("sum", sum);
        return hello;
    }

    function getMul() returns(bytes, int64, address) {
        return ("hello", 12, msg.sender);
    }

    function add(uint32 num1,uint32 num2) {sum = sum+num1+num2;}
}