// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package subscriptions

// [START pubsub_subscriber_sync_pull]
import (
	"context"
	"fmt"
	"io"

	pubsubV1 "cloud.google.com/go/pubsub/apiv1"
	pb "google.golang.org/genproto/googleapis/pubsub/v1"
)

func pullMsgsSync(w io.Writer, projectID, subName string, maxMessages int32) error {
	// projectID := "my-project-id"
	// subName := projectID + "-example-sub"
	// maxMessages := 5
	ctx := context.Background()

	// Instantiate a client. Note, this uses an autogenerated versioned client,
	// cloud.google.com/go/pubsub/apiv1, which differs than the usual
	// cloud.google.com/go/pubsub client that should be used for most other cases.
	// In this case, calling the synchronous pull method requires the underlying
	// versioned client.
	subClient, err := pubsubV1.NewSubscriberClient(ctx)
	if err != nil {
		return fmt.Errorf("Client instantiation error: %v", err)
	}
	sub := fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subName)

	req := &pb.PullRequest{
		Subscription:      sub,
		ReturnImmediately: false,
		MaxMessages:       maxMessages,
	}

	resp, err := subClient.Pull(ctx, req)
	if err != nil {
		return fmt.Errorf("Pull error: %v", err)
	}
	var ackIDs []string
	for _, msg := range resp.GetReceivedMessages() {
		fmt.Fprintf(w, "Got message %q\n", string(msg.GetMessage().Data))
		_ = msg // TODO: handle message.
		ackIDs = append(ackIDs, msg.GetAckId())
	}

	subClient.Acknowledge(ctx, &pb.AcknowledgeRequest{
		Subscription: sub,
		AckIds:       ackIDs,
	})

	return nil
}

// [END pubsub_subscriber_sync_pull]