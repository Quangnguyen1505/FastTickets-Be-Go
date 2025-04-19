package context

import (
	"context"
	"errors"
	"log"
)

type InfoUserUUID struct {
	UserId      string
	UserAccount string
}

func GetSubjectUUID(ctx context.Context) (string, error) {
	sUUID, ok := ctx.Value("subjectUUID").(string)
	if !ok {
		return "", errors.New("failed to get subject UUID")
	}

	return sUUID, nil
}

func GetUserIdFromUUID(ctx context.Context) (string, error) {
	sUUID, err := GetSubjectUUID(ctx)
	log.Println("sUUID::", sUUID)
	if err != nil {
		return "", err
	}

	//get infoUser from uuid
	// var infoUser InfoUserUUID
	// if err := cache.GetCache(ctx, sUUID, &infoUser); err != nil {
	// 	return "", err
	// }

	// return infoUser.UserId, nil
	return sUUID, nil
}
