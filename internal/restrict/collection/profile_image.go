package collection

import (
    "context"
    "fmt"

    "github.com/OpenSlides/openslides-go/datastore/dsfetch"
)

// ProfileImage handles permissions for the profile_image collection.
//
// Profile images can be seen if the request user can see the linked user.
//
// Mode A: The user can see the profile image.
type ProfileImage struct{}

// Name returns the collection name.
func (p ProfileImage) Name() string {
    return "profile_image"
}

// MeetingID returns no meeting.
func (p ProfileImage) MeetingID(ctx context.Context, ds *dsfetch.Fetch, id int) (int, bool, error) {
    return 0, false, nil
}

// Modes returns the field modes for the collection profile_image.
func (p ProfileImage) Modes(mode string) FieldRestricter {
    switch mode {
    case "A":
        return p.see
    }
    return nil
}

func (p ProfileImage) see(ctx context.Context, ds *dsfetch.Fetch, profileImageIDs ...int) ([]int, error) {
    return eachCondition(profileImageIDs, func(profileImageID int) (bool, error) {
        userID, err := ds.ProfileImage_UserID(profileImageID).Value(ctx)
        if err != nil {
            return false, fmt.Errorf("getting user_id of profile_image %d: %w", profileImageID, err)
        }

        canSeeUser, err := Collection(ctx, User{}.Name()).Modes("A")(ctx, ds, userID)
        if err != nil {
            return false, fmt.Errorf("can see user of profile image %d: %w", profileImageID, err)
        }

        return len(canSeeUser) >= 1, nil
    })
}
