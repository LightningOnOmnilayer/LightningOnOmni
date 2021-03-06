package rpc

import (
	"crypto/x509"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/omnilaboratory/obd/bean"
	"github.com/omnilaboratory/obd/bean/enum"
	"github.com/omnilaboratory/obd/lightclient"
	"log"
	"net/http/httptest"
	"strings"
	"time"
)

var (
	loginChan          = make(chan bean.ReplyMessage)
	logoutChan         = make(chan bean.ReplyMessage)
	changePasswordChan = make(chan bean.ReplyMessage)
	openChannelChan    = make(chan bean.ReplyMessage)
	fundChannelChan    = make(chan bean.ReplyMessage)
	rsmcChan           = make(chan bean.ReplyMessage)
	addInvoiceChan     = make(chan bean.ReplyMessage)
	payInvoiceChan     = make(chan bean.ReplyMessage)
	onceRequestChan    = make(chan bean.ReplyMessage)
)

var obcClient *lightclient.Client

func ConnToObd() (err error) {

	obcClient = &lightclient.Client{IsGRpcRequest: true}

	////u := url.URL{Scheme: "wss", Host: "127.0.0.1:60020", Path: "/ws" + config.ChainNodeType}
	//u := url.URL{Scheme: "ws", Host: "127.0.0.1:60020", Path: "/ws" + config.ChainNodeType}
	//log.Printf("grpc begin to connect to obd: %s", u.String())
	//
	//dailer := websocket.DefaultDialer
	////dailer.TLSClientConfig =&tls.Config{InsecureSkipVerify: true}
	//
	//header := http.Header{}
	//header.Add("session", tool.GetGRpcSession())
	//
	//connObd, _, err = dailer.Dial(u.String(), header)
	//if err != nil {
	//	log.Println("fail to dial obd:", err)
	//	return err
	//}
	//connObd.SetReadLimit(1 << 20)
	//
	//go readDataFromObd()

	return nil
}

func rootCAs(s *httptest.Server) *x509.CertPool {
	certs := x509.NewCertPool()
	for _, c := range s.TLS.Certificates {
		roots, err := x509.ParseCertificates(c.Certificate[len(c.Certificate)-1])
		if err != nil {
			log.Println("error parsing server's root cert: %v", err)
		}
		for _, root := range roots {
			certs.AddCert(root)
		}
	}
	return certs
}

func readDataFromObd() {
	defer func() {
		if connObd != nil {
			_ = connObd.Close()
		}
		ConnToObd()
	}()

	go func() {
		for {
			log.Println("waiting message...")
			_, message, err := connObd.ReadMessage()
			log.Println(string(message))
			if err != nil {
				log.Println(err)
				return
			}
			replyMessage := bean.ReplyMessage{}
			_ = json.Unmarshal(message, &replyMessage)

			if currUserInfo != nil {
				if strings.Contains(replyMessage.To, currUserInfo.UserPeerId) == false {
					continue
				}
			}

			switch replyMessage.Type {
			case enum.MsgType_GetMnemonic_2004:
				onceRequestChan <- replyMessage
			case enum.MsgType_UserLogin_2001:
				if strings.Contains(replyMessage.From, replyMessage.To) {
					loginChan <- replyMessage
				}
			case enum.MsgType_User_UpdateAdminToken_2008:
				changePasswordChan <- replyMessage
			case enum.MsgType_UserLogout_2002:
				logoutChan <- replyMessage

			case enum.MsgType_SendChannelOpen_32:
				openChannelChan <- replyMessage
			case enum.MsgType_RecvChannelAccept_33:
				openChannelChan <- replyMessage

			case enum.MsgType_Funding_134:
				fundChannelChan <- replyMessage
			case enum.MsgType_FundingCreate_SendBtcFundingCreated_340:
				fundChannelChan <- replyMessage
			case enum.MsgType_FundingSign_RecvBtcSign_350:
				fundChannelChan <- replyMessage
			case enum.MsgType_ClientSign_AssetFunding_AliceSignRD_1134:
				fundChannelChan <- replyMessage

			case enum.MsgType_CommitmentTx_SendCommitmentTransactionCreated_351:
				if replyMessage.Status == false {
					rsmcChan <- replyMessage
				}
				break
			case enum.MsgType_ClientSign_CommitmentTx_AliceSignC2a_360:
				rsmcChan <- replyMessage

			case enum.MsgType_HTLC_Invoice_402:
				addInvoiceChan <- replyMessage

			case enum.MsgType_HTLC_ParseInvoice_403:
				onceRequestChan <- replyMessage

			case enum.MsgType_HTLC_FindPath_401:
				if replyMessage.Status == false {
					payInvoiceChan <- replyMessage
				}
				break
			case enum.MsgType_HTLC_FinishTransferH_43:
				payInvoiceChan <- replyMessage
			default:
				continue
			}
		}
	}()

	ticker := time.NewTicker(time.Minute * 2)
	defer ticker.Stop()

	defer func(ticker *time.Ticker) {
		if r := recover(); r != nil {
			log.Println("grpc server goroutine recover")
			ticker.Stop()
			connObd = nil
		}
	}(ticker)

	for {
		select {
		case t := <-ticker.C:
			if connObd != nil {
				info := make(map[string]interface{})
				info["type"] = enum.MsgType_HeartBeat_2007
				info["data"] = t.String()
				bytes, err := json.Marshal(info)
				err = connObd.WriteMessage(websocket.TextMessage, bytes)
				if err != nil {
					connObd = nil
					return
				}
			} else {
				return
			}
		}
	}
}

func sendMsgToObd(info interface{}, RecipientNodePeerId, RecipientUserPeerId string, msgType enum.MsgType) {
	if connObd == nil {
		ConnToObd()
	}
	var infoBytes []byte
	if info != nil {
		infoBytes, _ = json.Marshal(info)
	}
	requestMessage := bean.RequestMessage{Data: string(infoBytes), Type: msgType}
	requestMessage.RecipientNodePeerId = RecipientNodePeerId
	requestMessage.RecipientUserPeerId = RecipientUserPeerId
	msg, _ := json.Marshal(requestMessage)

	err := connObd.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		connObd.Close()
		ConnToObd()
		sendMsgToObd(info, RecipientNodePeerId, RecipientUserPeerId, msgType)
		return
	}
}
