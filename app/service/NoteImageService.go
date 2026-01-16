package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"regexp"
	"time"
)

type NoteImageService struct {
}

// 通过id, userId得到noteIds
func (this *NoteImageService) GetNoteIds(imageId string) []string {
	query := "SELECT note_id FROM note_images WHERE file_id = $1"
	rows, err := db.DB.Query(query, imageId)
	if err != nil {
		Log(err.Error())
		return nil
	}
	defer rows.Close()

	noteIds := []string{}
	for rows.Next() {
		var noteImage info.NoteImage
		err := rows.Scan(&noteImage.NoteId)
		if err != nil {
			Log(err.Error())
			continue
		}
		noteIds = append(noteIds, noteImage.NoteId)
	}

	if len(noteIds) > 0 {
		return noteIds
	}

	return nil
}

// TODO 这个web可以用, 但api会传来, 不用用了
// 解析内容中的图片, 建立图片与note的关系
// <img src="/file/outputImage?fileId=12323232" />
// 图片必须是我的, 不然不添加
// imgSrc 防止博客修改了, 但内容删除了
func (this *NoteImageService) UpdateNoteImages(userId, noteId, imgSrc, content string) bool {
	// 让主图成为内容的一员
	if imgSrc != "" {
		content = "<img src=\"" + imgSrc + "\" >" + content
	}
	// life 添加getImage
	reg, _ := regexp.Compile("(outputImage|getImage)\\?fileId=([a-z0-9A-Z]{24})")
	find := reg.FindAllStringSubmatch(content, -1) // 查找所有的

	// 删除旧的
	query := "DELETE FROM note_images WHERE note_id = $1"
	_, err := db.DB.Exec(query, noteId)
	if err != nil {
		Log(err.Error())
	}

	// 添加新的
	var fileId string
	noteImage := info.NoteImage{NoteId: noteId, UserId: userId}
	hasAdded := make(map[string]bool)
	if find != nil && len(find) > 0 {
		for _, each := range find {
			if each != nil && len(each) == 3 {
				fileId = each[2] // 现在有两个子表达式了
				// 之前没能添加过的
				if _, ok := hasAdded[fileId]; !ok {
					Log(fileId)
					// 判断是否是我的文件
					if fileService.IsMyFile(userId, fileId) {
						noteImage.FileId = fileId
						noteImage.NoteImageId = db.NewUUID()
						noteImage.CreatedTime = time.Now()
						insertQuery := "INSERT INTO note_images (id, note_id, user_id, file_id, created_time) VALUES ($1, $2, $3, $4, $5)"
						_, err := db.DB.Exec(insertQuery, noteImage.NoteImageId, noteImage.NoteId, noteImage.UserId, noteImage.FileId, noteImage.CreatedTime)
						if err != nil {
							Log(err.Error())
						}
					}
					hasAdded[fileId] = true
				}
			}
		}
	}

	return true
}

// 复制图片, 把note的图片都copy给我, 且修改noteContent图片路径
func (this *NoteImageService) CopyNoteImages(fromNoteId, fromUserId, newNoteId, content, toUserId string) string {
	/* 弃用之
	// 得到fromNoteId的noteImages, 如果为空, 则直接返回content
	noteImages := []info.NoteImage{}
	db.ListByQWithFields(db.NoteImages, bson.M{"NoteId": bson.ObjectIdHex(fromNoteId)}, []string{"ImageId"}, &noteImages)
	if len(noteImages) == 0 {
		return content;
	}
	for _, noteImage := range noteImages {
		imageId := noteImage.ImageId.Hex()
		ok, newImageId := fileService.CopyImage(fromUserId, imageId, toUserId)
		if ok {
			replaceMap[imageId] = newImageId
		}
	}
	*/

	// 因为很多图片上传就会删除, 所以直接从内容中查看图片id进行复制

	// <img src="/file/outputImage?fileId=12323232" />
	// 把fileId=1232替换成新的
	replaceMap := map[string]string{}

	reg, _ := regexp.Compile("(outputImage|getImage)\\?fileId=([a-z0-9A-Z]{24})")
	content = reg.ReplaceAllStringFunc(content, func(each string) string {
		// each = outputImage?fileId=541bd2f599c37b4f3r000003
		// each = getImage?fileId=541bd2f599c37b4f3r000003

		fileId := each[len(each)-24:] // 得到后24位, 也即id

		if _, ok := replaceMap[fileId]; !ok {
			ok2, newImageId := fileService.CopyImage(fromUserId, fileId, toUserId)
			if ok2 {
				replaceMap[fileId] = newImageId
			} else {
				replaceMap[fileId] = ""
			}
		}

		replaceFileId := replaceMap[fileId]
		if replaceFileId != "" {
			if each[0] == 'o' {
				return "outputImage?fileId=" + replaceFileId
			}
			return "getImage?fileId=" + replaceFileId
		}
		return each
	})

	return content
}

func (this *NoteImageService) getImagesByNoteIds(noteIds []string) map[string][]info.File {
	// TODO: Implement PostgreSQL version
	return make(map[string][]info.File)
}
