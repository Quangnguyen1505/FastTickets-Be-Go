package grpc

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/ntquang/ecommerce/proto"
)

func UploadImageMiddleware(uploadClient pb.UploadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if fileHeader == nil {
			c.Next()
		}
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "cannot get image file"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "cannot open image file"})
			return
		}
		defer file.Close()

		// read file to []byte
		fileBytes, err := multipartFileToBytes(file)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "read file failed"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := uploadClient.UploadImage(ctx, &pb.UploadImageRequest{
			FileName:  fileHeader.Filename,
			ImageData: fileBytes,
		})
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "upload image failed"})
			return
		}

		fmt.Println("image_url ", res.GetImageUrl())
		c.Set("image_url", res.GetImageUrl())
		c.Next()
	}
}

func multipartFileToBytes(file multipart.File) ([]byte, error) {
	buf := make([]byte, 0)
	tmp := make([]byte, 1024)
	for {
		n, err := file.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			break
		}
	}
	return buf, nil
}
