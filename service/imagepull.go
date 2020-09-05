package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/rahultripathidev/docker-utility/registry"
	"strings"
)

func Pullimage(ctx *context.Context, targetNode *client.Client, imageUri string) error {
	//Using Custom Auth Token from username
	//auth := types.AuthConfig{
	//	Username: "UNAME",
	//	Password: "PASS",
	//}
	//_buf, err := json.Marshal(auth)
	//authToken := base64.URLEncoding.EncodeToString(_buf)
	//using AWS SDK to get auth token
	token ,err := registry.GetRegistryAuthorizationToken()
	if err != nil {
		return err
	}
	imagePullLog, err := targetNode.ImagePull(*(ctx), imageUri, types.ImagePullOptions{
		RegistryAuth: token,
		All:          false,
	})
	defer imagePullLog.Close()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Pulling Image")
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(imagePullLog)
		if err == nil {
			logStr := strings.Split(buf.String(), "\n")
			for _, log := range logStr {
				LogData := struct {
					Status string `json:"status"`
				}{}
				err := json.Unmarshal([]byte(log), &LogData)
				if err == nil {
					fmt.Println(LogData.Status)
				}

			}
		}
	}
	return nil
}

//
func GetImageId(imageName string, images []types.ImageSummary) string {
	imageId := ""
	for _, imageData := range images {
		for _, tags := range imageData.RepoTags {
			if tags == imageName {
				imageId = imageData.ID
				break
			}
			if imageId != "" {
				break
			}
		}
	}
	return imageId
}
