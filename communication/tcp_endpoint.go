package communication

import (
	"fmt"
	"net"
)
type TCPServer struct {
	Address []byte
	Port int
}
func (s TCPServer) StartTCP() {
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
		go handleConnection(conn)
	}

}
func handleConnection(conn net.Conn){
	defer conn.Close()
	/*buf := make([]byte, 512)
	for {
		nRead, err := conn.Read(buf[:])
		// Do stuff with the read bytes
		if err != nil {
			fmt.Print("Error: ", err)
		}
		fmt.Printf("packet-received: bytes=%d",
			nRead)

		buffer := string(buf[:])
		buffer = strings.TrimSpace(buffer) //Returnt ein slice vom string, ohne whitespaces
		splitBuffer := strings.Split(buffer, "\n")
		for _, data := range splitBuffer {
			log.Println(data)
			conn.WriteTo(buf[:], addrFrom)
		}

	}*/
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		addr := conn.RemoteAddr()
		handleTCPData(n,buf,addr)
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}

/**
Handles the data from the udp connection
*/
func handleTCPData(n int, buffer []byte, addr net.Addr) string{
	writeData := ""
	fmt.Printf("\n--------------\n")
	fmt.Printf("packet-received: bytes=%d from=%s over tcp\n",
		n, addr.String())
	fmt.Println("from", addr, "-", buffer[:n])
	fmt.Printf("\n--------------\n")
	return writeData
}
/**
Checks if errors are thrown. If yes, it prints the error and exits the program
*/
/*func checkError(err error) {

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		//Exits the program
		os.Exit(1)
	}
}*/