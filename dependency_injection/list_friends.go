package friends

import (
	"encoding/json"
	"net/http"

	"github.com/chr1sbest/api.mobl.ai/internal/auth"
	"github.com/chr1sbest/api.mobl.ai/internal/logger"
	"github.com/chr1sbest/api.mobl.ai/internal/storage"
	"github.com/chr1sbest/api.mobl.ai/internal/util"
	"github.com/chr1sbest/api.mobl.ai/internal/util/errors"
)

type ListFriendsHandler struct {
	Db   storage.Storage
	Auth auth.Authenticator
}

type ListFriendsRequest struct {
	Token string `json:"token"`
}

func (h *ListFriendsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req ListFriendsRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Sugar.Errorw("Error decoding request body", "error", err)
		util.WriteErrorResponse(w, http.StatusBadRequest, errors.InvalidRequestBody)
		return
	}

	if req.Token == "" {
		util.WriteErrorResponse(w, http.StatusBadRequest, errors.MissingUserToken)
		return
	}

	opts := &auth.ValidateOptions{
		TokenString: req.Token,
	}
	claims, err := h.Auth.ValidateIDToken(ctx, opts)
	if err != nil {
		logger.Sugar.Errorw("Error validating ID token", "error", err)
		util.WriteErrorResponse(w, http.StatusBadRequest, errors.IDTokenInvalid)
		return
	}

	friends, err := h.Db.ListFriends(ctx, claims.Email)
	if err != nil {
		logger.Sugar.Errorw("Failed to list friends", "email", claims.Email, "error", err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, errors.FailedFriendList)
		return
	}

	response := util.SuccessPayload(http.StatusOK, friends)
	util.WriteAPIResponse(w, response)
}
