package mqtt

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/desmond/rental-management-system/internal/db"
)

type CommandHandler struct {
	repo db.Repository
}

func NewCommandHandler(repo db.Repository) *CommandHandler {
	return &CommandHandler{repo: repo}
}

func (h *CommandHandler) HandleCommand(topic string, payload []byte) {
	log.Printf("MQTT: Received command on %s", topic)

	// Topics expected: rms/commands/reserve, rms/commands/update-tag, etc.
	parts := strings.Split(topic, "/")
	if len(parts) < 3 {
		return
	}

	command := parts[len(parts)-1]

	switch command {
	case "reserve":
		h.handleReserve(payload)
	case "update-tag":
		h.handleUpdateTag(payload)
	default:
		log.Printf("MQTT: Unknown command %s", command)
	}
}

func (h *CommandHandler) handleReserve(payload []byte) {
	var req struct {
		ItemTypeID int64  `json:"item_type_id"`
		Quantity   int    `json:"quantity"`
		StartDate  string `json:"start_date"`
		EndDate    string `json:"end_date"`
		ContactID  int64  `json:"contact_id"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("MQTT: Failed to unmarshal reserve command: %v", err)
		return
	}

	log.Printf("MQTT: Processing remote reservation request for ItemType %d (Qty: %d)", req.ItemTypeID, req.Quantity)

	// Implementation note: In a real system, we'd call a Service layer.
	// For now, we interact with the Repo directly or log success.
	// Actually, creating a reservation involves multiple steps if we want to be robust.
	// For this phase, we'll log the "intent" of the command ingest.
}

func (h *CommandHandler) handleUpdateTag(payload []byte) {
	var req struct {
		AssetID int64  `json:"asset_id"`
		NewTag  string `json:"new_tag"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("MQTT: Failed to unmarshal update-tag command: %v", err)
		return
	}

	log.Printf("MQTT: Processing remote tag update for Asset %d -> %s", req.AssetID, req.NewTag)
}
