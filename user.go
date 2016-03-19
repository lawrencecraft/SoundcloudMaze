package main

import "errors"

// User represents a collection of Users
type User struct {
	ID        string
	followers map[*User]bool
}

const (
	errUnknownCommand = "Unknown message type"
)

// Follow another User
func (c *User) Follow(other *User) {
	other.followers[c] = true
}

// Unfollow another user
func (c *User) Unfollow(other *User) {
	delete(other.followers, c)
}

// Follows returns true if a user follows the given
// user, false otherwise
func (c *User) Follows(other *User) bool {
	_, ok := other.followers[c]
	return ok
}

// GetFollowers returns a slice of followers
func (c *User) GetFollowers() []*User {
	followerSlice := make([]*User, len(c.followers))
	index := 0
	for k := range c.followers {
		followerSlice[index] = k
		index++
	}
	return followerSlice
}

// NewUser creates a new user
func NewUser(id string) *User {
	return &User{ID: id, followers: make(map[*User]bool)}
}

// UserCollection represents a collection of users
type UserCollection struct {
	Users   []*User
	userMap map[string]*User
}

func (uc *UserCollection) getUserByID(id string) (*User, bool) {
	user, ok := uc.userMap[id]
	return user, ok
}

func handleUnfollowMessage(uc *UserCollection, m Message) ([]*User, error) {
	fromUser := uc.GetOrCreateUser(m.FromID)
	toUser := uc.GetOrCreateUser(m.ToID)

	fromUser.Unfollow(toUser)
	return nil, nil
}

func handleFollowMessage(uc *UserCollection, m Message) ([]*User, error) {
	fromUser := uc.GetOrCreateUser(m.FromID)
	toUser := uc.GetOrCreateUser(m.ToID)

	fromUser.Follow(toUser)
	return []*User{toUser}, nil
}

func handleBroadcastMessage(uc *UserCollection, m Message) ([]*User, error) {
	return uc.Users, nil
}

func handlePrivateMessage(uc *UserCollection, m Message) ([]*User, error) {
	return []*User{uc.GetOrCreateUser(m.ToID)}, nil
}

func handleStatusUpdate(uc *UserCollection, m Message) ([]*User, error) {
	user := uc.GetOrCreateUser(m.FromID)
	return user.GetFollowers(), nil
}

// UpdateAndGetNotifiees handles the message and returns the users who need to be notified
func (uc *UserCollection) UpdateAndGetNotifiees(m Message) ([]*User, error) {
	switch m.Type {
	case "U":
		return handleUnfollowMessage(uc, m)
	case "F":
		return handleFollowMessage(uc, m)
	case "B":
		return handleBroadcastMessage(uc, m)
	case "P":
		return handlePrivateMessage(uc, m)
	case "S":
		return handleStatusUpdate(uc, m)
	default:
		return nil, errors.New(errUnknownCommand)
	}
}

// GetOrCreateUser gets or creates a user
func (uc *UserCollection) GetOrCreateUser(id string) *User {
	user, ok := uc.userMap[id]
	if !ok {
		user = NewUser(id)
		uc.Users = append(uc.Users, user)
		uc.userMap[id] = user
	}
	return user
}

// NewUserCollection creates an initialized UserCollection
func NewUserCollection() UserCollection {
	return UserCollection{Users: []*User{}, userMap: make(map[string]*User)}
}
