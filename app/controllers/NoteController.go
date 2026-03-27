package controllers

import (
	"go-app/db"
	"go-app/models"
	u "go-app/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"fmt"
	"context"
	"time"
)

var NoteCreate = func(w http.ResponseWriter, r *http.Request) {
	Note := &models.Note{}
	err := json.NewDecoder(r.Body).Decode(Note)

	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	database := db.GetDB()
	err = database.Create(Note).Error

	if err != nil {
		u.HandleBadRequest(w, err)
	} else {
		cacheKey:= fmt.Sprintf("note:%d", Note.ID)
		
		noteData, _ := json.Marshal(Note)

		ctx := context.Background()
        	db.GetRedis().Set(ctx, cacheKey, noteData, 10*time.Minute)
		res, _ := json.Marshal(Note)
		u.RespondJSON(w, res)
	}
}

var NoteRetrieve = func(w http.ResponseWriter, r *http.Request) {
	Note := &models.Note{}

	params := mux.Vars(r)
	id := params["id"]

	cacheKey := "note:" + id
	ctx := context.Background()

	val, err := db.GetRedis().Get(ctx, cacheKey).Result()
	if err == nil {
		 json.Unmarshal([]byte(val), Note)
        	 res, _ := json.Marshal(Note)
        	 u.RespondJSON(w, res)
		 return
	}

	database := db.GetDB()
	err = database.First(&Note, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.HandleNotFound(w)
		} else {
			u.HandleBadRequest(w, err)
		}
		return
	}

	noteData, _ := json.Marshal(Note)
    	db.GetRedis().Set(ctx, cacheKey, noteData, 10*time.Minute)

	res, err := json.Marshal(Note)
	if err != nil {
		u.HandleBadRequest(w, err)
	} else if Note.ID == 0 {
		u.HandleNotFound(w)
	} else {
		u.RespondJSON(w, res)
	}
}

var NoteUpdate = func(w http.ResponseWriter, r *http.Request) {
	Note := &models.Note{}

	params := mux.Vars(r)
	id := params["id"]
	cacheKey := "note:" + id
    	ctx := context.Background()

	database := db.GetDB()
	err := database.First(&Note, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.HandleNotFound(w)
		} else {
			u.HandleBadRequest(w, err)
		}
		return
	}

	newNote := &models.Note{}
	err = json.NewDecoder(r.Body).Decode(newNote)

	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	err = database.Model(&Note).Updates(newNote).Error
	
	if err != nil {
		u.HandleBadRequest(w, err)
	} else {
		db.GetRedis().Del(ctx, cacheKey)
		u.Respond(w, u.Message(true, "OK"))
	}
}

var NoteDelete = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	cacheKey := "note:" + id
    	ctx := context.Background()

	database := db.GetDB()
	err := database.Delete(&models.Note{}, id).Error

	if err != nil {
		u.HandleBadRequest(w, err)
	} else {
		db.GetRedis().Del(ctx, cacheKey)
		u.Respond(w, u.Message(true, "OK"))
	}
}
