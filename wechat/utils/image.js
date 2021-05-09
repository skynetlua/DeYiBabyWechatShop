function imageUtil(originalWidth, originalHeight) {
    let imageSize = {}
    wx.getSystemInfo({
        success: function(res) {
            const windowWidth = res.windowWidth
            imageSize.x = 0
            imageSize.y = 0
            imageSize.windowWidth = windowWidth
            imageSize.imageWidth = originalWidth
            imageSize.imageHeight = originalHeight
            if (originalWidth > windowWidth) {
                imageSize.imageWidth = windowWidth
                imageSize.imageHeight = windowWidth * originalHeight / originalWidth
            } else {
                imageSize.x = (windowWidth - originalWidth) / 2
            }
        }
    })
    return imageSize
}

function imageSizeIsLessLimitSize(imagePath, limitSize, lessCallBack, moreCallBack) {
    wx.getFileInfo({
        filePath: imagePath,
        success(res) {
            console.log("压缩前图片大小:", res.size / 1024, 'kb');
            if (res.size > 1024 * limitSize) {
                moreCallBack();
            } else {
                lessCallBack();
            }
        }
    })
}

function getLessLimitSizeImage(canvasId, imagePath, limitSize = 100, drawWidth, callBack) {
    imageSizeIsLessLimitSize(imagePath, limitSize,
        (lessRes) => {
            callBack(imagePath);
        },
        (moreRes) => {
            wx.getImageInfo({
                src: imagePath,
                success: function(imageInfo) {
                    var maxSide = Math.max(imageInfo.width, imageInfo.height);
                    //画板的宽高默认是windowWidth
                    var windowW = drawWidth;
                    var scale = 1;
                    if (maxSide > windowW) {
                        scale = windowW / maxSide;
                    }
                    var imageW = Math.floor(imageInfo.width * scale);
                    var imageH = Math.floor(imageInfo.height * scale);
                    console.log('调用压缩', imageW, imageH);
                    getCanvasImage(canvasId, imagePath, imageW, imageH,
                        (pressImgPath) => {
                            getLessLimitSizeImage(canvasId, pressImgPath, limitSize, drawWidth * 0.7, callBack);
                        }
                    );
                }
            })
        }
    )
}

function getCanvasImage(canvasId, imagePath, imageW, imageH, getImgsuccess) {
    const ctx = wx.createCanvasContext(canvasId);
    ctx.drawImage(imagePath, 0, 0, imageW, imageH);
    ctx.draw(false, () => {
        wx.canvasToTempFilePath({
            canvasId: canvasId,
            x: 0,
            y: 0,
            width: imageW,
            height: imageH,
            quality: 1,
            success(res) {
                getImgsuccess(res.tempFilePath);
            }
        });
    });
}


module.exports = {
    imageUtil: imageUtil,
    getLessLimitSizeImage: getLessLimitSizeImage,
}