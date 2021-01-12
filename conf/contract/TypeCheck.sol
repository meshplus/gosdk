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