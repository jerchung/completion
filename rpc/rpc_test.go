package main

import (
	"fmt"
	"github.com/quarnster/completion/content"
	"github.com/quarnster/completion/java"
	"github.com/quarnster/completion/net"
	"net/rpc/jsonrpc"
	"reflect"
	"testing"
)

func TestRpc(t *testing.T) {
	if c, err := jsonrpc.Dial("tcp", fmt.Sprintf("127.0.0.1%s", port)); err != nil {
		t.Error(err)
	} else {
		defer c.Close()
		tests := []struct {
			complete content.Complete
			rpcname  string
			abs      string
		}{
			{
				&java.Java{},
				"Java.Complete",
				"java://type/java/util/jar/JarEntry",
			},
			{
				&net.Net{},
				"Net.Complete",
				"net://type/System.String",
			},
		}
		for _, test := range tests {
			var a content.CompleteArgs
			var cmp1, cmp2 content.CompletionResult
			a.Location.Absolute = test.abs
			if err := test.complete.Complete(&a, &cmp1); err != nil {
				t.Error(err)
			}
			t.Log("calling", test.rpcname)
			if err := c.Call(test.rpcname, &a, &cmp2); err != nil {
				t.Error(err)
			} else if !reflect.DeepEqual(cmp1, cmp2) {
				t.Errorf("Results aren't equal: %v\n==============\n:%v", cmp1, cmp2)
			}
		}
	}
}

func TestRpcInvalid(t *testing.T) {
	if c, err := jsonrpc.Dial("tcp", fmt.Sprintf("127.0.0.1%s", port)); err != nil {
		t.Error(err)
	} else {
		defer c.Close()
		tests := []struct {
			complete content.Complete
			rpcname  string
			abs      string
		}{
			{
				&java.Java{},
				"Java.Complete",
				"", // Intentional bad request
			},
		}
		for _, test := range tests {
			var a content.CompleteArgs
			var cmp2 content.CompletionResult
			a.Location.Absolute = test.abs
			t.Log("calling", test.rpcname)
			if err := c.Call(test.rpcname, &a, &cmp2); err == nil {
				t.Error("Expected an error, but didn't receive one")
			}
		}
	}
}
