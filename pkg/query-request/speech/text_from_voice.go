package speech

import (
	"context"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/api/option"
)

var KeyJSON = ""



func TranscribeAudio(audioBytes []byte) string {
	ctx := context.Background()

	client, err := speech.NewClient(ctx, option.WithCredentialsFile(KeyJSON))
	if err != nil {
		return ""
	}
	defer client.Close()

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 44100,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audioBytes},
		},
	}

	resp, _ := client.Recognize(ctx, req)

	var sb strings.Builder
	for _, result := range resp.Results {
		for _, alternative := range result.Alternatives {
			sb.WriteString(alternative.Transcript + " ")
		}
	}

	return sb.String()

}
