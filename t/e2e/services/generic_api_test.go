package services

import (
	"encoding/json"
	"fmt"

	"github.com/status-im/status-go/geth/params"
	"github.com/status-im/status-go/sign"
	"github.com/status-im/status-go/t/e2e"

	. "github.com/status-im/status-go/t/utils"
)

type GenericRPCApiSuite struct {
	e2e.BackendTestSuite
}

func (s *GenericRPCApiSuite) testAPIExported(method string, expectExported bool) {
	cmd := fmt.Sprintf(`{"jsonrpc":"2.0", "method": "%s", "params": []}`, method)

	result := s.Backend.CallRPC(cmd)

	var response struct {
		Error *rpcError `json:"error"`
	}

	s.NoError(json.Unmarshal([]byte(result), &response))

	hidden := (response.Error != nil && response.Error.Code == methodNotFoundErrorCode)

	s.Equal(expectExported, !hidden,
		"method %s should be %s, but it isn't",
		method, map[bool]string{true: "exported", false: "hidden"}[expectExported])
}

func (s *GenericRPCApiSuite) initTest(upstreamEnabled bool, statusServiceEnabled bool) error {
	nodeConfig, err := MakeTestNodeConfig(GetNetworkID())
	s.NoError(err)

	nodeConfig.IPCEnabled = false
	nodeConfig.StatusServiceEnabled = statusServiceEnabled
	nodeConfig.HTTPHost = "" // to make sure that no HTTP interface is started

	if upstreamEnabled {
		networkURL, err := GetRemoteURL()
		s.NoError(err)

		nodeConfig.UpstreamConfig.Enabled = true
		nodeConfig.UpstreamConfig.URL = networkURL
	}

	return s.Backend.StartNode(nodeConfig)
}

func (s *GenericRPCApiSuite) notificationHandlerSuccess(account string, pass string) func(string) {
	return func(jsonEvent string) {
		s.notificationHandler(account, pass, nil)(jsonEvent)
	}
}

func (s *GenericRPCApiSuite) notificationHandler(account string, pass string, expectedError error) func(string) {
	return func(jsonEvent string) {
		envelope := unmarshalEnvelope(jsonEvent)
		if envelope.Type == sign.EventSignRequestAdded {
			event := envelope.Event.(map[string]interface{})
			id := event["id"].(string)
			s.T().Logf("Sign request added (will be completed shortly): {id: %s}\n", id)

			//check for the correct method name
			method := event["method"].(string)
			s.Equal(params.PersonalSignMethodName, method)
			//check the event data
			args := event["args"].(map[string]interface{})
			s.Equal(signDataString, args["data"].(string))
			s.Equal(account, args["account"].(string))

			e := s.Backend.ApproveSignRequest(id, pass).Error
			s.T().Logf("Sign request approved. {id: %s, acc: %s, err: %v}", id, account, e)
			if expectedError == nil {
				s.NoError(e, "cannot complete sign reauest[%v]: %v", id, e)
			} else {
				s.EqualError(e, expectedError.Error())
			}
		}
	}
}
