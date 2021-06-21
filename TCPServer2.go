package main

import (
	"fmt"
	"net"

	"encoding/binary"
	"math"
)

var Addr []string = make([]string, 0)
var scalex float32
var scaley float32
var scalez float32

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("listen err=", err)
		return
	}
	defer listen.Close()
	for {
		fmt.Println("等待")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("err=", err)

		} else {
			fmt.Println("con=", conn.RemoteAddr().String())
			Addr = append(Addr, conn.RemoteAddr().String())
		}
		go process(conn)

	}
}
func process(conn net.Conn) {
	defer conn.Close()
	for {
		if conn.RemoteAddr().String() == Addr[0] {
			fmt.Println("host")
			buf := make([]byte, 12)
			_, err := conn.Write(buf)
			if err != nil {
				fmt.Println("User1 logout", err)
				Addr = Addr[1:]
				return
			}
			for {
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println("User1 logout", err)
					Addr = Addr[1:]
					return
				}
				fmt.Print(BytesToFloats32(buf[:n]))
				scalex = BytesToFloats32(buf[:n])[0]
				scaley = BytesToFloats32(buf[:n])[1]
				scalez = BytesToFloats32(buf[:n])[2]
			}
		} else {
			fmt.Println("client")
			buf := make([]byte, 12)
			buf[0] = 1
			_, err := conn.Write(buf)
			if err != nil {
				fmt.Println("User1 logout", err)
				Addr = Addr[1:]
				return
			}
			// SendMsg := make([]byte, 1)
			// var SendMsg float64 = 60.8
			temScalex := scalex
			temScaley := scaley
			temScalez := scalez
			for {
				if temScalex != scalex || temScaley != scaley || temScalez != scalez {

					floats := []float32{scalex, scaley, scalez}
					_, err := conn.Write(Floats32ToBytes(floats))
					if err != nil {
						for i, v := range Addr {
							if v == conn.RemoteAddr().String() {
								Addr = append(Addr[0:i], Addr[i+1:]...)
								break
							}
						}
						break
					}
					temScalex = scalex
					temScaley = scaley
					temScalez = scalez
				}
			}
		}
	}
}

func Floats32ToBytes(floats []float32) []byte {
	bytes := make([]byte, 0)
	byte := make([]byte, 4)
	for i := 0; i < len(floats); i++ {
		bits := math.Float32bits(floats[i])
		binary.LittleEndian.PutUint32(byte, bits)
		bytes = append(bytes, byte...)
	}

	return bytes
}

func BytesToFloats32(bytes []byte) []float32 {
	floats32 := make([]float32, 0)
	var float float32
	for i := 0; i < 3; i++ {

		bits := binary.LittleEndian.Uint32(bytes[i*4 : (i*4 + 4)])
		float = math.Float32frombits(bits)
		floats32 = append(floats32, float)
	}
	return floats32
}
