package room

import (
	"net/http"
	"time"

	"github.com/m1k1o/neko/server/pkg/utils"
)

// KeySequenceRequest represents a sequence of key combos to be injected.
// Each inner slice represents keys pressed simultaneously (a combo).
type KeySequenceRequest struct {
	Sequences [][]uint32 `json:"sequences"`
	// Delay in milliseconds between combos
	DelayMs int `json:"delay_ms,omitempty"`
}

// controlKeys allows the host (or admins) to inject a sequence of key
// combinations into the desktop. Each combo (inner slice) is passed to
// DesktopManager.KeyPress which presses and releases the provided keysyms.
func (h *RoomHandler) controlKeys(w http.ResponseWriter, r *http.Request) error {
	payload := KeySequenceRequest{}
	if err := utils.HttpJsonRequest(w, r, &payload); err != nil {
		return err
	}

	if len(payload.Sequences) == 0 {
		return utils.HttpBadRequest("no sequences provided")
	}

	for _, combo := range payload.Sequences {
		if len(combo) == 0 {
			// skip empty combos
			continue
		}

		if err := h.desktop.KeyPress(combo...); err != nil {
			return utils.HttpInternalServerError().WithInternalErr(err)
		}

		if payload.DelayMs > 0 {
			time.Sleep(time.Duration(payload.DelayMs) * time.Millisecond)
		}
	}

	return utils.HttpSuccess(w)
}
