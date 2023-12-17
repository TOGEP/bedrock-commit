package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"

	"github.com/joho/godotenv"
)

const ModelID = "amazon.titan-text-express-v1"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		// NOTE: 現在はクレデンシャルを環境変数で設定しているため、.envがなくても動く
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")), config.WithSharedConfigProfile(os.Getenv("AWS_PROFILE")))
	if err != nil {
		log.Fatal(err)
	}

	brc := bedrockruntime.NewFromConfig(cfg)

	cmd := exec.Command("git", "diff")
	var diffBuffer bytes.Buffer
	cmd.Stdout = &diffBuffer

	if err := cmd.Run(); err != nil {
		log.Fatal("Error running git diff:", err)
	}

	// FIXME: プロンプト調整
	prompt := fmt.Sprintf(
		"Analyze the following git diff and generate concise, present-tense commit messages based on the changes.\n"+
			"- Generate up to 5 commit messages, each as a single line.\n"+
			"- Write your response using the imperative tense following the kernel git commit style guide.\n"+
			"\n%s",
		diffBuffer.String(),
	)

	payload := Request{
		InputText: prompt,
		TextGenerationConfig: TextGenerationConfig{
			MaxTokenCount: 4096,
			StopSequences: []string{},
			Temperature:   0.9,
			TopP:          1.0,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	output, err := brc.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(ModelID),
		ContentType: aws.String("application/json"),
		Body:        payloadBytes,
	})
	if err != nil {
		log.Fatal("failed to invoke model: ", err)
	}

	var resp Response
	err = json.Unmarshal(output.Body, &resp)
	if err != nil {
		log.Fatal("failed to unmarshal", err)
	}

	if len(resp.Results) > 0 {
		fmt.Println(resp.Results[0].OutputText)
	} else {
		fmt.Println("No results found")
	}
}

type TextGenerationConfig struct {
	MaxTokenCount int      `json:"maxTokenCount"`
	StopSequences []string `json:"stopSequences"`
	Temperature   float64  `json:"temperature"`
	TopP          float64  `json:"topP"`
}

type Request struct {
	InputText            string               `json:"inputText"`
	TextGenerationConfig TextGenerationConfig `json:"textGenerationConfig"`
}

type Response struct {
	InputTextTokenCount int `json:"inputTextTokenCount"`
	Results             []struct {
		TokenCount       int    `json:"tokenCount"`
		OutputText       string `json:"outputText"`
		CompletionReason string `json:"completionReason"`
	} `json:"results"`
}
