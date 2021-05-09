package handle

import (
	"bestsell/common"
	"bestsell/module"
	"fmt"
	"github.com/kataras/iris/v12"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"
)

//=>/upload/file true post {} 
func Upload_file(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	//maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	maxSize := int64(1 << 20)
	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		//ctx.WriteString(err.Error())
		fmt.Println("Upload_file err =", err.Error())
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	form := ctx.Request().MultipartForm
	if form.File["upfile"] == nil || len(form.File["upfile"]) == 0 {
		fmt.Println("Upload_file  no file")
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	file := form.File["upfile"][0]
	playerId := strconv.Itoa(player.ID)
	destDirectory := common.CreateTokenDir(common.UploadPath, playerId)
	fileName := time.Now().Format("20060102150405")+".png"
	saveFile := path.Join(destDirectory, fileName)
	_, err = saveUploadedFile(file, saveFile)
	if err != nil {
		fmt.Println("failed to upload:", file.Filename)
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	data := map[string]interface{}{
		"url": fileName,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

func saveUploadedFile(fh *multipart.FileHeader, saveFile string) (int64, error) {
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()
	out, err := os.OpenFile(saveFile, os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return 0, err
	}
	defer out.Close()
	return io.Copy(out, src)
}