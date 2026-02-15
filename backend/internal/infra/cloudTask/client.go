package cloudtask

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/o-ga09/zenn-hackthon-2026/internal/queue"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
)

type CloudTaskClient struct {
	client              *cloudtasks.Client
	projectID           string
	location            string
	queue               string
	serviceAccountEmail string
	baseURL             string
}

func NewClient(ctx context.Context) (*CloudTaskClient, error) {
	env := config.GetCtxEnv(ctx)
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud tasks client: %w", err)
	}

	return &CloudTaskClient{
		client:              client,
		projectID:           env.ProjectID,
		location:            env.CLOUD_TASKS_LOCATION,
		queue:               env.CLOUD_TASKS_QUEUE_NAME,
		serviceAccountEmail: env.SERVICE_ACCOUNT_EMAIL,
		baseURL:             env.BASE_URL,
	}, nil
}

func (c *CloudTaskClient) Enqueue(ctx context.Context, task *queue.Task) error {
	payloadBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// タスクタイプに応じてエンドポイントを決定
	var endpoint string
	switch task.Type {
	case "ProcessVLogTask":
		endpoint = "/internal/tasks/create-vlog"
	case "ProcessMediaAnalysisTask":
		endpoint = "/internal/tasks/analyze-media"
	default:
		return fmt.Errorf("unknown task type: %s", task.Type)
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", c.projectID, c.location, c.queue)

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        fmt.Sprintf("%s%s", c.baseURL, endpoint),
					Body:       payloadBytes,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					AuthorizationHeader: &cloudtaskspb.HttpRequest_OidcToken{
						OidcToken: &cloudtaskspb.OidcToken{
							ServiceAccountEmail: c.serviceAccountEmail,
						},
					},
				},
			},
		},
	}

	// NOTE: 親ctxのデッドラインをCloud Tasks APIリクエストに引き継がないよう、
	// 値のみ伝播する新しいコンテキストを使用する
	apiCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()

	_, err = c.client.CreateTask(apiCtx, req)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}
