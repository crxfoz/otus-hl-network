package delivery

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/crxfoz/otus-hl-network/internal/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type ChatHandler struct {
	repo     domain.ChatRepo
	chatHubs map[string]*wsHub
}

func NewChatHandler(repo domain.ChatRepo) *ChatHandler {
	return &ChatHandler{
		repo:     repo,
		chatHubs: make(map[string]*wsHub),
	}
}

func isExist(list []int64, search int64) bool {
	for _, item := range list {
		if item == search {
			return true
		}
	}

	return false
}

type ChatRequest struct {
	Name  string       `json:"name"`
	Users domain.Users `json:"users"`
}

func (h *ChatHandler) CreateChat(c echo.Context, userContext domain.UserContext) error {
	var chatRequest ChatRequest

	if err := c.Bind(&chatRequest); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.HTTPError{Error: err.Error()})
	}

	users := chatRequest.Users
	if !users.Contains(userContext.ID) {
		users = append(users, userContext.ID)
	}

	chatID := uuid.New().String()
	chatName := chatID

	if err := h.repo.CreateChat(chatID, chatName, users); err != nil {
		return c.JSON(http.StatusInternalServerError, domain.HTTPError{Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, domain.HTTPok{Status: "created"})
}

func (h *ChatHandler) newOnSendFn(chatID string, userContext domain.UserContext) func(message domain.Message) {
	return func(msg domain.Message) {
		if err := h.repo.SendMessage(chatID, msg); err != nil {
			logrus.WithError(err).WithField("chatID", chatID).Error("could not save message")
		}
	}
}

func (h *ChatHandler) ConnectChat(c echo.Context, userContext domain.UserContext) error {
	chatID := c.Param("id")

	chatInfo, err := h.repo.GetChat(chatID)
	if err != nil {
		// TODO: check mongo.ErrNoDocuments
		logrus.WithField("chatID", chatID).WithError(err).Error("could not retrieve chat info")
		return c.JSON(http.StatusInternalServerError, "")
	}

	if !chatInfo.Users.Contains(userContext.ID) {
		return c.JSON(http.StatusForbidden, "")
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.WithError(err).Error("could not upgrade connect")
		return c.JSON(http.StatusBadGateway, "")
	}

	var hub *wsHub

	if vv, ok := h.chatHubs[chatID]; ok {
		hub = vv
	} else {
		hub = &wsHub{
			clients:    make(map[*Client]bool),
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			closed:     make(chan struct{}),
		}

		h.chatHubs[chatID] = hub
		go hub.Run()
	}

	client := &Client{
		userName: userContext.Username,
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		onSend:   h.newOnSendFn(chatID, userContext),
	}

	// TODO: wrap logic around websocket

	history, err := h.repo.GetHistory(chatID, 1000)
	if err != nil {
		logrus.WithError(err).WithField("chatID", chatID).Error("could not get chat history")
	} else {
		sort.Sort(history)

		historyPayload, err := json.Marshal(history)
		if err != nil {
			logrus.WithError(err).Error("could not marshal chat history")
		} else {
			client.send <- historyPayload
		}
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()

	<-hub.CloseNotify()
	return c.JSON(http.StatusOK, domain.HTTPok{Status: "ok"})
}
