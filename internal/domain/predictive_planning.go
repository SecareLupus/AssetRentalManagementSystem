package domain

import (
	"context"
	"fmt"
)

// SeasonRepository defines the data access methods needed for predictive planning
type SeasonRepository interface {
	GetRingsForShow(ctx context.Context, showID int64) ([]ShowRing, error)
	AddRingToShow(ctx context.Context, showRing *ShowRing) error
	SetShowRingLoadout(ctx context.Context, showRingID int64, items []RingLoadoutItem) error
}

// PredictivePlanner handles the business logic of forecasting equipment needs
// for upcoming seasons and shows based on historical configurations.
type PredictivePlanner struct {
	repo SeasonRepository
}

func NewPredictivePlanner(repo SeasonRepository) *PredictivePlanner {
	return &PredictivePlanner{repo: repo}
}

// PredictShowLoadout fetches a historical show's ring configuration and loadouts,
// strips the specific IDs, and returns a template ready to be saved for a new Show.
func (p *PredictivePlanner) PredictShowLoadout(ctx context.Context, historicalShowID int64) ([]ShowRing, error) {
	historicalRings, err := p.repo.GetRingsForShow(ctx, historicalShowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical show rings: %w", err)
	}

	var predictedRings []ShowRing
	for _, hr := range historicalRings {
		// Zero out IDs for the new template
		var newLoadouts []RingLoadoutItem
		for _, loadout := range hr.LoadoutItems {
			newLoadouts = append(newLoadouts, RingLoadoutItem{
				ItemTypeID: loadout.ItemTypeID,
				Quantity:   loadout.Quantity,
			})
		}

		predictedRings = append(predictedRings, ShowRing{
			RingID:       hr.RingID,
			Ring:         hr.Ring, // Keep the Ring details for UI display
			LoadoutItems: newLoadouts,
		})
	}

	return predictedRings, nil
}

// ApplyPredictedLoadout saves the predicted/edited loadout to an actual new Show.
func (p *PredictivePlanner) ApplyPredictedLoadout(ctx context.Context, newShowID int64, rings []ShowRing) error {
	for _, sr := range rings {
		sr.ShowID = newShowID
		if err := p.repo.AddRingToShow(ctx, &sr); err != nil {
			return fmt.Errorf("failed to add ring to show: %w", err)
		}

		if len(sr.LoadoutItems) > 0 {
			if err := p.repo.SetShowRingLoadout(ctx, sr.ID, sr.LoadoutItems); err != nil {
				return fmt.Errorf("failed to set loadout for show ring: %w", err)
			}
		}
	}
	return nil
}
