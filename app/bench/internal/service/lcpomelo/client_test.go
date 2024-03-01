package lcpomelo

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"
)

func BenchmarkClientConnector_AsyncCustomSend(b *testing.B) {
	c := NewClientConnector("ws://localhost:3051", 15735181677, 2, "roomid")

	c.pomeloChatAddress = "ws://localhost:3051"
	err := c.RunChatConnectorAndWaitConnect(context.Background(), time.Second)
	if err != nil {
		b.Fatal(err)
	}

	req := map[string]interface{}{
		"liveId":     "liveId",
		"lecturerId": "lecturerId",
		"tutorId":    "tutorId",
		"stuId":      "stuId",
	}

	data, err := json.Marshal(req)
	if err != nil {
		b.Fatal(err)
	}

	wg := sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		err = c.AsyncCustomSend(context.Background(), "recover.recoverHandler.msgRecover", data, func(data []byte) {
			wg.Done()
		})
		wg.Add(1)
		if err != nil {
			b.Fatal(err)
		}
	}

	wg.Wait()

}
