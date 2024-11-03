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
	//投稿一覧をHTMLテンプレートに渡す
	return c.Render(http.StatusOK, "dogpost_list.html", posts)
}

// 投稿編集画面
func EditPost(c echo.Context) error {
	// URLからIDを取得
	id := c.Param("id")
	// 取得したデーターを格納
	var post models.Post
	// IDに該当するデーターを取得
	if err := database.DB.First(&post, id).Error; err != nil {
		// IDが見つからなかった場合は404 Not Found
		return c.String(http.StatusNotFound, "投稿が見つかりませんでした")
	}
	// 編集画面表示
	return c.Render(http.StatusOK, "dogpost_edit.html", post)
}

// 編集処理
func UpdatePost(c echo.Context) error {
	// URLから投稿IDを取得
	id := c.Param("id")
	// 取得したデータを格納
	var post models.Post
	// ログ
	log.Println("更新処理開始 - 投稿ID", id)
	// 投稿IDに該当するデータを取得
	if err := database.DB.First(&post, id).Error; err != nil {
		// IDが見つからなかった場合は404 Not Found
		return c.String(http.StatusNotFound, "投稿が見つかりませんでした")
	}

	// フォームデータ取得(更新後のデータ)
	name := c.FormValue("name")
	description := c.FormValue("description")
	// 入力値を投稿データに設定
	post.Name = name
	post.Description = description

	// 写真ファイルを取得しアップロードがあれば更新
	file, err := c.FormFile("photo")
	// エラーが発生していない場合(画像の処理)
	if err == nil {
		// 取得したファイルをsrcとしてOpen
		src, err := file.Open()
		if err != nil {
			return err
		}
		// 関数が終了時点でファイルを閉じる
		defer src.Close()

		// 保存先パスを作成
		filePath := filepath.Join("uploads", file.Filename)
		// filePathに新しいファイルを作成してdst変数に格納
		dst, err := os.Create(filePath)
		// ファイル作成失敗後にエラー返す
		if err != nil {
			return err
		}
		defer dst.Close()

		// ファイル内容をコピーして保存
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// 新しい画像のパス指定
		post.PhotoPath = file.Filename
	}

	// データベースの投稿を更新
	if err := database.DB.Save(&post).Error; err != nil {
		return c.String(http.StatusInternalServerError, "投稿の更新に失敗しました")
	}
	// 更新後、投稿の詳細画面へリダイレクト
	return c.Redirect(http.StatusSeeOther, "/posts/"+id)
}

// 詳細画面
func ShowPost(c echo.Context) error {
	// 投稿のIDをパラメータから取得
	id := c.Param("id")
	// 投稿データを格納する変数を制限
	var post models.Post
	// データベースから投稿IDで検索
	if err := database.DB.First(&post, id).Error; err != nil {
		// IDが見つからなかった場合エラー返す
		return c.String(http.StatusNotFound, "投稿が見つかりませんでした")
	}
	// 詳細ページのレンタリング、投稿データを返す
	return c.Render(http.StatusOK, "show_post.html", post)
}

// 削除処理
func DeletePost(c echo.Context) error {
	// 投稿のIDをパラメータから取得
	id := c.Param("id")
	// 削除対象の投稿データを格納する変数を制限
	var post models.Post
	// データベースから投稿IDで検索
	if err := database.DB.First(&post, id).Error; err != nil {
		return c.String(http.StatusNotFound, "投稿が見つかりませんでした")
	}
	// 投稿を削除
	if err := database.DB.Delete(&post).Error; err != nil {
		// 削除にした場合エラーを返す
		return c.String(http.StatusInternalServerError, "投��の削除に失��しました")
	}
	// 削除成功後
	return c.Redirect(http.StatusSeeOther, "/index")
}