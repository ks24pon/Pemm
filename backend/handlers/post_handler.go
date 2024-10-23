package handlers

import (
    "net/http"
    "os"
    "path/filepath"
    "github.com/labstack/echo/v4"
    "io"
    "pemm/models"
    "pemm/database"
)

// データーベースに保存する処理を行う関数
func CreatePost(c echo.Context) error {
	// フォームデータ取得
	name := c.FormValue("name")
	description := c.FormValue("description")

	// 写真ファイルを取得
	file, err := c.FormFile("photo")
	if err != nil {
		return c.String(http.StatusBadRequest, "写真のアップロードに失敗しました")
	}

	// ファイルをサーバーに保存
	src, err := file.Open()
	if err != nil {
		return err
	}
	// ファイルを閉じる
	defer src.Close()

	// 保存先パスを作成
	filePath := filepath.Join("uploads", file.Filename)
	dst, err := os.Create(filePath) //保存先ファイルを作成
	// エラーハンドリング
    if err!= nil {
        return err
    }
	// ファイルを閉じる
	defer dst.Close()

	// ファイルの内容をコピーして保存
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// データーベースに保存するための投稿データーを作成
	post := models.Post{
		Name:		name,
		Description: description,
		PhotoPath: file.Filename,
	}

	// データーベースに接続
	if err := database.DB.Create(&post).Error; err != nil {
		return c.String(http.StatusInternalServerError, "データーベースへの保存に失敗しました")
	}

	// 成功後のレスポンス
	return c.Redirect(http.StatusSeeOther, "/album")
}