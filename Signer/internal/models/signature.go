package models

import "time"

// Signature represents the signature of a set of answers given by a user.
type Signature struct {
	ID        int64     // Unique identifier for the signature
	UserEmail string    // Identifier of the user who provided the answers
	Timestamp time.Time // Moment when the signature was made
	Sign      string    // The actual signature (could be a hash)
}
