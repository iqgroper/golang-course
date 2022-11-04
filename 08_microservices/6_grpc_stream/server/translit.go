package main

import (
	"fmt"
	"io"

	"gitlab.com/vk-go/lectures-2022-2/08_microservices/6_grpc_stream/translit"

	tr "github.com/gen1us2k/go-translit"
)

type TrServer struct {
	translit.UnimplementedTransliterationServer
}

func (srv *TrServer) EnRu(inStream translit.Transliteration_EnRuServer) error {
	// go func() {

	// 	for {
	// 		inStream.Send(&translit.Word{
	// 			Word: "stat",
	// 		})
	// 		time.Sleep(time.Second)
	// 	}
	// }()
	for {
		// time.Sleep(1 * time.Second)
		inWord, err := inStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		out := &translit.Word{
			Word: tr.Translit(inWord.Word),
		}
		fmt.Println(inWord.Word, "->", out.Word)
		inStream.Send(out)
	}
	return nil
}

func NewTr() *TrServer {
	return &TrServer{}
}
