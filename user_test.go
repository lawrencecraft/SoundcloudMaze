package main

import (
	"testing"
)

func TestFollowingAUserAddsYouToTheirFollowers(t *testing.T) {
	user1 := NewUser("001")
	user2 := NewUser("002")

	user1.Follow(user2)

	followers := user2.GetFollowers()

	if len(followers) != 1 {
		t.Error("User2 does not have 1 follower")
	}

	if user1 != followers[0] {
		t.Error("User1 is not user2's follower")
	}
}

func TestUnfollowingAUserRemovesYouFromTheirFollowers(t *testing.T) {
	user1 := NewUser("001")
	user2 := NewUser("002")

	user1.Follow(user2)
	user1.Unfollow(user2)

	followers := user2.GetFollowers()

	if len(followers) != 0 {
		t.Error("Followers was not zero")
	}
}

func TestMultipleFollowersAddsAllFollowers(t *testing.T) {
	user1 := NewUser("001")
	user2 := NewUser("002")
	user3 := NewUser("003")

	user1.Follow(user3)
	user2.Follow(user3)

	followers := user3.GetFollowers()

	if len(followers) != 2 {
		t.Error("Expected 2 followers, got", len(followers))
	}
}

func TestFollowsReturnsTrueForFollower(t *testing.T) {
	user1 := NewUser("001")
	user2 := NewUser("002")

	user1.Follow(user2)

	if !user1.Follows(user2) {
		t.Error("Expected user1 to follow user2")
	}
}

func TestFollowMessageActuallyFollowsTheUser(t *testing.T) {
	uc := NewUserCollection()

	user1 := uc.GetOrCreateUser("001")
	user2 := uc.GetOrCreateUser("002")

	test, _ := uc.UpdateAndGetNotifiees(Message{Type: "F", FromID: "001", ToID: "002", Timestamp: 3})

	if !user1.Follows(user2) {
		t.Error("User1 should follow user2")
	}

	if !(len(test) == 1 && user2 == test[0]) {
		t.Error("only User2 should be notified")
	}
}

func TestUnfollowMessageUnfollowsTheUser(t *testing.T) {
	uc := NewUserCollection()

	user1 := uc.GetOrCreateUser("001")
	user2 := uc.GetOrCreateUser("002")

	user1.Follow(user2)

	test, _ := uc.UpdateAndGetNotifiees(Message{Type: "U", FromID: "001", ToID: "002", Timestamp: 3})

	if user1.Follows(user2) {
		t.Error("Should have unfollowed user2")
	}

	if len(test) > 0 {
		t.Error("No users should be notified")
	}
}
