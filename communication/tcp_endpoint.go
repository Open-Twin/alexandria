package communication

import (
	"fmt"
	"net"
)
type TCPServer struct {
	Address []byte
	Port int
}
type TCPHandler func(addr net.Addr, buf []byte) []byte

func (s TCPServer) Start(handler TCPHandler) {
	addr := net.TCPAddr{
		IP: s.Address,
		Port: s.Port,
		Zone: "",
	}

	// TCP Listener
	listener, err := net.ListenTCP("tcp", &addr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, handler)
	}

}
func handleConnection(conn net.Conn, handler TCPHandler){
	defer conn.Close()

	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		addr := conn.RemoteAddr()
		result := handleTCPData(n,buf,addr, handler)
		_, err2 := conn.Write(result)
		if err2 != nil {
			return
		}
	}
}

/**
Handles the data from the udp connection
*/
func handleTCPData(n int, buffer []byte, addr net.Addr, handler TCPHandler) []byte{
	fmt.Printf("\n--------------\n")
	fmt.Printf("packet-received: bytes=%d from=%s over tcp\n",
		n, addr.String())
	fmt.Println("from", addr, "-", buffer[:n])
	fmt.Printf("\n--------------\n")
	result := handler(addr,buffer)
	return result
}
/**
Checks if errors are thrown. If yes, it prints the error and exits the program
*/
/*func checkError(err error) {

	if err != nil {
		log.Fatal(os.Stderr, "Fatal error: %s", err.Error())
		//Exits the program
		//os.Exit(1)
	}
}*/