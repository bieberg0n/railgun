package railgun

import "fmt"

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Log(args... interface{}) {
	fmt.Println(args...)
}

//func checkSum(msg []byte) int {
//	sum := 0
//	for n := 1; n < len(msg)-1; n += 2 {
//		sum += int(msg[n])*256 + int(msg[n+1])
//	}
//	sum = (sum >> 16) + (sum & 0xffff)
//	sum += sum >> 16
//	return int(^sum)
//}

func checkSum1(msg []byte) uint16 {
	sum := 0
	for n := 1; n < len(msg)-1; n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var ans = uint16(^sum)
	return ans
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	//以每16位为单位进行求和，直到所有的字节全部求完或者只剩下一个8位字节（如果剩余一个8位字节说明字节数为奇数个）
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	//如果字节数为奇数个，要加上最后剩下的那个8位字节
	if length > 0 {
		sum += uint32(data[index])
	}
	//加上高16位进位的部分
	sum += (sum >> 16)
	//别忘了返回的时候先求反
	return uint16(^sum)
}