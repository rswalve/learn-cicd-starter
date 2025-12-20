package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerNotesGet(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetNotesForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user", err)
		return
	}

	postsResp, err := databasePostsToPosts(posts)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert posts", err)
		return
	}

	respondWithJSON(w, http.StatusOK, postsResp)
}

func (cfg *apiConfig) handlerNotesCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Note string `json:"note"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	id := uuid.New().String()
	// 1. Fix the assignment (ensure you aren't reusing := if 'note' was already declared)
	note, err := cfg.DB.CreateNote(r.Context(), database.CreateNoteParams{
		ID:        uuid.New(),       // Ensure this matches the type in models.go
		CreatedAt: time.Now().UTC(), // No longer needs to be cast to string
		UpdatedAt: time.Now().UTC(),
		Note:      params.Note,
		UserID:    user.ID, // Ensure user.ID is also the correct type
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create note")
		return
	}

	note, err := cfg.DB.GetNote(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get note", err)
		return
	}

	noteResp, err := databaseNoteToNote(note)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert note", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, noteResp)
}
