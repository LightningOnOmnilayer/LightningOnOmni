package main

import (
	"context"
	"encoding/json"
	proxy "github.com/omnilaboratory/obd/proxy/pb"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestAddInvoice(t *testing.T) {

	client, conn := getHtlcClient()
	defer conn.Close()

	invoice, err := client.AddInvoice(context.Background(), &proxy.Invoice{
		CltvExpiry: "2021-08-15",
		Value:      0.001,
		PropertyId: 137,
		Private:    false,
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(invoice.PaymentRequest)
}

func TestParseInvoice(t *testing.T) {
	client, conn := getHtlcClient()
	defer conn.Close()
	invoice, err := client.ParseInvoice(context.Background(), &proxy.ParseInvoiceRequest{
		PaymentRequest: "obtb100000s1pqzyfnpwQmZPzUh7Q6PQg6gXB4XheaoZMMhHA9JNeCrJsp3FWjFrAFuzqa5f24dc5d5414d961bba98c98624b87222da3984b324bcab7cfd7fd63aee33b3hzz03cd311a46e7100f775607e231d2538cbfca01c8746ef8a0900e75a15e33456339xq8ps306yqtqp0dqtdescription3hg",
	})
	if err != nil {
		log.Println(err)
		return
	}
	marshal, _ := json.Marshal(invoice)
	log.Println(string(marshal))
}

func TestSendPayment(t *testing.T) {

	client, conn := getHtlcClient()
	defer conn.Close()

	htlcPayment, err := client.SendPayment(context.Background(), &proxy.SendPaymentRequest{
		PaymentRequest: "obtb100000s1pqzyfnpwQmZPzUh7Q6PQg6gXB4XheaoZMMhHA9JNeCrJsp3FWjFrAFuzqa5f24dc5d5414d961bba98c98624b87222da3984b324bcab7cfd7fd63aee33b3hzz03cd311a46e7100f775607e231d2538cbfca01c8746ef8a0900e75a15e33456339xq8ps306yqtqp0dqtdescription3hg",
		InvoiceDetail: &proxy.ParseInvoiceResponse{
			PropertyId:          137,
			Value:               0.001,
			Memo:                "description",
			CltvExpiry:          "2021-08-15",
			H:                   "03367ac2ec89f567433729a448bd0f115cba8f4f48f3b97c6c0e0a19bb226d2aac",
			RecipientNodePeerId: "QmccE4s2uhEXrJXE778NChn1ed8NyWNyAHH23mP7f9NM3L",
			RecipientUserPeerId: "63167817c979ade9e42f3204404c1513a4b1b4e9eea654c9498ed9cc920dbb36",
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	marshal, _ := json.Marshal(htlcPayment)
	log.Println(string(marshal))
}

func getHtlcClient() (proxy.HtlcClient, *grpc.ClientConn) {

	opts := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	return proxy.NewHtlcClient(conn), conn
}
