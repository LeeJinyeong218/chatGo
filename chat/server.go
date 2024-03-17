package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_assign/dto"
	"go_assign/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func verify_user(id string, db *gorm.DB, c echo.Context) (*model.User, error) {
	user_id, _ := strconv.ParseInt(id, 10, 64) //string -> int
	user := new(model.User)
	db.First(user, "userid = ?", user_id)
	if user.Userid == 0 { // No existing User
		return nil, errors.New("{}")
	}
	return user, nil
}

func verify_chatroom(id string, db *gorm.DB, c echo.Context) (*model.Chatroom, error) {
	chatroom_id, _ := strconv.ParseInt(id, 10, 64)
	chatroom := new(model.Chatroom)
	db.First(chatroom, "chatroom_id = ?", chatroom_id)
	if chatroom.ChatroomId == 0 {
		return nil, errors.New("{}")
	}
	return chatroom, nil
}

func friendlist_string_to_int(friend_list string) ([]int64, error) {
	friendlist_split_array := strings.Split(friend_list, ",")

	friendlist_int := []int64{}

	for _, item := range friendlist_split_array {
		if item != "" {
			item_int, _ := strconv.ParseInt(item, 10, 64)
			fmt.Printf("s > i > %+v %+v\n", item, item_int)

			friendlist_int = append(friendlist_int, item_int)
		}
	}
	return friendlist_int, nil
}

func friendlist_int_to_string(friendlist []int64) ([]string, error) {
	friendlist_string := []string{}
	for _, item := range friendlist {
		item_string := strconv.FormatInt(item, 10)
		friendlist_string = append(friendlist_string, item_string)
	}

	return friendlist_string, nil
}

func add_friend(friend_list string, opponent_id int64) ([]string, error) {
	user_friendlist_items := strings.Split(friend_list, ",")

	new_user_friendlist_items := []int64{}
	fmt.Printf("%+v", new_user_friendlist_items)
	for _, item := range user_friendlist_items {
		if item != "" {
			item_int, _ := strconv.ParseInt(item, 10, 64)
			fmt.Printf("s > i > %+v %+v\n", item, item_int)

			new_user_friendlist_items = append(new_user_friendlist_items, item_int)
		}
	}

	for _, item := range new_user_friendlist_items {
		if item == opponent_id {
			return nil, errors.New("중복")
		}
	}

	new_user_friendlist_items = append(new_user_friendlist_items, opponent_id)
	new_user_friendlist_items_string := []string{}
	for _, item := range new_user_friendlist_items {
		item_string := strconv.FormatInt(item, 10)
		fmt.Printf("i>s %s %+v\n", item_string, new_user_friendlist_items_string)
		new_user_friendlist_items_string = append(new_user_friendlist_items_string, item_string)
	}

	return new_user_friendlist_items_string, nil
}

func delete_friend(friend_list string, opponent_id int64) ([]string, error) {
	user_friendlist_int, _ := friendlist_string_to_int(friend_list)

	existence_delete_user := false

	for _, item := range user_friendlist_int {
		if item == opponent_id {
			existence_delete_user = true
		}
	}

	if existence_delete_user == false {
		return nil, errors.New("user to delete is not on the friendlist")
	}

	for index, item := range user_friendlist_int {
		if item == opponent_id {
			user_friendlist_int = append(user_friendlist_int[:index], user_friendlist_int[index+1:]...)
		}
	}
	user_friendlist_string, _ := friendlist_int_to_string(user_friendlist_int)
	return user_friendlist_string, nil
}

func main() {
	//echo instance
	e := echo.New()

	//middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Chat code

	//db
	dsn := "test:0000@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//1. Make Account
	e.POST("/signup", func(c echo.Context) error {
		newUser := new(model.User)
		if err = c.Bind(newUser); err != nil {
			panic(err)
		}

		searchUser := new(model.User)
		db.First(searchUser, "name = ? and email = ?", newUser.Name, newUser.Email)
		if searchUser.Userid != 0 {
			return c.JSON(http.StatusBadRequest, "{}")
		}

		newUser.Created = time.Now()
		newUser.Updated = time.Now()

		db.Create(newUser)

		return c.JSON(http.StatusOK, newUser)
	})
	// //2. Login
	e.POST("/login", func(c echo.Context) error {
		user := new(model.User)
		if err = c.Bind(user); err != nil {
			panic(err)
		}
		searchUser := new(model.User)
		db.First(searchUser, "email = ?", user.Email)
		if searchUser.Userid == 0 { // No existing User
			return c.JSON(http.StatusBadRequest, "{}")
		} else if searchUser.Password != user.Password { //Wrong Password
			return c.String(http.StatusBadRequest, "wrong password")
		}

		return c.JSON(http.StatusOK, searchUser)
	})

	//3. Add Friend .split->to array /',' 구분 /parsing / []int -> string
	e.POST("/friendlist", func(c echo.Context) error {
		user, err := verify_user(c.FormValue("user_id"), db, c)
		add_user, err := verify_user(c.FormValue("add_id"), db, c)

		user_friendlist, err := add_friend(user.FriendList, add_user.Userid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		add_friendlist, err := add_friend(add_user.FriendList, user.Userid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		//update
		db.Model(&model.User{}).Where("userid = ?", user.Userid).Update("friend_list", strings.Join(user_friendlist, ","))
		db.Model(&model.User{}).Where("userid = ?", add_user.Userid).Update("friend_list", strings.Join(add_friendlist, ","))
		return c.JSON(http.StatusOK, "added successfully")
	})

	//4. Delete Friend in list
	e.POST("/friendlist_delete", func(c echo.Context) error {
		user, err := verify_user(c.FormValue("user_id"), db, c)
		fmt.Printf("%+v\n", err)
		delete_user, err := verify_user(c.FormValue("delete_id"), db, c)
		fmt.Printf("%+v\n", err)

		new_user_friendlist, err := delete_friend(user.FriendList, delete_user.Userid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		new_delete_user_friendlist, err := delete_friend(delete_user.FriendList, user.Userid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		db.Model(&model.User{}).Where("userid = ?", user.Userid).Update("friend_list", strings.Join(new_user_friendlist, ","))
		db.Model(&model.User{}).Where("userid = ?", delete_user.Userid).Update("friend_list", strings.Join(new_delete_user_friendlist, ","))

		return c.JSON(http.StatusOK, "deleted successfully")
	})

	//chatrooms list
	e.GET("/chatting_Room", func(c echo.Context) error {
		user, err := verify_user(c.QueryParam("user_id"), db, c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, user)
		}

		chatrooms := new([]model.Chatroom)
		chatrooms_user := []model.Chatroom{}
		db.Find(chatrooms)
		for index, item := range *chatrooms {
			fmt.Printf("%+v %+v\n", index, item)
			userIds_int, err := friendlist_string_to_int(item.UserIds)
			if err != nil {
				return c.JSON(http.StatusBadRequest, user)
			}
			for _, item2 := range userIds_int {
				if item2 == user.Userid {
					chatrooms_user = append(chatrooms_user, item)
					break
				}
			}

		}

		return c.JSON(http.StatusOK, chatrooms_user)
	})

	//Chatting Room
	e.POST("/chatting_Room", func(c echo.Context) error {
		newChatroom := new(model.Chatroom)
		newChatroom.UserIds = c.FormValue("userIds")

		newChatroom.Created = time.Now()
		newChatroom.Updated = time.Now()

		db.Create(newChatroom)

		return c.JSON(http.StatusOK, newChatroom)
	})

	//text create
	e.POST("/chatting_Room/:chatroomid", func(c echo.Context) error {
		chatroomid, _ := strconv.ParseInt(c.FormValue("chatroomid"), 10, 64)
		chattextDto := new(dto.ChattextDto)
		if err = c.Bind(chattextDto); err != nil {
			panic(err)
		}
		chattext := new(model.Chattext)
		chattext.ChatroomId = int(chatroomid)
		chattext.SenderId = chattextDto.SenderId
		chattext.Text = chattextDto.Text

		chattext.Created = time.Now()
		chattext.Updated = time.Now()

		db.Create(chattext)

		//

		chattextLast := new(model.Chattext)
		db.Last(chattextLast)
		db.Model(&model.Chatroom{}).Where("chatroom_id = ?", chatroomid).Updates(model.Chatroom{LasttextId: chattextLast.TextId, Updated: time.Now()})

		return c.JSON(http.StatusOK, chattext)
	})
	//text load
	e.GET("/chatting_Room/:chatroomid", func(c echo.Context) error {
		chatroom, err := verify_chatroom(c.Param("chatroomid"), db, c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, chatroom)
		}

		chattexts := new([]model.Chattext)
		db.Where("chatroom_id = ?", chatroom.ChatroomId).Order("created ASC").Find(chattexts)

		return c.JSON(http.StatusOK, chattexts)
	})
	//text edit
	e.PUT("/chatting_Room/:chatroomid", func(c echo.Context) error {
		chatroom, err := verify_chatroom(c.Param("chatroomid"), db, c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, chatroom)
		}
		edit_textId, _ := strconv.ParseInt(c.FormValue("edit_text_id"), 10, 64)
		edit_text := c.FormValue("edit_text")

		db.Model(&model.Chattext{}).Where("text_id =?", edit_textId).Updates(model.Chattext{Text: edit_text, Updated: time.Now()})

		return c.JSON(http.StatusOK, "edit complete")
	})

	//text delete
	e.DELETE("/chatting_Room/:chatroomid", func(c echo.Context) error {
		chatroom, err := verify_chatroom(c.Param("chatroomid"), db, c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, chatroom)
		}

		delete_textId, _ := strconv.ParseInt(c.FormValue("delete_text_id"), 10, 64)
		delete_text := new(model.Chattext)
		db.First(delete_text, "text_id = ?", delete_textId)
		if delete_text.TextId == 0 {
			return c.JSON(http.StatusBadRequest, "no existing text")
		}
		//delete
		db.Delete(&model.Chattext{}, delete_textId)
		//chatroom update
		if delete_textId == int64(chatroom.LasttextId) {
			db.Model(&model.Chatroom{}).Where("chatroom_id = ?", delete_text.ChatroomId).Update("lasttext_id", chatroom.LasttextId-1)
		}

		return c.JSON(http.StatusOK, "delete complete")
	})
	//
	// /Chat code
	//start server
	e.Logger.Fatal(e.Start(":1323"))
}
