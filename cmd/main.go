package main

import "github.com/Open-Twin/alexandria/raft"

func main(){
    /*cfg := raft.RawConfig{
        BindAddress: "192.168.0.1",
        JoinAddress: "192.168.0.1",
        RaftPort: 1000,
		HTTPPort: 8080,
		DataDir: "./raft",
		Bootstrap: false,
    }
	//cfg := raft.ReadRawConfig()
    result, err := cfg.ValidateConfig()

    fmt.Println(result)
    fmt.Println(err)*/
	raft.Main()

}
