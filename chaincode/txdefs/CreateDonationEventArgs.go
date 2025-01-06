package txdefs

import "fmt"

// valid the crateDonationEvent 
type CreateDonationEventArgs struct {
    ID          string
    EventName   string
    Recipient   string
    Description string
    Timestamp   string
}

func ValidateCreateDonationEventArgs(args CreateDonationEventArgs) error {
    if args.ID == "" || args.EventName == "" || args.Recipient == "" || args.Timestamp == "" {
        return fmt.Errorf("missing required fields")
    }
    return nil
}
