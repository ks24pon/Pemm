package handlers

import (
    "net/http"
    "os"
    "path/filepath"
    "github.com/labstack/echo/v4"
    "io"
	"log"
    "pemm/models"
    "pemm/database"
)

// データーベースに保存する処理を行う関数
func CreatePost(c echo.Context) error {
	// デバッグ
	log.Println("フォームデータ取得開始")
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
	// デバッグ
	log.Println("投稿成功。リダイレクトします。")
	// 成功後のレスポンス
	return c.Redirect(http.StatusSeeOther, "/index")
}

// 投稿一覧表示
func ListPosts(c echo.Context) error {
	// スライスは投稿のデータを格納
	var posts []models.Post
	// データベースから投稿データを取得し.Errorを使ってエラーチェック
	if err := database.DB.Find(&posts).Error; err != nil {
		return c.String(http.StatusInternalServerError, "投稿一覧を取得できませんでした")
	}
	//デバッグ
	log.Println("================取得した投稿データ",posts)
	//投稿一覧をHTMLテンプレートに渡す
	return c.Render(http.StatusOK, "dogpost_list.html", posts)
}

// 投稿編集
func EditPost(c echo.Context) error {
	// URLからIDを取得
	id := c.Param("id")
	// 取得したデーターを格納
	var post models.Post
	// IDに該当するデーターを取得
	if err := database.DB.First(&post, id).Error; err!= nil {
		// IDが見つからなかった場合は404 Not Found
		return c.String(http.StatusNotFound, "投��が見つかりませんでした")
	}
	// 編集画面表示
	return c.Render(http.StatusOK, "dogpost_edit.html", post)
}