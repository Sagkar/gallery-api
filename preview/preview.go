package preview

import (
	"github.com/prplecake/go-thumbnail"
	"os"
)

func GeneratePreview(imagePath, uuid string) (string, error) {
	var config = thumbnail.Generator{
		DestinationPath:   imagePath,
		DestinationPrefix: "thumb_",
		Scaler:            "CatmullRom",
	}
	gen := thumbnail.NewGenerator(config)

	thumbName := gen.DestinationPrefix + uuid + ".png"
	dest := "resources/storage/preview/" + thumbName

	i, err := gen.NewImageFromFile(imagePath)
	if err != nil {
		return "", err
	}

	thumbBytes, err := gen.CreateThumbnail(i)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(dest, thumbBytes, 0644)
	if err != nil {
		return "", err
	}

	return thumbName, nil
}
